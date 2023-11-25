package common

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	js "github.com/osamingo/jsonrpc/v2"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
)

// Auxiliary functions used in RPC functions.

const (
	EmailAndVerificationCodeSeparator = ":"
)

const (
	ErrNetParseIP                      = "net.ParseIP error" // Golang's built-in parser does not return an error. What a shame again.
	ErrDecodeEmailWithVerificationCode = "e-mail with verification code decoding error"
	ErrCaptchaAnswerIsNotSet           = "answer is not set"
	ErrFPasswordIsNotAllowed           = "password is not allowed: %s" // Template.
	ErrUnsupportedRCSMode              = "unsupported RCS mode"
	ErrFakeJWToken                     = "fake JWT token"
	ErrFMakeVerificationLink           = "makeVerificationLink: %v. E-mail: %s, VC: %s."
)

func GetRequestId(c context.Context) (rid am.RequestId, jerr *js.Error) {
	jerr = js.Unmarshal(js.RequestID(c), &rid)
	if jerr != nil {
		return rid, jerr
	}

	return rid, nil
}

// DatabaseError checks the database error and returns the JSON RPC error.
func DatabaseError(err error, dbErrorsChan *chan error) (jerr *js.Error) {
	CheckForNetworkError(err, dbErrorsChan)
	log.Println(fmt.Sprintf(MsgFDBError, err.Error()))
	return &js.Error{Code: RpcErrorCode_DatabaseError, Message: RpcErrorMsg_DatabaseError}
}

// CheckForNetworkError informs the system about database network errors.
func CheckForNetworkError(err error, dbErrorsChan *chan error) {
	var nerr net.Error
	isNetworkError := errors.As(err, &nerr)

	if isNetworkError {
		*dbErrorsChan <- nerr
	}
}
