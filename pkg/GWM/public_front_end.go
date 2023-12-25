package gwm

import (
	"log"
	"net/http"

	"github.com/vault-thirteen/SimpleBB/pkg/GWM/models"
)

func (srv *Server) handlePublicFrontEnd(rw http.ResponseWriter, req *http.Request) {
	up, err := models.NewUrlParameterFromHttpRequest(req)
	if err != nil {
		srv.processBadRequest(rw)
		return
	}

	//TODO
	rw.WriteHeader(http.StatusTeapot)
	if up.ForumId != nil {
		// Showing a forum.
		_, err = rw.Write([]byte("Forum"))
	} else if up.ThreadId != nil {
		// Showing a thread.
		_, err = rw.Write([]byte("Thread"))
	} else {
		// Showing sections & forums.
		_, err = rw.Write([]byte("Sections & Forums"))
	}
	if err != nil {
		log.Println(err)
	}
}
