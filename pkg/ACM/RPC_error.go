package acm

import (
	"net/http"
)

// RPC errors.

// Error codes must not exceed 999.

// Codes.
const (
	RpcErrorCode_Exception                          = 1
	RpcErrorCode_AuthChallengeResponseIsNotSet      = 2
	RpcErrorCode_BPP_GenerateRandomSalt             = 3
	RpcErrorCode_BPP_PackSymbols                    = 4
	RpcErrorCode_CaptchaAnswerIsNotSet              = 5
	RpcErrorCode_CaptchaAnswerIsWrong               = 6
	RpcErrorCode_EmailAddressIsNotValid             = 7
	RpcErrorCode_EmailAddressIsUsed                 = 8
	RpcErrorCode_EmailChangeIsNotFound              = 9
	RpcErrorCode_JWTCreation                        = 10
	RpcErrorCode_NameIsNotSet                       = 11
	RpcErrorCode_NameIsTooLong                      = 12
	RpcErrorCode_NameIsUsed                         = 13
	RpcErrorCode_NewPasswordIsNotSet                = 14
	RpcErrorCode_NewEmailIsNotSet                   = 15
	RpcErrorCode_PasswordChangeIsNotFound           = 16
	RpcErrorCode_PasswordCheckError                 = 17
	RpcErrorCode_PasswordIsNotValid                 = 18
	RpcErrorCode_PasswordIsNotSet                   = 19
	RpcErrorCode_PasswordIsTooLong                  = 20
	RpcErrorCode_PasswordIsWrong                    = 21
	RpcErrorCode_RCS_CheckCaptcha                   = 22
	RpcErrorCode_RCS_CreateCaptcha                  = 23
	RpcErrorCode_RequestIdGenerator                 = 24
	RpcErrorCode_RequestIdIsNotSet                  = 25
	RpcErrorCode_SmtpModule                         = 26
	RpcErrorCode_StepIsUnknown                      = 27
	RpcErrorCode_ThisActionIsNotForLoggedUsers      = 28
	RpcErrorCode_UserAlreadyStartedToChangePassword = 29
	RpcErrorCode_UserAlreadyStartedToChangeEmail    = 30
	RpcErrorCode_UserCanNotLogIn                    = 31
	RpcErrorCode_UserHasAlreadyStartedToLogIn       = 32
	RpcErrorCode_UserHasNotStartedToLogIn           = 33
	RpcErrorCode_UserIdIsNotSet                     = 34
	RpcErrorCode_UserIsAlreadyLoggedIn              = 35
	RpcErrorCode_UserIsNotFound                     = 36
	RpcErrorCode_UserPreSessionIsNotFound           = 37
	RpcErrorCode_VerificationCodeGenerator          = 38
	RpcErrorCode_VerificationCodeIsNotSet           = 39
	RpcErrorCode_VerificationCodeIsWrong            = 40
)

// Messages.
const (
	RpcErrorMsg_Exception                          = "exception"
	RpcErrorMsg_AuthChallengeResponseIsNotSet      = "authorisation challenge response is not set"
	RpcErrorMsg_BPP_GenerateRandomSalt             = "BPP: GenerateRandomSalt error"
	RpcErrorMsgF_BPP_PackSymbols                   = "BPP: PackSymbols error: %s" // Template.
	RpcErrorMsg_CaptchaAnswerIsNotSet              = "captcha answer is not set"
	RpcErrorMsg_CaptchaAnswerIsWrong               = "captcha answer is wrong"
	RpcErrorMsg_EmailAddressIsNotValid             = "e-mail address is not valid"
	RpcErrorMsg_EmailAddressIsUsed                 = "e-mail address is used"
	RpcErrorMsg_EmailChangeIsNotFound              = "request for e-mail address change is not found"
	RpcErrorMsgF_JWTCreation                       = "JWT creation error: %s" // Template.
	RpcErrorMsg_NameIsNotSet                       = "name is not set"
	RpcErrorMsg_NameIsTooLong                      = "name is too long"
	RpcErrorMsg_NameIsUsed                         = "name is already used"
	RpcErrorMsg_NewPasswordIsNotSet                = "new password is not set"
	RpcErrorMsg_NewEmailIsNotSet                   = "new e-mail address is not set"
	RpcErrorMsg_PasswordChangeIsNotFound           = "request for password change is not found"
	RpcErrorMsgF_PasswordCheckError                = "password check error: %s"  // Template.
	RpcErrorMsgF_PasswordIsNotValid                = "password is not valid: %s" // Template.
	RpcErrorMsg_PasswordIsNotSet                   = "password is not set"
	RpcErrorMsg_PasswordIsTooLong                  = "password is too long"
	RpcErrorMsg_PasswordIsWrong                    = "password is wrong"
	RpcErrorMsgF_RCS_CheckCaptcha                  = "RCS: CheckCaptcha error: %s"  // Template.
	RpcErrorMsgF_RCS_CreateCaptcha                 = "RCS: CreateCaptcha error: %s" // Template.
	RpcErrorMsg_RequestIdGenerator                 = "error generating request ID"
	RpcErrorMsg_RequestIdIsNotSet                  = "request ID is not set"
	RpcErrorMsgF_SmtpModule                        = "SMTP module error: %s" // Template.
	RpcErrorMsg_StepIsUnknown                      = "unknown step"
	RpcErrorMsg_ThisActionIsNotForLoggedUsers      = "this action is not for logged in users"
	RpcErrorMsg_UserAlreadyStartedToChangePassword = "user has already started to change password"
	RpcErrorMsg_UserAlreadyStartedToChangeEmail    = "user has already started to change e-mail address"
	RpcErrorMsg_UserCanNotLogIn                    = "user can not log in"
	RpcErrorMsg_UserHasAlreadyStartedToLogIn       = "user has already started to log in"
	RpcErrorMsg_UserHasNotStartedToLogIn           = "user has not started to log in"
	RpcErrorMsg_UserIdIsNotSet                     = "user ID is not set"
	RpcErrorMsg_UserIsAlreadyLoggedIn              = "user is already logged in"
	RpcErrorMsg_UserIsNotFound                     = "user is not found"
	RpcErrorMsg_UserPreSessionIsNotFound           = "user's preliminary session is not found"
	RpcErrorMsg_VerificationCodeGenerator          = "verification code generator error"
	RpcErrorMsg_VerificationCodeIsNotSet           = "verification code is not set"
	RpcErrorMsg_VerificationCodeIsWrong            = "verification code is wrong"
)

