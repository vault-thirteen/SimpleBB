package models

import "net/http"

type ApiRequestHandler = func(*ApiRequest, http.ResponseWriter, *http.Request)
