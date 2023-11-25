package mm

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

	js "github.com/osamingo/jsonrpc/v2"
	ac "github.com/vault-thirteen/SimpleBB/pkg/ACM/client"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	"github.com/vault-thirteen/SimpleBB/pkg/MM/dbo"
	s "github.com/vault-thirteen/SimpleBB/pkg/MM/settings"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cdd "github.com/vault-thirteen/SimpleBB/pkg/common/DiagnosticData"
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/client"
)

const (
	DbReconnectCoolDownPeriodSec = 15
)

const (
	ErrACMModuleIsBroken = "ACM module is broken"
)

type Server struct {
	// Settings.
	settings *s.Settings

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

	// Database Object.
	dbo *dbo.DatabaseObject

	// JSON-RPC handlers.
	jsonRpcHandlers *js.MethodRepository

	// ACM service client.
	acmServiceClient *cc.Client

	// Diagnostic data.
	diag *cdd.DiagnosticData
}

func NewServer(stn *s.Settings) (srv *Server, err error) {
	err = stn.Check()
	if err != nil {
		return nil, err
	}

	dbErrorsChannel := make(chan error, c.DbErrorsChannelSize)

	srv = &Server{
		settings:        stn,
		listenDsn:       net.JoinHostPort(stn.HttpsSettings.Host, strconv.FormatUint(uint64(stn.HttpsSettings.Port), 10)),
		mustBeStopped:   make(chan bool, c.MustBeStoppedChannelSize),
		subRoutines:     new(sync.WaitGroup),
		mustStop:        new(atomic.Bool),
		httpErrors:      make(chan error, c.HttpErrorsChannelSize),
		dbErrors:        &dbErrorsChannel,
		jsonRpcHandlers: js.NewMethodRepository(),
	}
	srv.mustStop.Store(false)

	err = srv.initJsonRpcHandlers()
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

	err = srv.initACMServiceClient()
	if err != nil {
		return nil, err
	}

	err = srv.initDiagnosticData()
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
	srv.startHttpServer()

	srv.subRoutines.Add(2)
	go srv.listenForHttpErrors()
	go srv.listenForDbErrors()

	return nil
}

func (srv *Server) Stop() (err error) {
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

	return nil
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

			time.Sleep(time.Second * DbReconnectCoolDownPeriodSec)
		}
	}

	log.Println(c.MsgDbNetworkErrorListenerHasStopped)
}

func (srv *Server) httpRouter(rw http.ResponseWriter, req *http.Request) {
	srv.jsonRpcHandlers.ServeHTTP(rw, req)
}

func (srv *Server) initACMServiceClient() (err error) {
	srv.acmServiceClient, err = ac.NewClient(
		srv.settings.ACMSettings.Host,
		srv.settings.ACMSettings.Port,
		srv.settings.ACMSettings.Path,
		srv.settings.ACMSettings.EnableSelfSignedCertificate,
	)
	if err != nil {
		return err
	}

	fmt.Print(c.MsgPingingAcmModule)

	var params = am.PingParams{}
	var result am.PingResult
	err = srv.acmServiceClient.MakeRequest(context.Background(), &result, ac.FuncPing, params)
	if err != nil {
		return err
	}

	if !result.OK {
		return errors.New(ErrACMModuleIsBroken)
	}

	fmt.Println(c.MsgOK)

	return nil
}

func (srv *Server) initDiagnosticData() (err error) {
	srv.diag = &cdd.DiagnosticData{}

	return nil
}