// Unique HTTP status codes used in the map:
// - 400 (Bad request);
// - 403 (Forbidden);
// - 404 (Not found);
// - 409 (Conflict);
// - 500 (Internal server error).
func GetMapOfHttpStatusCodesByRpcErrorCodes() map[int]int {
	return map[int]int{
		RpcErrorCode_AuthChallengeResponseIsNotSet:      http.StatusBadRequest,
		RpcErrorCode_BPP_GenerateRandomSalt:             http.StatusInternalServerError,
		RpcErrorCode_BPP_PackSymbols:                    http.StatusInternalServerError,
		RpcErrorCode_CaptchaAnswerIsNotSet:              http.StatusBadRequest,
		RpcErrorCode_CaptchaAnswerIsWrong:               http.StatusForbidden,
		RpcErrorCode_EmailAddressIsNotValid:             http.StatusBadRequest,
		RpcErrorCode_EmailAddressIsUsed:                 http.StatusConflict,
		RpcErrorCode_EmailChangeIsNotFound:              http.StatusNotFound,
		RpcErrorCode_JWTCreation:                        http.StatusInternalServerError,
		RpcErrorCode_NameIsNotSet:                       http.StatusBadRequest,
		RpcErrorCode_NameIsTooLong:                      http.StatusBadRequest,
		RpcErrorCode_NameIsUsed:                         http.StatusConflict,
		RpcErrorCode_NewPasswordIsNotSet:                http.StatusBadRequest,
		RpcErrorCode_NewEmailIsNotSet:                   http.StatusBadRequest,
		RpcErrorCode_PasswordChangeIsNotFound:           http.StatusNotFound,
		RpcErrorCode_PasswordCheckError:                 http.StatusInternalServerError,
		RpcErrorCode_PasswordIsNotValid:                 http.StatusBadRequest,
		RpcErrorCode_PasswordIsNotSet:                   http.StatusBadRequest,
		RpcErrorCode_PasswordIsTooLong:                  http.StatusBadRequest,
		RpcErrorCode_PasswordIsWrong:                    http.StatusForbidden,
		RpcErrorCode_RCS_CheckCaptcha:                   http.StatusInternalServerError,
		RpcErrorCode_RCS_CreateCaptcha:                  http.StatusInternalServerError,
		RpcErrorCode_RequestIdGenerator:                 http.StatusInternalServerError,
		RpcErrorCode_RequestIdIsNotSet:                  http.StatusBadRequest,
		RpcErrorCode_SmtpModule:                         http.StatusInternalServerError,
		RpcErrorCode_StepIsUnknown:                      http.StatusBadRequest,
		RpcErrorCode_ThisActionIsNotForLoggedUsers:      http.StatusForbidden,
		RpcErrorCode_UserAlreadyStartedToChangePassword: http.StatusForbidden,
		RpcErrorCode_UserAlreadyStartedToChangeEmail:    http.StatusForbidden,
		RpcErrorCode_UserCanNotLogIn:                    http.StatusForbidden,
		RpcErrorCode_UserHasAlreadyStartedToLogIn:       http.StatusForbidden,
		RpcErrorCode_UserHasNotStartedToLogIn:           http.StatusForbidden,
		RpcErrorCode_UserIdIsNotSet:                     http.StatusBadRequest,
		RpcErrorCode_UserIsAlreadyLoggedIn:              http.StatusForbidden,
		RpcErrorCode_UserIsNotFound:                     http.StatusNotFound,
		RpcErrorCode_UserPreSessionIsNotFound:           http.StatusNotFound,
		RpcErrorCode_VerificationCodeGenerator:          http.StatusInternalServerError,
		RpcErrorCode_VerificationCodeIsNotSet:           http.StatusBadRequest,
		RpcErrorCode_VerificationCodeIsWrong:            http.StatusForbidden,
	}
}
