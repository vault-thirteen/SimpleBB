package smtp

import (
	"fmt"

	js "github.com/osamingo/jsonrpc/v2"
	sm "github.com/vault-thirteen/SimpleBB/pkg/SMTP/models"
)

// RPC functions.

func (srv *Server) sendEmailMessage(recipient, subject, message string) (result *sm.SendMessageResult, jerr *js.Error) {
	srv.mailerGuard.Lock()
	defer srv.mailerGuard.Unlock()

	err := srv.mailer.SendMail([]string{recipient}, subject, message)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_MailerError, Message: fmt.Sprintf(RpcErrorMsgF_MailerError, err.Error())}
	}

	result = &sm.SendMessageResult{}

	return result, nil
}

func (srv *Server) showDiagnosticData() (result *sm.ShowDiagnosticDataResult, jerr *js.Error) {
	result = &sm.ShowDiagnosticDataResult{
		TotalRequestsCount:      srv.diag.GetTotalRequestsCount(),
		SuccessfulRequestsCount: srv.diag.GetSuccessfulRequestsCount(),
	}

	return result, nil
}
