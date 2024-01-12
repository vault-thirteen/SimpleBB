package gwm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	a "github.com/vault-thirteen/SimpleBB/pkg/ACM"
	"github.com/vault-thirteen/SimpleBB/pkg/GWM/dbo"
	"github.com/vault-thirteen/SimpleBB/pkg/GWM/models/api"
	gs "github.com/vault-thirteen/SimpleBB/pkg/GWM/settings"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	"github.com/vault-thirteen/SimpleBB/pkg/common/app"
	"github.com/vault-thirteen/SimpleBB/pkg/common/avm"
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/client"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cset "github.com/vault-thirteen/SimpleBB/pkg/common/settings"
)

type Server struct {
	// Settings.
	settings *gs.Settings

	// HTTP server for internal requests.
	listenDsnInt  string
	httpServerInt *http.Server

	// HTTPS server for external requests.
	listenDsnExt     string
	httpServerExt    *http.Server
	apiFunctionNames []string
	apiHandlers      map[string]api.RequestHandler

	// Channel for an external controller. When a message comes from this
	// channel, a controller must stop this server. The server does not stop
	// itself.
	mustBeStopped chan bool

	// Internal control structures.
	subRoutines *sync.WaitGroup
	mustStop    *atomic.Bool
	httpErrors  chan error
	dbErrors    *chan error
	ssp         *avm.SSP

	// Database Object.
	dbo *dbo.DatabaseObject

	// JSON-RPC server.
	js *jrm1.Processor

	// Clients for external services.
	acmServiceClient *cc.Client
	mmServiceClient  *cc.Client

	// Mapping of HTTP status codes by RPC error code for various services.
	commonHttpStatusCodesByRpcErrorCode map[int]int
	acmHttpStatusCodesByRpcErrorCode    map[int]int
}

func NewServer(s cm.ISettings) (srv *Server, err error) {
	stn := s.(*gs.Settings)

	err = stn.Check()
	if err != nil {
		return nil, err
	}

	dbErrorsChannel := make(chan error, c.DbErrorsChannelSize)

	srv = &Server{
		settings:      stn,
		listenDsnInt:  net.JoinHostPort(stn.IntHttpSettings.Host, strconv.FormatUint(uint64(stn.IntHttpSettings.Port), 10)),
		listenDsnExt:  net.JoinHostPort(stn.ExtHttpsSettings.Host, strconv.FormatUint(uint64(stn.ExtHttpsSettings.Port), 10)),
		mustBeStopped: make(chan bool, c.MustBeStoppedChannelSize),
		subRoutines:   new(sync.WaitGroup),
		mustStop:      new(atomic.Bool),
		httpErrors:    make(chan error, c.HttpErrorsChannelSize),
		dbErrors:      &dbErrorsChannel,
		ssp:           avm.NewSSP(),
	}
	srv.mustStop.Store(false)

	if srv.settings.SystemSettings.IsFirewallUsed {
		fmt.Println(c.MsgFirewallIsEnabled)
	} else {
		fmt.Println(c.MsgFirewallIsDisabled)
	}

	// RPC server.
	err = srv.initRpc()
	if err != nil {
		return nil, err
	}

	// Database.
	srv.dbo = dbo.NewDatabaseObject(srv.settings.DbSettings)

	err = srv.dbo.Init()
	if err != nil {
		return nil, err
	}

	// HTTP server for internal requests.
	srv.httpServerInt = &http.Server{
		Addr:    srv.listenDsnInt,
		Handler: http.Handler(http.HandlerFunc(srv.httpRouterInt)),
	}

	// HTTPS server for external requests.
	srv.httpServerExt = &http.Server{
		Addr:    srv.listenDsnExt,
		Handler: http.Handler(http.HandlerFunc(srv.httpRouterExt)),
	}

	err = srv.initApiFunctions()
	if err != nil {
		return nil, err
	}

	err = srv.createClientsForExternalServices()
	if err != nil {
		return nil, err
	}

	err = srv.initStatusCodeMapper()
	if err != nil {
		return nil, err
	}

	return srv, nil
}

