package mm

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/SimpleBB/pkg/MM/dbo"
	ms "github.com/vault-thirteen/SimpleBB/pkg/MM/settings"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	"github.com/vault-thirteen/SimpleBB/pkg/common/app"
	"github.com/vault-thirteen/SimpleBB/pkg/common/avm"
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/client"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cset "github.com/vault-thirteen/SimpleBB/pkg/common/settings"
)

type Server struct {
	// Settings.
	settings *ms.Settings

	// HTTP server.
	listenDsn  string
	httpServer *http.Server

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
	nmServiceClient  *cc.Client

	// CRC32 Table.
	crcTable *crc32.Table

	// DKeys.
	dKeyForNM *string
}

func NewServer(s cm.ISettings) (srv *Server, err error) {
	stn := s.(*ms.Settings)

	err = stn.Check()
	if err != nil {
		return nil, err
	}

	dbErrorsChannel := make(chan error, c.DbErrorsChannelSize)

	srv = &Server{
		settings:      stn,
		listenDsn:     net.JoinHostPort(stn.HttpsSettings.Host, strconv.FormatUint(uint64(stn.HttpsSettings.Port), 10)),
		mustBeStopped: make(chan bool, c.MustBeStoppedChannelSize),
		subRoutines:   new(sync.WaitGroup),
		mustStop:      new(atomic.Bool),
		httpErrors:    make(chan error, c.HttpErrorsChannelSize),
		dbErrors:      &dbErrorsChannel,
		ssp:           avm.NewSSP(),
	}
	srv.mustStop.Store(false)

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

	// HTTP Server.
	srv.httpServer = &http.Server{
		Addr:    srv.listenDsn,
		Handler: http.Handler(http.HandlerFunc(srv.httpRouter)),
	}

	err = srv.initCRC()
	if err != nil {
		return nil, err
	}

	err = srv.createClientsForExternalServices()
	if err != nil {
		return nil, err
	}

	return srv, nil
}

func (srv *Server) GetListenDsn() (dsn string) {
	return srv.listenDsn
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

	srv.startHttpServer()

	srv.subRoutines.Add(2)
	go srv.listenForHttpErrors()
	go srv.listenForDbErrors()

	err = srv.pingClientsForExternalServices()
	if err != nil {
		return err
	}

	err = srv.synchroniseModules(true)
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

	ctx, cf := context.WithTimeout(context.Background(), time.Minute)
	defer cf()
	err = srv.httpServer.Shutdown(ctx)
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

func (srv *Server) GetSubRoutinesWG() *sync.WaitGroup {
	return srv.subRoutines
}

func (srv *Server) GetMustStopAB() *atomic.Bool {
	return srv.mustStop
}

func (srv *Server) startHttpServer() {
	go func() {
		var listenError error
		listenError = srv.httpServer.ListenAndServeTLS(srv.settings.HttpsSettings.CertFile, srv.settings.HttpsSettings.KeyFile)

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

func (srv *Server) httpRouter(rw http.ResponseWriter, req *http.Request) {
	srv.js.ServeHTTP(rw, req)
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

	// NM module.
	var nmSCS = &cset.ServiceClientSettings{
		Schema:                      srv.settings.NmSettings.Schema,
		Host:                        srv.settings.NmSettings.Host,
		Port:                        srv.settings.NmSettings.Port,
		Path:                        srv.settings.NmSettings.Path,
		EnableSelfSignedCertificate: srv.settings.NmSettings.EnableSelfSignedCertificate,
	}

	srv.nmServiceClient, err = cc.NewClientWithSCS(nmSCS, app.ServiceShortName_NM)
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

	// NM module.
	err = srv.nmServiceClient.Ping(true)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) synchroniseModules(verbose bool) (err error) {
	// NM module.
	if verbose {
		fmt.Print(fmt.Sprintf(c.MsgFSynchronisingWithModule, app.ServiceShortName_NM))
	}

	var re *jrm1.RpcError
	srv.dKeyForNM, re = srv.getDKeyForNM()
	if re != nil {
		return re.AsError()
	}

	if verbose {
		fmt.Println(c.MsgOK)
	}

	return nil
}

func (srv *Server) initCRC() (err error) {
	srv.crcTable = crc32.IEEETable

	return nil
}

func (srv *Server) ReportStart() {
	fmt.Println(c.MsgHttpsServer + srv.GetListenDsn())
}

func (srv *Server) UseConstructor(stn cm.ISettings) (cm.IServer, error) {
	return NewServer(stn)
}
