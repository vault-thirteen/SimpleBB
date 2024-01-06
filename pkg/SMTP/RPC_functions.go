package smtp

import (
	"fmt"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	sm "github.com/vault-thirteen/SimpleBB/pkg/SMTP/models"
)

// RPC functions.

func (srv *Server) sendMessage(p *sm.SendMessageParams) (result *sm.SendMessageResult, re *jrm1.RpcError) {
	srv.mailerGuard.Lock()
	defer srv.mailerGuard.Unlock()

	err := srv.mailer.SendMail([]string{p.Recipient}, p.Subject, p.Message)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MailerError, fmt.Sprintf(RpcErrorMsgF_MailerError, err.Error()), err)
	}

	result = &sm.SendMessageResult{}

	return result, nil
}

func (srv *Server) showDiagnosticData() (result *sm.ShowDiagnosticDataResult, re *jrm1.RpcError) {
	result = &sm.ShowDiagnosticDataResult{}
	result.TotalRequestsCount, result.SuccessfulRequestsCount = srv.js.GetRequestsCount()

	return result, nil
}
