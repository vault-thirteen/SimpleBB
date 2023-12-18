package models

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/vault-thirteen/auxie/number"
)

const (
	UrlParameter_ForumId    = "f"
	UrlParameter_ThreadId   = "t"
	UrlParameter_PageNumber = "p"
)

const (
	ErrRequestIsNotSet              = "request is not set"
	ErrUnsupportedParameter         = "unsupported parameter: %v"
	ErrDuplicateParameter           = "duplicate parameter: %v"
	ErrMultipleObjectTypesRequested = "multiple object types were requested"
)

type UrlParameter struct {
	ForumId    *uint
	ThreadId   *uint
	PageNumber *uint
}

func NewUrlParameterFromHttpRequest(req *http.Request) (up *UrlParameter, err error) {
	if req == nil {
		return nil, errors.New(ErrRequestIsNotSet)
	}

	up = &UrlParameter{}
	var tmpU uint

	for pn, pv := range req.URL.Query() {
		if len(pv) != 1 {
			return nil, fmt.Errorf(ErrDuplicateParameter, pn)
		}

		tmpU, err = number.ParseUint(pv[0])
		if err != nil {
			return nil, err
		}

		switch pn {
		case UrlParameter_ForumId:
			up.ForumId = &tmpU
		case UrlParameter_ThreadId:
			up.ThreadId = &tmpU
		case UrlParameter_PageNumber:
			up.PageNumber = &tmpU
		default:
			return nil, fmt.Errorf(ErrUnsupportedParameter, pn)
		}
	}

	if (up.ForumId != nil) && (up.ThreadId != nil) {
		return nil, errors.New(ErrMultipleObjectTypesRequested)
	}

	return up, nil
}
