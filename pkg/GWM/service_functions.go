package gwm

import (
	"net/http"

	gm "github.com/vault-thirteen/SimpleBB/pkg/GWM/models"
)

// Service functions.

func (srv *Server) GetProductVersion(ar *gm.ApiRequest, rw http.ResponseWriter, req *http.Request) {
	srv.respondWithPlainText(rw, srv.settings.VersionInfo.ProgramVersionString())
}
