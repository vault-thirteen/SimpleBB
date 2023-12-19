package gwm

import (
	"fmt"
	"log"

	js "github.com/osamingo/jsonrpc/v2"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
)

// Auxiliary functions used in RPC functions.

// logError logs error if debug mode is enabled.
func (srv *Server) logError(err error) {
	if srv.settings.SystemSettings.IsDebugMode {
		log.Println(err)
	}
}

// databaseError processes the database error and returns a JSON RPC error.
func (srv *Server) databaseError(err error) (jerr *js.Error) {
	if c.IsNetworkError(err) {
		log.Println(fmt.Sprintf(c.ErrFDatabaseNetwork, err.Error()))
		*(srv.dbErrors) <- err
	} else {
		srv.logError(err)
	}

	return &js.Error{Code: c.RpcErrorCode_DatabaseError, Message: c.RpcErrorMsg_DatabaseError}
}
