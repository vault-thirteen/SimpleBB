package acm

import (
	gc "github.com/vault-thirteen/SimpleBB/pkg/GWM/client"
	gm "github.com/vault-thirteen/SimpleBB/pkg/GWM/models"
)

func (srv *Server) initJsonRpcHandlers() (err error) {
	// Ping.
	err = srv.jsonRpcHandlers.RegisterMethod(gc.FuncPing, PingHandler{Server: srv}, gm.PingParams{}, gm.PingResult{})
	if err != nil {
		return err
	}

	// IP address list.
	err = srv.jsonRpcHandlers.RegisterMethod(gc.FuncBlockIPAddress, BlockIPAddressHandler{Server: srv}, gm.BlockIPAddressParams{}, gm.BlockIPAddressResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(gc.FuncIsIPAddressBlocked, IsIPAddressBlockedHandler{Server: srv}, gm.IsIPAddressBlockedParams{}, gm.IsIPAddressBlockedResult{})
	if err != nil {
		return err
	}

	// Other.
	err = srv.jsonRpcHandlers.RegisterMethod(gc.FuncShowDiagnosticData, ShowDiagnosticDataHandler{Server: srv}, gm.ShowDiagnosticDataParams{}, gm.ShowDiagnosticDataResult{})
	if err != nil {
		return err
	}

	return nil
}
