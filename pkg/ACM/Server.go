package acm

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
	"github.com/vault-thirteen/SimpleBB/pkg/ACM/dbo"
	"github.com/vault-thirteen/SimpleBB/pkg/ACM/im"
	"github.com/vault-thirteen/SimpleBB/pkg/ACM/km"
	as "github.com/vault-thirteen/SimpleBB/pkg/ACM/settings"
	gc "github.com/vault-thirteen/SimpleBB/pkg/GWM/client"
	gm "github.com/vault-thirteen/SimpleBB/pkg/GWM/models"
	rc "github.com/vault-thirteen/SimpleBB/pkg/RCS/client"
	rm "github.com/vault-thirteen/SimpleBB/pkg/RCS/models"
	sc "github.com/vault-thirteen/SimpleBB/pkg/SMTP/client"
	sm "github.com/vault-thirteen/SimpleBB/pkg/SMTP/models"
	"github.com/vault-thirteen/SimpleBB/pkg/avm"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cdd "github.com/vault-thirteen/SimpleBB/pkg/common/DiagnosticData"
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/client"
	rp "github.com/vault-thirteen/auxie/rpofs"
)

type Server struct {
	// Settings.
	settings *as.Settings

	// HTTPS server.
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

	// JSON-RPC handlers.
	jsonRpcHandlers *js.MethodRepository

	// Verification code generator.
	vcg *rp.Generator

	// SMTP Module client.
	smtpModuleClient *cc.Client

	// Captcha.
	captchaServiceClient *cc.Client

	// Generator of request IDs.
	ridg *rp.Generator

	// JWT key maker.
	jwtkm *km.KeyMaker

	// Diagnostic data.
	diag *cdd.DiagnosticData

	// Incident manager.
	incidentManager *im.IncidentManager

	// Gateway module client.
	gwmClient *cc.Client
}

func NewServer(stn *as.Settings) (srv *Server, err error) {
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
		ssp:             avm.NewSSP(),
		jsonRpcHandlers: js.NewMethodRepository(),
	}
	srv.mustStop.Store(false)

	err = srv.initJsonRpcHandlers()
	if err != nil {
		return nil, err
	}

	// Database.
	sp := dbo.SystemParameters{
		PreSessionExpirationTime: srv.settings.SystemSettings.PreSessionExpirationTime,
	}

	srv.dbo = dbo.NewDatabaseObject(srv.settings.DbSettings, sp)

	err = srv.dbo.Init()
	if err != nil {
		return nil, err
	}

	// HTTPS Server.
	srv.httpServer = &http.Server{
		Addr:    srv.listenDsn,
		Handler: http.Handler(http.HandlerFunc(srv.httpRouter)),
	}

	err = srv.initVerificationCodeGenerator()
	if err != nil {
		return nil, err
	}

	err = srv.initSmtpModuleClient()
	if err != nil {
		return nil, err
	}

	err = srv.initCaptchaServiceClient()
	if err != nil {
		return nil, err
	}

	err = srv.initRequestIdGenerator()
	if err != nil {
		return nil, err
	}

	err = srv.initJwtKeyMaker()
	if err != nil {
		return nil, err
	}

	err = srv.initDiagnosticData()
	if err != nil {
		return nil, err
	}

	if stn.SystemSettings.IsTableOfIncidentsUsed {
		fmt.Println(c.MsgIncidentsTableIsEnabled)

		err = srv.initGWModuleClient()
		if err != nil {
			return nil, err
		}
	} else {
		fmt.Println(c.MsgIncidentsTableIsDisabled)
	}

	err = srv.initIncidentManager()
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

	srv.subRoutines.Add(3)
	go srv.listenForHttpErrors()
	go srv.listenForDbErrors()
	go srv.runScheduler()

	err = srv.incidentManager.Start()
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

	err = srv.incidentManager.Stop()
	if err != nil {
		return err
	}

	err = srv.dbo.Fin()
	if err != nil {
		return err
	}

	srv.ssp.CompleteStop()

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
			err = srv.clearPreRegUsersTable()
			if err != nil {
				log.Println(err)
			}

			err = srv.clearPasswordChangesTable()
			if err != nil {
				log.Println(err)
			}

			err = srv.clearEmailChangesTable()
			if err != nil {
				log.Println(err)
			}

			err = srv.clearSessions()
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

