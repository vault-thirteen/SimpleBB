package acm

// RPC errors.

// Error codes must not exceed 999.

// Codes.
const (
	RpcErrorCode_AuthChallengeResponseIsNotSet      = 1
	RpcErrorCode_BPP_GenerateRandomSalt             = 2
	RpcErrorCode_BPP_PackSymbols                    = 3
	RpcErrorCode_CaptchaAnswerIsNotSet              = 4
	RpcErrorCode_CaptchaAnswerIsWrong               = 5
	RpcErrorCode_EmailAddressIsNotValid             = 6
	RpcErrorCode_EmailAddressIsUsed                 = 7
	RpcErrorCode_EmailChangeIsNotFound              = 8
	RpcErrorCode_JWTCreation                        = 9
	RpcErrorCode_NameIsNotSet                       = 10
	RpcErrorCode_NameIsTooLong                      = 11
	RpcErrorCode_NameIsUsed                         = 12
	RpcErrorCode_NewPasswordIsNotSet                = 13
	RpcErrorCode_NewEmailIsNotSet                   = 14
	RpcErrorCode_PasswordChangeIsNotFound           = 15
	RpcErrorCode_PasswordCheckError                 = 16
	RpcErrorCode_PasswordError                      = 17
	RpcErrorCode_PasswordIsNotSet                   = 18
	RpcErrorCode_PasswordIsTooLong                  = 19
	RpcErrorCode_PasswordIsWrong                    = 20
	RpcErrorCode_RCS_CheckCaptcha                   = 21
	RpcErrorCode_RCS_CreateCaptcha                  = 22
	RpcErrorCode_RequestIdGenerator                 = 23
	RpcErrorCode_RequestIdIsNotSet                  = 24
	RpcErrorCode_SmtpModule                         = 25
	RpcErrorCode_StepIsUnknown                      = 26
	RpcErrorCode_ThisActionIsNotForLoggedUsers      = 27
	RpcErrorCode_UserAlreadyStartedToChangePassword = 28
	RpcErrorCode_UserAlreadyStartedToChangeEmail    = 29
	RpcErrorCode_UserCanNotLogIn                    = 30
	RpcErrorCode_UserHasAlreadyStartedToLogIn       = 31
	RpcErrorCode_UserHasNotStartedToLogIn           = 32
	RpcErrorCode_UserIdIsNotSet                     = 33
	RpcErrorCode_UserIsAlreadyLoggedIn              = 34
	RpcErrorCode_UserIsNotFound                     = 35
	RpcErrorCode_UserPreSessionIsNotFound           = 36
	RpcErrorCode_VerificationCodeGenerator          = 37
	RpcErrorCode_VerificationCodeIsNotSet           = 38
	RpcErrorCode_VerificationCodeIsWrong            = 39
)

// Messages.
const (
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
	RpcErrorMsgF_PasswordCheckError                = "password check error: %s" // Template.
	RpcErrorMsgF_PasswordError                     = "password error: %s"       // Template.
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
