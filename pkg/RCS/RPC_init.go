package rcs

import (
	rc "github.com/vault-thirteen/SimpleBB/pkg/RCS/client"
	rm "github.com/vault-thirteen/SimpleBB/pkg/RCS/models"
)

func (srv *Server) initJsonRpcHandlers() (err error) {
	// Ping.
	err = srv.jsonRpcHandlers.RegisterMethod(rc.FuncPing, PingHandler{Server: srv}, rm.PingParams{}, rm.PingResult{})
	if err != nil {
		return err
	}

	// Captcha.
	err = srv.jsonRpcHandlers.RegisterMethod(rc.FuncCreateCaptcha, CreateCaptchaHandler{Server: srv}, rm.CreateCaptchaParams{}, rm.CreateCaptchaResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(rc.FuncCheckCaptcha, CheckCaptchaHandler{Server: srv}, rm.CheckCaptchaParams{}, rm.CheckCaptchaResult{})
	if err != nil {
		return err
	}

	// Other.
	err = srv.jsonRpcHandlers.RegisterMethod(rc.FuncShowDiagnosticData, ShowDiagnosticDataHandler{Server: srv}, rm.ShowDiagnosticDataParams{}, rm.ShowDiagnosticDataResult{})
	if err != nil {
		return err
	}

	return nil
}
