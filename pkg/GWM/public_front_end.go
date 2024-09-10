package gwm

import (
	"log"
	"net/http"

	"github.com/vault-thirteen/SimpleBB/pkg/GWM/models"
	"github.com/vault-thirteen/auxie/header"
)

func (srv *Server) handleFrontEndStaticFile(rw http.ResponseWriter, req *http.Request, fedf models.FrontEndFileData) {
	if req.Method != http.MethodGet {
		srv.respondMethodNotAllowed(rw)
		return
	}

	rw.Header().Set(header.HttpHeaderContentType, fedf.ContentType)

	_, err := rw.Write(fedf.CachedFile)
	if err != nil {
		log.Println(err.Error())
		return
	}
}