func (srv *Server) GetListenDsnInt() (dsn string) {
	return srv.listenDsnInt
}

func (srv *Server) GetListenDsnExt() (dsn string) {
	return srv.listenDsnExt
}

func (srv *Server) GetStopChannel() *chan bool {
	return &srv.mustBeStopped
}

func (srv *Server) Start() (err error) {
	srv.ssp.Lock()
	defer srv.ssp.Unlock()

	err = srv.ssp.BeginStart()
	if err != nil {
		return err
	}

	srv.startHttpServerInt()
	srv.startHttpServerExt()

	srv.subRoutines.Add(3)
	go srv.listenForHttpErrors()
	go srv.listenForDbErrors()
	go srv.runScheduler()

	err = srv.pingClientsForExternalServices()
	if err != nil {
		return err
	}

	srv.ssp.CompleteStart()

	return nil
}

func (srv *Server) Stop() (err error) {
	srv.ssp.Lock()
	defer srv.ssp.Unlock()

	err = srv.ssp.BeginStop()
	if err != nil {
		return err
	}

	srv.mustStop.Store(true)

	ctxInt, cfInt := context.WithTimeout(context.Background(), time.Minute)
	defer cfInt()
	err = srv.httpServerInt.Shutdown(ctxInt)
	if err != nil {
		return err
	}

	ctxExt, cfExt := context.WithTimeout(context.Background(), time.Minute)
	defer cfExt()
	err = srv.httpServerExt.Shutdown(ctxExt)
	if err != nil {
		return err
	}

	close(srv.httpErrors)
	close(*srv.dbErrors)

	srv.subRoutines.Wait()

	err = srv.dbo.Fin()
	if err != nil {
		return err
	}

	srv.ssp.CompleteStop()

	return nil
}

func (srv *Server) startHttpServerInt() {
	go func() {
		var listenError error
		listenError = srv.httpServerInt.ListenAndServe()

		if (listenError != nil) && (!errors.Is(listenError, http.ErrServerClosed)) {
			srv.httpErrors <- listenError
		}
	}()
}

func (srv *Server) startHttpServerExt() {
	go func() {
		var listenError error
		listenError = srv.httpServerExt.ListenAndServeTLS(srv.settings.ExtHttpsSettings.CertFile, srv.settings.ExtHttpsSettings.KeyFile)

		if (listenError != nil) && (!errors.Is(listenError, http.ErrServerClosed)) {
			srv.httpErrors <- listenError
		}
	}()
}

func (srv *Server) listenForHttpErrors() {
	defer srv.subRoutines.Done()

	for err := range srv.httpErrors {
		log.Println(c.MsgServerError + err.Error())
		srv.mustBeStopped <- true
	}

	log.Println(c.MsgHttpErrorListenerHasStopped)
}

func (srv *Server) listenForDbErrors() {
	defer srv.subRoutines.Done()

	var err error
	for dbErr := range *srv.dbErrors {
		// When a network error occurs, it may be followed by numerous other
		// errors. If we try to fix each of them, we can make a flood. So,
		// we make a smart thing here.

		// 1. Ensure that the problem still exists.
		err = srv.dbo.ProbeDb()
		if err == nil {
			// Network is now fine. Ignore the error.
			continue
		}

		// 2. Log the error and try to reconnect.
		log.Println(c.MsgDatabaseNetworkError + dbErr.Error())

		for {
			log.Println(c.MsgReconnectingDatabase)
			// While we have prepared statements,
			// the simple reconnect will not help.
			err = srv.dbo.Init()
			if err != nil {
				// Network is still bad.
				log.Println(c.MsgReconnectionHasFailed + err.Error())
			} else {
				log.Println(c.MsgConnectionToDatabaseWasRestored)
				break
			}

			time.Sleep(time.Second * c.DbReconnectCoolDownPeriodSec)
		}
	}

	log.Println(c.MsgDbNetworkErrorListenerHasStopped)
}

