package gwm

import (
	"fmt"
	"log"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
)

// Auxiliary functions used in RPC functions.

// logError logs error if debug mode is enabled.
func (srv *Server) logError(err error) {
	if err == nil {
		return
	}

	if srv.settings.SystemSettings.IsDebugMode {
		log.Println(err)
	}
}

// databaseError processes the database error and returns it.
func (srv *Server) databaseError(err error) (re *jrm1.RpcError) {
	if err == nil {
		return nil
	}

	if c.IsNetworkError(err) {
		log.Println(fmt.Sprintf(c.ErrFDatabaseNetwork, err.Error()))
		*(srv.dbErrors) <- err
	} else {
		srv.logError(err)
	}

	return jrm1.NewRpcErrorByUser(c.RpcErrorCode_DatabaseError, c.RpcErrorMsg_DatabaseError, err)
}
