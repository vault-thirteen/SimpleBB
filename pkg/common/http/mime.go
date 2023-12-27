package http

import (
	"net/http"

	mime "github.com/vault-thirteen/auxie/MIME"
	"github.com/vault-thirteen/auxie/header"
	hh "github.com/vault-thirteen/auxie/http-helper"
)

func CheckClientSupportForJson(req *http.Request) (ok bool, err error) {
	var httpAcceptHeader string
	httpAcceptHeader, err = GetSingleHeader(req, header.HttpHeaderAccept)
	if err != nil {
		return false, err
	}

	var amts *hh.AcceptedMimeTypes
	amts, err = hh.NewAcceptedMimeTypesFromHeader(httpAcceptHeader)
	if err != nil {
		return false, err
	}

	var amt *hh.AcceptedMimeType
	for {
		amt, err = amts.Next()
		if err != nil {
			break
		}

		switch amt.MimeType {
		case mime.TypeApplicationJson,
			mime.TypeApplicationAny,
			mime.TypeAny:
			return true, nil
		}
	}

	return false, nil
}
