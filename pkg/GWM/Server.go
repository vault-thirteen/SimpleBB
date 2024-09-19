package gwm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	a "github.com/vault-thirteen/SimpleBB/pkg/ACM"
	"github.com/vault-thirteen/SimpleBB/pkg/GWM/dbo"
	"github.com/vault-thirteen/SimpleBB/pkg/GWM/models"
	"github.com/vault-thirteen/SimpleBB/pkg/GWM/models/api"
	gs "github.com/vault-thirteen/SimpleBB/pkg/GWM/settings"
	m "github.com/vault-thirteen/SimpleBB/pkg/MM"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	"github.com/vault-thirteen/SimpleBB/pkg/common/app"
	"github.com/vault-thirteen/SimpleBB/pkg/common/avm"
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/client"
	ch "github.com/vault-thirteen/SimpleBB/pkg/common/http"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cn "github.com/vault-thirteen/SimpleBB/pkg/common/net"
	cset "github.com/vault-thirteen/SimpleBB/pkg/common/settings"
)

const ErrUrlIsTooShort = "URL is too short"

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
	captchaProxy     *httputil.ReverseProxy

	// Mapping of HTTP status codes by RPC error code for various services.
	commonHttpStatusCodesByRpcErrorCode map[int]int
	acmHttpStatusCodesByRpcErrorCode    map[int]int
	mmHttpStatusCodesByRpcErrorCode     map[int]int

	// Public settings.
	publicSettingsFileData []byte // Cached file contents in JSON format.

	// Front End Data.
	frontEnd *models.FrontEndData
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

	err = srv.initPublicSettings()
	if err != nil {
		return nil, err
	}

	if srv.settings.SystemSettings.IsFrontEndEnabled {
		err = srv.initFrontEndData()
		if err != nil {
			return nil, err
		}
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
	// Firewall (optional).
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

	isFrontEndEnabled := srv.settings.SystemSettings.IsFrontEndEnabled

	if isFrontEndEnabled {
		switch req.URL.Path {
		case gs.FrontEndRoot: // <- /
			srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.IndexHtmlPage)
			return

		case srv.frontEnd.FavIcon.UrlPath: // <- /favicon.png
			srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.FavIcon)
			return
		}
	}

	var err error
	urlParts := cn.SplitUrlPath(req.URL.Path)
	if len(urlParts) < 1 {
		err = errors.New(ErrUrlIsTooShort)
		srv.logError(err)
		return
	}
	category := urlParts[0]

	switch category {
	case srv.settings.SystemSettings.ApiFolder: // <- /api
		srv.handlePublicApi(rw, req, clientIPA)
		return

	case srv.settings.SystemSettings.PublicSettingsFileName: // <- /settings.json
		srv.handlePublicSettings(rw, req)
		return

	case srv.settings.SystemSettings.CaptchaFolder: // <- /captcha
		srv.handleCaptcha(rw, req)
		return

	case srv.settings.SystemSettings.FrontEndStaticFilesFolder: // <- /fe
		if !isFrontEndEnabled {
			srv.respondNotFound(rw)
			return
		}
		srv.handleFrontEnd(rw, req, clientIPA)
		return

	case gs.FrontEndStaticFileName_AdminHtmlPage: // <- /admin.html
		srv.handleAdminFrontEnd(rw, req, clientIPA)
		return

	default:
		srv.respondNotFound(rw)
		return
	}
}

func (srv *Server) initStatusCodeMapper() (err error) {
	srv.commonHttpStatusCodesByRpcErrorCode = c.GetMapOfHttpStatusCodesByRpcErrorCodes()
	srv.acmHttpStatusCodesByRpcErrorCode = a.GetMapOfHttpStatusCodesByRpcErrorCodes()
	srv.mmHttpStatusCodesByRpcErrorCode = m.GetMapOfHttpStatusCodesByRpcErrorCodes()

	return nil
}

