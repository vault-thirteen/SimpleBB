package acm

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	js "github.com/osamingo/jsonrpc/v2"
	"github.com/vault-thirteen/SimpleBB/pkg/GWM/dbo"
	gs "github.com/vault-thirteen/SimpleBB/pkg/GWM/settings"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cdd "github.com/vault-thirteen/SimpleBB/pkg/common/DiagnosticData"
)

type Server struct {
	// Settings.
	settings *gs.Settings

	// HTTP server for internal requests.
	listenDsnInt  string
	httpServerInt *http.Server

	// HTTPS server for external requests.
	listenDsnExt  string
	httpServerExt *http.Server

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

	// Diagnostic data.
	diag *cdd.DiagnosticData
}

func NewServer(stn *gs.Settings) (srv *Server, err error) {
	err = stn.Check()
	if err != nil {
		return nil, err
	}

	dbErrorsChannel := make(chan error, c.DbErrorsChannelSize)

	srv = &Server{
		settings:        stn,
		listenDsnInt:    net.JoinHostPort(stn.IntHttpSettings.Host, strconv.FormatUint(uint64(stn.IntHttpSettings.Port), 10)),
		listenDsnExt:    net.JoinHostPort(stn.ExtHttpsSettings.Host, strconv.FormatUint(uint64(stn.ExtHttpsSettings.Port), 10)),
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

	err = srv.initDiagnosticData()
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
	srv.startHttpServerInt()
	srv.startHttpServerExt()

	srv.subRoutines.Add(3)
	go srv.listenForHttpErrors()
	go srv.listenForDbErrors()
	go srv.runScheduler()

	return nil
}

func (srv *Server) Stop() (err error) {
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
			err = srv.clearIPAddresses()
			if err != nil {
				log.Println(err)
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
	srv.jsonRpcHandlers.ServeHTTP(rw, req)
}

// HTTP router for external requests.
func (srv *Server) httpRouterExt(rw http.ResponseWriter, req *http.Request) {
	//TODO
}

func (srv *Server) initDiagnosticData() (err error) {
	srv.diag = &cdd.DiagnosticData{}

	return nil
}