func (srv *Server) httpRouter(rw http.ResponseWriter, req *http.Request) {
	srv.jsonRpcHandlers.ServeHTTP(rw, req)
}

func (srv *Server) initVerificationCodeGenerator() (err error) {
	symbols := c.MakeSymbolsNumbersAndCapitalLatinLetters()

	srv.vcg, err = rp.NewGenerator(srv.settings.SystemSettings.VerificationCodeLength, symbols)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) initSmtpModuleClient() (err error) {
	srv.smtpModuleClient, err = sc.NewClient(
		srv.settings.SmtpModuleSettings.Host,
		srv.settings.SmtpModuleSettings.Port,
		srv.settings.SmtpModuleSettings.Path,
	)
	if err != nil {
		return err
	}

	var params = sm.PingParams{}
	var result sm.PingResult

	fmt.Print(c.MsgPingingSmtpModule)

	err = srv.smtpModuleClient.MakeRequest(context.Background(), &result, sc.FuncPing, params)
	if err != nil {
		return err
	}

	if !result.OK {
		return errors.New(c.MsgSmtpModuleIsBroken)
	}

	fmt.Println(c.MsgOK)

	return nil
}

func (srv *Server) initCaptchaServiceClient() (err error) {
	srv.captchaServiceClient, err = rc.NewClient(
		srv.settings.CaptchaServiceSettings.RpcHost,
		srv.settings.CaptchaServiceSettings.RpcPort,
		srv.settings.CaptchaServiceSettings.RpcPath,
	)
	if err != nil {
		return err
	}

	var params = rm.PingParams{}
	var result rm.PingResult

	fmt.Print(c.MsgPingingCaptchaModule)

	err = srv.captchaServiceClient.MakeRequest(context.Background(), &result, rc.FuncPing, params)
	if err != nil {
		return err
	}

	if !result.OK {
		return errors.New(c.MsgCaptchaServiceIsBroken)
	}

	fmt.Println(c.MsgOK)

	return nil
}

func (srv *Server) initRequestIdGenerator() (err error) {
	symbols := c.MakeSymbolsNumbersAndCapitalLatinLetters()

	srv.ridg, err = rp.NewGenerator(srv.settings.SystemSettings.LogInRequestIdLength, symbols)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) initJwtKeyMaker() (err error) {
	srv.jwtkm, err = km.New(
		srv.settings.JWTSettings.SigningMethod,
		srv.settings.JWTSettings.PrivateKeyFilePath,
		srv.settings.JWTSettings.PublicKeyFilePath,
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

func (srv *Server) initGWModuleClient() (err error) {
	srv.gwmClient, err = sc.NewClient(
		srv.settings.GatewayModuleSettings.Host,
		srv.settings.GatewayModuleSettings.Port,
		srv.settings.GatewayModuleSettings.Path,
	)
	if err != nil {
		return err
	}

	var params = gm.PingParams{}
	var result gm.PingResult

	fmt.Print(c.MsgPingingGatewayModule)

	// While several services are inter-dependent, we make several attempts to
	// ping the Gateway module.
	for i := 1; i <= c.ServicePingAttemptsCount; i++ {
		err = srv.gwmClient.MakeRequest(context.Background(), &result, gc.FuncPing, params)
		if err == nil {
			break
		}

		fmt.Print(c.MsgPingAttempt)

		if i < c.ServicePingAttemptsCount {
			time.Sleep(time.Second * time.Duration(c.ServiceNextPingAttemptDelaySec))
		}
	}

	if err != nil {
		return err
	}

	if !result.OK {
		return errors.New(c.MsgGatewayModuleIsBroken)
	}

	fmt.Println(c.MsgOK)

	return nil
}

func (srv *Server) initIncidentManager() (err error) {
	srv.incidentManager = im.NewIncidentManager(srv.settings.SystemSettings.IsTableOfIncidentsUsed, srv.dbo, srv.gwmClient, &srv.settings.SystemSettings.BlockTimePerIncident)

	return nil
}
