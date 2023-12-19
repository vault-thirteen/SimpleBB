package common

import (
	"context"
	"errors"
	"net"

	js "github.com/osamingo/jsonrpc/v2"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
)

// Auxiliary functions used in RPC functions.

const (
	ErrCaptchaAnswerIsNotSet = "answer is not set"
	ErrFPasswordIsNotAllowed = "password is not allowed: %s" // Template.
	ErrUnsupportedRCSMode    = "unsupported RCS mode"
	ErrFakeJWToken           = "fake JWT token"
	ErrFDatabaseNetwork      = "database network error: %v" // Template.
)

func GetRequestId(c context.Context) (rid am.RequestId, jerr *js.Error) {
	jerr = js.Unmarshal(js.RequestID(c), &rid)
	if jerr != nil {
		return rid, jerr
	}

	return rid, nil
}

// IsNetworkError checks if an error is a network error.
func IsNetworkError(err error) (isNetworkError bool) {
	var nerr net.Error
	return errors.As(err, &nerr)
}
