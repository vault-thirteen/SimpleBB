package common

import (
	"context"

	js "github.com/osamingo/jsonrpc/v2"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
)

const (
	ErrNetParseIP = "net.ParseIP error" // Golang's built-in parser does not return an error. What a shame again.
)

func GetRequestId(c context.Context) (rid am.RequestId, jerr *js.Error) {
	jerr = js.Unmarshal(js.RequestID(c), &rid)
	if jerr != nil {
		return rid, jerr
	}

	return rid, nil
}
