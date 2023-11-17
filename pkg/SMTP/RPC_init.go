package smtp

import (
	sc "github.com/vault-thirteen/SimpleBB/pkg/SMTP/client"
	sm "github.com/vault-thirteen/SimpleBB/pkg/SMTP/models"
)

func (srv *Server) initJsonRpcHandlers() (err error) {
	// Ping.
	err = srv.jsonRpcHandlers.RegisterMethod(sc.FuncPing, PingHandler{Server: srv}, sm.PingParams{}, sm.PingResult{})
	if err != nil {
		return err
	}

	// Message.
	err = srv.jsonRpcHandlers.RegisterMethod(sc.FuncSendMessage, SendMessageHandler{Server: srv}, sm.SendMessageParams{}, sm.SendMessageResult{})
	if err != nil {
		return err
	}

	// Other.
	err = srv.jsonRpcHandlers.RegisterMethod(sc.FuncShowDiagnosticData, ShowDiagnosticDataHandler{Server: srv}, sm.ShowDiagnosticDataParams{}, sm.ShowDiagnosticDataResult{})
	if err != nil {
		return err
	}

	return nil
}
