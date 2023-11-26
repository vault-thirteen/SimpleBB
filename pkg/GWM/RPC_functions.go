package acm

import (
	"fmt"
	"net"

	js "github.com/osamingo/jsonrpc/v2"
	gm "github.com/vault-thirteen/SimpleBB/pkg/GWM/models"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
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

	userIPAB := net.ParseIP(p.UserIPA)
	if userIPAB == nil {
		return nil, &js.Error{Code: c.RpcErrorCode_IPAddressError, Message: fmt.Sprintf(c.RpcErrorMsgF_IPAddressError, c.ErrNetParseIP)}
	}

	// Search for an existing record.
	var err error
	var n int
	n, err = srv.dbo.CountBlocksByIPAddress(userIPAB)
	if err != nil {
		return nil, c.DatabaseError(err, srv.dbErrors)
	}

	// Insert a new record when none is found, or update an existing block.
	if n == 0 {
		err = srv.dbo.InsertBlock(userIPAB, p.BlockTimeSec)
		if err != nil {
			return nil, c.DatabaseError(err, srv.dbErrors)
		}
	} else {
		err = srv.dbo.IncreaseBlockDuration(userIPAB, p.BlockTimeSec)
		if err != nil {
			return nil, c.DatabaseError(err, srv.dbErrors)
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

	userIPAB := net.ParseIP(p.UserIPA)
	if userIPAB == nil {
		return nil, &js.Error{Code: c.RpcErrorCode_IPAddressError, Message: fmt.Sprintf(c.RpcErrorMsgF_IPAddressError, c.ErrNetParseIP)}
	}

	// Search for an existing record.
	var err error
	var n int
	n, err = srv.dbo.CountBlocksByIPAddress(userIPAB)
	if err != nil {
		return nil, c.DatabaseError(err, srv.dbErrors)
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