func (srv *Server) runScheduler() {
	defer srv.subRoutines.Done()

	var funcsToRunEveryMinute = []c.ScheduledFnSimple{
		srv.clearIPAddresses,
	}

	// Time counter.
	// It counts seconds and resets every 24 hours.
	var tc uint = 1
	const SecondsInDay = 86400 // 60*60*24.
	var err error

	for {
		if srv.mustStop.Load() {
			break
		}

		// Periodical tasks.

		// Every minute.
		if tc%60 == 0 {
			for _, fn := range funcsToRunEveryMinute {
				err = fn()
				if err != nil {
					log.Println(err)
				}
			}
		}

		// Every hour.
		if tc%3600 == 0 {
			// ...
		}

		// Next tick.
		if tc == SecondsInDay {
			tc = 0
		}
		tc++
		time.Sleep(time.Second)
	}

	log.Println(c.MsgSchedulerHasStopped)
}

// HTTP router for internal requests.
func (srv *Server) httpRouterInt(rw http.ResponseWriter, req *http.Request) {
	srv.js.ServeHTTP(rw, req)
}

// HTTP router for external requests.
func (srv *Server) httpRouterExt(rw http.ResponseWriter, req *http.Request) {
	// Firewall.
	var clientIPA string
	if srv.settings.SystemSettings.IsFirewallUsed {
		var ok bool
		var err error
		ok, clientIPA, err = srv.isIPAddressAllowed(req)
		if err != nil {
			srv.processInternalServerError(rw, err)
			return
		}

		if !ok {
			srv.respondForbidden(rw)
			return
		}
	}

	// API handlers.
	if req.URL.Path == srv.settings.SystemSettings.ApiPath {
		srv.handlePublicApi(rw, req, clientIPA)
		return
	}

	// Front end handlers (optional).
	if srv.settings.SystemSettings.IsFrontEndEnabled {
		if req.URL.Path == srv.settings.SystemSettings.FrontEndPath {
			srv.handlePublicFrontEnd(rw, req)
			return
		}
	}

	srv.respondNotFound(rw)
}

func (srv *Server) initStatusCodeMapper() (err error) {
	srv.acmHttpStatusCodesByRpcErrorCode = a.GetMapOfHttpStatusCodesByRpcErrorCodes()
	srv.commonHttpStatusCodesByRpcErrorCode = c.GetMapOfHttpStatusCodesByRpcErrorCodes()

	return nil
}

