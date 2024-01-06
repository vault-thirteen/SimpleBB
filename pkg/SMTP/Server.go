package smtp

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	ss "github.com/vault-thirteen/SimpleBB/pkg/SMTP/settings"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	"github.com/vault-thirteen/SimpleBB/pkg/common/avm"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

type Server struct {
	settings *ss.Settings

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
	ssp         *avm.SSP

	// Mailer.
	mailer      *Mailer
	mailerGuard sync.Mutex

	// JSON-RPC server.
	js *jrm1.Processor
}

func NewServer(s cm.ISettings) (srv *Server, err error) {
	stn := s.(*ss.Settings)

	err = stn.Check()
	if err != nil {
		return nil, err
	}

	srv = &Server{
		settings:      stn,
		listenDsn:     net.JoinHostPort(stn.HttpSettings.Host, strconv.FormatUint(uint64(stn.HttpSettings.Port), 10)),
		mustBeStopped: make(chan bool, c.MustBeStoppedChannelSize),
		subRoutines:   new(sync.WaitGroup),
		mustStop:      new(atomic.Bool),
		httpErrors:    make(chan error, c.HttpErrorsChannelSize),
		ssp:           avm.NewSSP(),
	}
	srv.mustStop.Store(false)

	// RPC server.
	err = srv.initRpc()
	if err != nil {
		return nil, err
	}

	err = srv.initMailer()
	if err != nil {
		return nil, err
	}

	// HTTP Server.
	srv.httpServer = &http.Server{
		Addr:    srv.listenDsn,
		Handler: http.Handler(http.HandlerFunc(srv.httpRouter)),
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

	srv.subRoutines.Add(1)
	go srv.listenForHttpErrors()

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

	srv.subRoutines.Wait()

	srv.ssp.CompleteStop()

	return nil
}

func (srv *Server) startHttpServer() {
	go func() {
		var listenError error
		listenError = srv.httpServer.ListenAndServe()

		if (listenError != nil) && (listenError != http.ErrServerClosed) {
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

func (srv *Server) httpRouter(rw http.ResponseWriter, req *http.Request) {
	srv.js.ServeHTTP(rw, req)
}

func (srv *Server) initMailer() (err error) {
	srv.mailer, err = NewMailer(
		srv.settings.SmtpSettings.Host,
		srv.settings.SmtpSettings.Port,
		srv.settings.SmtpSettings.User,
		srv.settings.SmtpSettings.Password,
		srv.settings.SmtpSettings.UserAgent,
	)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) ReportStart() {
	fmt.Println(c.MsgHttpServer + srv.GetListenDsn())
}

func (srv *Server) UseConstructor(stn cm.ISettings) (cm.IServer, error) {
	return NewServer(stn)
}
