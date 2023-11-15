package smtp

import (
	"context"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	js "github.com/osamingo/jsonrpc/v2"
	ss "github.com/vault-thirteen/SimpleBB/pkg/SMTP/settings"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cdd "github.com/vault-thirteen/SimpleBB/pkg/common/DiagnosticData"
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

	// Mailer.
	mailer      *Mailer
	mailerGuard sync.Mutex

	// JSON-RPC handlers.
	jsonRpcHandlers *js.MethodRepository

	// Diagnostic data.
	diag *cdd.DiagnosticData
}

func NewServer(stn *ss.Settings) (srv *Server, err error) {
	err = stn.Check()
	if err != nil {
		return nil, err
	}

	srv = &Server{
		settings:        stn,
		listenDsn:       net.JoinHostPort(stn.HttpSettings.Host, strconv.FormatUint(uint64(stn.HttpSettings.Port), 10)),
		mustBeStopped:   make(chan bool, 2),
		subRoutines:     new(sync.WaitGroup),
		mustStop:        new(atomic.Bool),
		httpErrors:      make(chan error, 8),
		jsonRpcHandlers: js.NewMethodRepository(),
	}
	srv.mustStop.Store(false)

	err = srv.initJsonRpcHandlers()
	if err != nil {
		return nil, err
	}

	err = srv.initMailer()
	if err != nil {
		return nil, err
	}

	err = srv.initDiagnosticData()
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
	srv.startHttpServer()

	srv.subRoutines.Add(1)
	go srv.listenForHttpErrors()

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

	srv.subRoutines.Wait()

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
	srv.jsonRpcHandlers.ServeHTTP(rw, req)
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

func (srv *Server) initDiagnosticData() (err error) {
	srv.diag = &cdd.DiagnosticData{}

	return nil
}
