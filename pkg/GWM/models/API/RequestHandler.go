package api

import (
	"net/http"
)

type RequestHandler = func(*Request, http.ResponseWriter, *http.Request)
