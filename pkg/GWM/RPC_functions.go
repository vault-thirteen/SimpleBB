package gwm

import (
	"fmt"

	js "github.com/osamingo/jsonrpc/v2"
	gm "github.com/vault-thirteen/SimpleBB/pkg/GWM/models"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cn "github.com/vault-thirteen/SimpleBB/pkg/common/net"
)

// RPC functions.

func (srv *Server) blockIPAddress(p *gm.BlockIPAddressParams) (result *gm.BlockIPAddressResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	if !srv.settings.SystemSettings.IsFirewallUsed {
		return nil, &js.Error{Code: RpcErrorCode_FirewallIsDisabled, Message: RpcErrorMsg_FirewallIsDisabled}
	}

	// Check parameters.
	if len(p.UserIPA) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_IPAddressIsNotSet, Message: RpcErrorMsg_IPAddressIsNotSet}
	}

	if p.BlockTimeSec == 0 {
		return nil, &js.Error{Code: RpcErrorCode_BlockTimeIsNotSet, Message: RpcErrorMsg_BlockTimeIsNotSet}
	}

	userIPAB, err := cn.ParseIPA(p.UserIPA)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: c.RpcErrorCode_IPAddressError, Message: fmt.Sprintf(c.RpcErrorMsgF_IPAddressError, err.Error())}
	}

	// Search for an existing record.
	var n int
	n, err = srv.dbo.CountBlocksByIPAddress(userIPAB)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Insert a new record when none is found, or update an existing block.
	if n == 0 {
		err = srv.dbo.InsertBlock(userIPAB, p.BlockTimeSec)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	} else {
		err = srv.dbo.IncreaseBlockDuration(userIPAB, p.BlockTimeSec)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	return &gm.BlockIPAddressResult{OK: true}, nil
}

func (srv *Server) isIPAddressBlocked(p *gm.IsIPAddressBlockedParams) (result *gm.IsIPAddressBlockedResult, jerr *js.Error) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	if !srv.settings.SystemSettings.IsFirewallUsed {
		return nil, &js.Error{Code: RpcErrorCode_FirewallIsDisabled, Message: RpcErrorMsg_FirewallIsDisabled}
	}

	// Check parameters.
	if len(p.UserIPA) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_IPAddressIsNotSet, Message: RpcErrorMsg_IPAddressIsNotSet}
	}

	userIPAB, err := cn.ParseIPA(p.UserIPA)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: c.RpcErrorCode_IPAddressError, Message: fmt.Sprintf(c.RpcErrorMsgF_IPAddressError, err.Error())}
	}

	// Search for an existing record.
	var n int
	n, err = srv.dbo.CountBlocksByIPAddress(userIPAB)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	return &gm.IsIPAddressBlockedResult{IsBlocked: n > 0}, nil
}

func (srv *Server) showDiagnosticData() (result *gm.ShowDiagnosticDataResult, jerr *js.Error) {
	result = &gm.ShowDiagnosticDataResult{
		TotalRequestsCount:      srv.diag.GetTotalRequestsCount(),
		SuccessfulRequestsCount: srv.diag.GetSuccessfulRequestsCount(),
	}

	return result, nil
}
