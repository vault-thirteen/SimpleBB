package gwm

import (
	"fmt"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	gm "github.com/vault-thirteen/SimpleBB/pkg/GWM/models"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cn "github.com/vault-thirteen/SimpleBB/pkg/common/net"
)

// RPC functions.

func (srv *Server) blockIPAddress(p *gm.BlockIPAddressParams) (result *gm.BlockIPAddressResult, re *jrm1.RpcError) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	if !srv.settings.SystemSettings.IsFirewallUsed {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_FirewallIsDisabled, RpcErrorMsg_FirewallIsDisabled, nil)
	}

	// Check parameters.
	if len(p.UserIPA) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_IPAddressIsNotSet, RpcErrorMsg_IPAddressIsNotSet, nil)
	}

	if p.BlockTimeSec == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_BlockTimeIsNotSet, RpcErrorMsg_BlockTimeIsNotSet, nil)
	}

	userIPAB, err := cn.ParseIPA(p.UserIPA)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_IPAddressError, fmt.Sprintf(c.RpcErrorMsgF_IPAddressError, err.Error()), nil)
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

func (srv *Server) isIPAddressBlocked(p *gm.IsIPAddressBlockedParams) (result *gm.IsIPAddressBlockedResult, re *jrm1.RpcError) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	if !srv.settings.SystemSettings.IsFirewallUsed {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_FirewallIsDisabled, RpcErrorMsg_FirewallIsDisabled, nil)
	}

	// Check parameters.
	if len(p.UserIPA) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_IPAddressIsNotSet, RpcErrorMsg_IPAddressIsNotSet, nil)
	}

	userIPAB, err := cn.ParseIPA(p.UserIPA)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_IPAddressError, fmt.Sprintf(c.RpcErrorMsgF_IPAddressError, err.Error()), nil)
	}

	// Search for an existing record.
	var n int
	n, err = srv.dbo.CountBlocksByIPAddress(userIPAB)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	return &gm.IsIPAddressBlockedResult{IsBlocked: n > 0}, nil
}

func (srv *Server) showDiagnosticData() (result *gm.ShowDiagnosticDataResult, re *jrm1.RpcError) {
	result = &gm.ShowDiagnosticDataResult{}
	result.TotalRequestsCount, result.SuccessfulRequestsCount = srv.js.GetRequestsCount()

	return result, nil
}
