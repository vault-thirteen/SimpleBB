package http

import (
	"fmt"
	"net/http"

	mime "github.com/vault-thirteen/auxie/MIME"
)

const (
	ContentType_TextPlain       = mime.TypeTextPlain
	ContentType_ApplicationJson = mime.TypeApplicationJson
)

const (
	ErrFMultipleHeaders  = "multiple headers: %s"
	ErrFHeaderIsNotFound = "header is not found: %s"
)

// GetSingleHeader reads exactly one, single HTTP header. If multiple headers
// are found, an error is returned. Unfortunately, Go language does not do this
// by default and returns only a first header value even when there are
// multiple values available. Such a behaviour may lead to unexpected errors.
func GetSingleHeader(req *http.Request, headerName string) (h string, err error) {
	headers := req.Header[headerName]

	if len(headers) > 1 {
		return "", fmt.Errorf(ErrFMultipleHeaders, headerName)
	}

	if len(headers) == 0 {
		return "", fmt.Errorf(ErrFHeaderIsNotFound, headerName)
	}

	return headers[0], nil
}