func (srv *Server) createClientsForExternalServices() (err error) {
	// ACM module.
	var acmSCS = &cset.ServiceClientSettings{
		Schema:                      srv.settings.AcmSettings.Schema,
		Host:                        srv.settings.AcmSettings.Host,
		Port:                        srv.settings.AcmSettings.Port,
		Path:                        srv.settings.AcmSettings.Path,
		EnableSelfSignedCertificate: srv.settings.AcmSettings.EnableSelfSignedCertificate,
	}

	srv.acmServiceClient, err = cc.NewClientWithSCS(acmSCS, app.ServiceShortName_ACM)
	if err != nil {
		return err
	}

	// MM module.
	var mmSCS = &cset.ServiceClientSettings{
		Schema:                      srv.settings.MmSettings.Schema,
		Host:                        srv.settings.MmSettings.Host,
		Port:                        srv.settings.MmSettings.Port,
		Path:                        srv.settings.MmSettings.Path,
		EnableSelfSignedCertificate: srv.settings.MmSettings.EnableSelfSignedCertificate,
	}

	srv.mmServiceClient, err = cc.NewClientWithSCS(mmSCS, app.ServiceShortName_MM)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) pingClientsForExternalServices() (err error) {
	// ACM module.
	err = srv.acmServiceClient.Ping(true)
	if err != nil {
		return err
	}

	// MM module.
	err = srv.mmServiceClient.Ping(true)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) initApiFunctions() (err error) {
	srv.apiFunctionNames = []string{
		ApiFunctionName_GetProductVersion,
	}

	srv.apiHandlers = map[string]api.RequestHandler{
		ApiFunctionName_GetProductVersion: srv.GetProductVersion,

		// ACM.
		ApiFunctionName_RegisterUser:           srv.RegisterUser,
		ApiFunctionName_ApproveAndRegisterUser: srv.ApproveAndRegisterUser,
		ApiFunctionName_LogUserIn:              srv.LogUserIn,
		ApiFunctionName_LogUserOut:             srv.LogUserOut,
		ApiFunctionName_GetListOfLoggedUsers:   srv.GetListOfLoggedUsers,
		ApiFunctionName_IsUserLoggedIn:         srv.IsUserLoggedIn,
		ApiFunctionName_ChangePassword:         srv.ChangePassword,
		ApiFunctionName_ChangeEmail:            srv.ChangeEmail,
		ApiFunctionName_GetUserRoles:           srv.GetUserRoles,
		ApiFunctionName_ViewUserParameters:     srv.ViewUserParameters,
		ApiFunctionName_SetUserRoleAuthor:      srv.SetUserRoleAuthor,
		ApiFunctionName_SetUserRoleWriter:      srv.SetUserRoleWriter,
		ApiFunctionName_SetUserRoleReader:      srv.SetUserRoleReader,
		ApiFunctionName_GetSelfRoles:           srv.GetSelfRoles,
		ApiFunctionName_BanUser:                srv.BanUser,
		ApiFunctionName_UnbanUser:              srv.UnbanUser,

		// MM.
		ApiFunctionName_AddSection:                  srv.AddSection,
		ApiFunctionName_ChangeSectionName:           srv.ChangeSectionName,
		ApiFunctionName_ChangeSectionParent:         srv.ChangeSectionParent,
		ApiFunctionName_GetSection:                  srv.GetSection,
		ApiFunctionName_DeleteSection:               srv.DeleteSection,
		ApiFunctionName_AddForum:                    srv.AddForum,
		ApiFunctionName_ChangeForumName:             srv.ChangeForumName,
		ApiFunctionName_ChangeForumSection:          srv.ChangeForumSection,
		ApiFunctionName_GetForum:                    srv.GetForum,
		ApiFunctionName_DeleteForum:                 srv.DeleteForum,
		ApiFunctionName_AddThread:                   srv.AddThread,
		ApiFunctionName_ChangeThreadName:            srv.ChangeThreadName,
		ApiFunctionName_ChangeThreadForum:           srv.ChangeThreadForum,
		ApiFunctionName_GetThread:                   srv.GetThread,
		ApiFunctionName_DeleteThread:                srv.DeleteThread,
		ApiFunctionName_AddMessage:                  srv.AddMessage,
		ApiFunctionName_ChangeMessageText:           srv.ChangeMessageText,
		ApiFunctionName_ChangeMessageThread:         srv.ChangeMessageThread,
		ApiFunctionName_GetMessage:                  srv.GetMessage,
		ApiFunctionName_DeleteMessage:               srv.DeleteMessage,
		ApiFunctionName_ListThreadAndMessages:       srv.ListThreadAndMessages,
		ApiFunctionName_ListThreadAndMessagesOnPage: srv.ListThreadAndMessagesOnPage,
		ApiFunctionName_ListForumAndThreads:         srv.ListForumAndThreads,
		ApiFunctionName_ListForumAndThreadsOnPage:   srv.ListForumAndThreadsOnPage,
		ApiFunctionName_ListSectionsAndForums:       srv.ListSectionsAndForums,
	}

	return nil
}

func (srv *Server) ReportStart() {
	fmt.Println(c.MsgHttpServer + srv.GetListenDsnInt())
	fmt.Println(c.MsgHttpsServer + srv.GetListenDsnExt())
}

func (srv *Server) UseConstructor(stn cm.ISettings) (cm.IServer, error) {
	return NewServer(stn)
}