func (srv *Server) initPublicSettings() (err error) {
	// File with public settings is a special virtual file which lies in the
	// root folder. It contains useful settings for client applications.
	var publicSettings = &models.Settings{
		Version:                   srv.settings.SystemSettings.SettingsVersion,
		ProductVersion:            srv.settings.VersionInfo.ProgramVersionString(),
		SiteName:                  srv.settings.SystemSettings.SiteName,
		SiteDomain:                srv.settings.SystemSettings.SiteDomain,
		CaptchaFolder:             srv.settings.SystemSettings.CaptchaFolder,
		SessionMaxDuration:        srv.settings.SystemSettings.SessionMaxDuration,
		MessageEditTime:           srv.settings.SystemSettings.MessageEditTime,
		PageSize:                  srv.settings.SystemSettings.PageSize,
		ApiFolder:                 srv.settings.SystemSettings.ApiFolder,
		PublicSettingsFileName:    srv.settings.SystemSettings.PublicSettingsFileName,
		IsFrontEndEnabled:         srv.settings.SystemSettings.IsFrontEndEnabled,
		FrontEndStaticFilesFolder: srv.settings.SystemSettings.FrontEndStaticFilesFolder,
	}

	srv.publicSettingsFileData, err = json.Marshal(publicSettings)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) initFrontEndData() (err error) {
	frontendAssetsFolder := cm.NormalisePath(srv.settings.SystemSettings.FrontEndAssetsFolder)
	fep := gs.FrontEndRoot + srv.settings.SystemSettings.FrontEndStaticFilesFolder + "/"

	srv.frontEnd = &models.FrontEndData{}

	srv.frontEnd.AdminHtmlPage, err = models.NewFrontEndFileData(gs.FrontEndRoot, gs.FrontEndStaticFileName_AdminHtmlPage, ch.ContentType_HtmlPage, frontendAssetsFolder)
	if err != nil {
		return err
	}

	srv.frontEnd.AdminJs, err = models.NewFrontEndFileData(fep, gs.FrontEndStaticFileName_AdminJs, ch.ContentType_JavaScript, frontendAssetsFolder)
	if err != nil {
		return err
	}

	srv.frontEnd.ArgonJs, err = models.NewFrontEndFileData(fep, gs.FrontEndStaticFileName_ArgonJs, ch.ContentType_JavaScript, frontendAssetsFolder)
	if err != nil {
		return err
	}

	srv.frontEnd.ArgonWasm, err = models.NewFrontEndFileData(fep, gs.FrontEndStaticFileName_ArgonWasm, ch.ContentType_Wasm, frontendAssetsFolder)
	if err != nil {
		return err
	}

	srv.frontEnd.BppJs, err = models.NewFrontEndFileData(fep, gs.FrontEndStaticFileName_BppJs, ch.ContentType_JavaScript, frontendAssetsFolder)
	if err != nil {
		return err
	}

	srv.frontEnd.CssStyles, err = models.NewFrontEndFileData(fep, gs.FrontEndStaticFileName_CssStyles, ch.ContentType_CssStyle, frontendAssetsFolder)
	if err != nil {
		return err
	}

	srv.frontEnd.FavIcon, err = models.NewFrontEndFileData(gs.FrontEndRoot, gs.FrontEndStaticFileName_FavIcon, ch.ContentType_PNG, frontendAssetsFolder)
	if err != nil {
		return err
	}

	srv.frontEnd.IndexHtmlPage, err = models.NewFrontEndFileData(fep, gs.FrontEndStaticFileName_IndexHtmlPage, ch.ContentType_HtmlPage, frontendAssetsFolder)
	if err != nil {
		return err
	}

	srv.frontEnd.LoaderScript, err = models.NewFrontEndFileData(fep, gs.FrontEndStaticFileName_LoaderScript, ch.ContentType_JavaScript, frontendAssetsFolder)
	if err != nil {
		return err
	}

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

	// Proxy for captcha images.
	targetAddr := cc.UrlSchemeHttp + "://" + net.JoinHostPort(srv.settings.SystemSettings.CaptchaImgServerHost, strconv.FormatUint(uint64(srv.settings.SystemSettings.CaptchaImgServerPort), 10))
	var targetUrl *url.URL
	targetUrl, err = url.Parse(targetAddr)
	if err != nil {
		return err
	}

	srv.captchaProxy = httputil.NewSingleHostReverseProxy(targetUrl)

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
		// ACM.
		ApiFunctionName_RegisterUser,
		ApiFunctionName_ApproveAndRegisterUser,
		ApiFunctionName_LogUserIn,
		ApiFunctionName_LogUserOut,
		ApiFunctionName_GetListOfLoggedUsers,
		ApiFunctionName_IsUserLoggedIn,
		ApiFunctionName_ChangePassword,
		ApiFunctionName_ChangeEmail,
		ApiFunctionName_GetUserRoles,
		ApiFunctionName_ViewUserParameters,
		ApiFunctionName_SetUserRoleAuthor,
		ApiFunctionName_SetUserRoleWriter,
		ApiFunctionName_SetUserRoleReader,
		ApiFunctionName_GetSelfRoles,
		ApiFunctionName_BanUser,
		ApiFunctionName_UnbanUser,

		// MM.
		ApiFunctionName_AddSection,
		ApiFunctionName_ChangeSectionName,
		ApiFunctionName_ChangeSectionParent,
		ApiFunctionName_GetSection,
		ApiFunctionName_DeleteSection,
		ApiFunctionName_AddForum,
		ApiFunctionName_ChangeForumName,
		ApiFunctionName_ChangeForumSection,
		ApiFunctionName_GetForum,
		ApiFunctionName_DeleteForum,
		ApiFunctionName_AddThread,
		ApiFunctionName_ChangeThreadName,
		ApiFunctionName_ChangeThreadForum,
		ApiFunctionName_GetThread,
		ApiFunctionName_DeleteThread,
		ApiFunctionName_AddMessage,
		ApiFunctionName_ChangeMessageText,
		ApiFunctionName_ChangeMessageThread,
		ApiFunctionName_GetMessage,
		ApiFunctionName_DeleteMessage,
		ApiFunctionName_ListThreadAndMessages,
		ApiFunctionName_ListThreadAndMessagesOnPage,
		ApiFunctionName_ListForumAndThreads,
		ApiFunctionName_ListForumAndThreadsOnPage,
		ApiFunctionName_ListSectionsAndForums,
	}

	srv.apiHandlers = map[string]api.RequestHandler{
		// ACM.
		ApiFunctionName_RegisterUser:                           srv.RegisterUser,
		ApiFunctionName_GetListOfRegistrationsReadyForApproval: srv.GetListOfRegistrationsReadyForApproval,
		ApiFunctionName_ApproveAndRegisterUser:                 srv.ApproveAndRegisterUser,
		ApiFunctionName_LogUserIn:                              srv.LogUserIn,
		ApiFunctionName_LogUserOut:                             srv.LogUserOut,
		ApiFunctionName_GetListOfLoggedUsers:                   srv.GetListOfLoggedUsers,
		ApiFunctionName_IsUserLoggedIn:                         srv.IsUserLoggedIn,
		ApiFunctionName_ChangePassword:                         srv.ChangePassword,
		ApiFunctionName_ChangeEmail:                            srv.ChangeEmail,
		ApiFunctionName_GetListOfAllUsers:                      srv.GetListOfAllUsers,
		ApiFunctionName_GetUserRoles:                           srv.GetUserRoles,
		ApiFunctionName_ViewUserParameters:                     srv.ViewUserParameters,
		ApiFunctionName_SetUserRoleAuthor:                      srv.SetUserRoleAuthor,
		ApiFunctionName_SetUserRoleWriter:                      srv.SetUserRoleWriter,
		ApiFunctionName_SetUserRoleReader:                      srv.SetUserRoleReader,
		ApiFunctionName_GetSelfRoles:                           srv.GetSelfRoles,
		ApiFunctionName_BanUser:                                srv.BanUser,
		ApiFunctionName_UnbanUser:                              srv.UnbanUser,

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
