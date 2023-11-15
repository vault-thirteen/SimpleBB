package acm

// RPC errors.

// Error codes must not exceed 999.

// Codes.
const (
	RpcErrorCode_RegisterUser_StepIsUnknown              = 3
	RpcErrorCode_RegisterUser_EmailAddresIsNotValid      = 4
	RpcErrorCode_RegisterUser_EmailAddresIsUsed          = 5
	RpcErrorCode_RegisterUser_VerificationCodeGenerator  = 6
	RpcErrorCode_RegisterUser_SmtpModule                 = 7
	RpcErrorCode_RegisterUser_VerificationCodeIsNotSet   = 8
	RpcErrorCode_RegisterUser_VerificationCodeIsNotValid = 9
	RpcErrorCode_RegisterUser_NameIsNotSet               = 10
	RpcErrorCode_RegisterUser_PasswordIsNotSet           = 11
	RpcErrorCode_RegisterUser_UserNameIsTooLong          = 12
	RpcErrorCode_RegisterUser_UserNameIsUsed             = 13
	RpcErrorCode_RegisterUser_PasswordError              = 14
	RpcErrorCode_BPP_PackSymbols                         = 15
	RpcErrorCode_RegisterUser_UserPasswordIsTooLong      = 16
	RpcErrorCode_LogUserIn_StepIsUnknown                 = 17
	RpcErrorCode_LogUserIn_UserCanNotLogIn               = 18
	RpcErrorCode_LogUserIn_UserIsAlreadyLoggedIn         = 19
	RpcErrorCode_LogUserIn_UserHasAlreadyStartedToLogIn  = 20
	RpcErrorCode_BPP_GenerateRandomSalt                  = 21
	RpcErrorCode_RCS_CreateCaptcha                       = 22
	RpcErrorCode_LogUserIn_RequestIdGenerator            = 23
	RpcErrorCode_LogUserIn_UserHasNotStartedToLogIn      = 24
	RpcErrorCode_LogUserIn_UserPreSessionIsNotFound      = 25
	RpcErrorCode_RCS_CheckCaptcha                        = 26
	RpcErrorCode_LogUserIn_CaptchaAnswerIsWrong          = 27
	RpcErrorCode_LogUserIn_PasswordIsWrong               = 28
	RpcErrorCode_LogUserIn_PasswordCheckError            = 29
	RpcErrorCode_LogUserIn_VerificationCodeGenerator     = 30
	RpcErrorCode_LogUserIn_VerificationCodeIsNotSet      = 31
	RpcErrorCode_LogUserIn_VerificationCodeIsWrong       = 32
	RpcErrorCode_LogUserIn_JWTCreation                   = 33
	RpcErrorCode_ThisActionIsNotForLoggedUsers           = 34
	RpcErrorCode_LogUserIn_EmailAddresIsNotValid         = 35
	RpcErrorCode_LogUserIn_RequestIdIsNotSet             = 36
	RpcErrorCode_LogUserIn_AuthChallengeResponseIsNotSet = 37
	RpcErrorCode_UserIdIsNotSet                          = 38
	RpcErrorCode_UserIsNotFound                          = 39
)

// Messages.
const (
	RpcErrorMsg_RegisterUser_StepIsUnknown              = "RegisterUser: unknown step"
	RpcErrorMsg_RegisterUser_EmailAddresIsNotValid      = "RegisterUser: e-mail address is not valid"
	RpcErrorMsg_RegisterUser_EmailAddresIsUsed          = "RegisterUser: e-mail address is used"
	RpcErrorMsg_RegisterUser_VerificationCodeGenerator  = "RegisterUser: VCG error"
	RpcErrorMsgF_RegisterUser_SmtpModule                = "RegisterUser: SMTP module error: %s" // Template.
	RpcErrorMsg_RegisterUser_VerificationCodeIsNotSet   = "RegisterUser: verification code is not set"
	RpcErrorMsg_RegisterUser_VerificationCodeIsNotValid = "RegisterUser: verification code is not valid"
	RpcErrorMsg_RegisterUser_NameIsNotSet               = "RegisterUser: name is not set"
	RpcErrorMsg_RegisterUser_PasswordIsNotSet           = "RegisterUser: password is not set"
	RpcErrorMsg_RegisterUser_UserNameIsTooLong          = "RegisterUser: user name is too long"
	RpcErrorMsg_RegisterUser_UserNameIsUsed             = "RegisterUser: user name is already used"
	RpcErrorMsgF_RegisterUser_PasswordError             = "RegisterUser: password error: %s" // Template.
	RpcErrorMsgF_BPP_PackSymbols                        = "BPP: PackSymbols error: %s"       // Template.
	RpcErrorMsg_RegisterUser_UserPasswordIsTooLong      = "RegisterUser: password is too long"
	RpcErrorMsg_LogUserIn_StepIsUnknown                 = "LogUserIn: unknown step"
	RpcErrorMsg_LogUserIn_UserCanNotLogIn               = "LogUserIn: user can not log in"
	RpcErrorMsg_LogUserIn_UserIsAlreadyLoggedIn         = "LogUserIn: user is already logged in"
	RpcErrorMsg_LogUserIn_UserHasAlreadyStartedToLogIn  = "LogUserIn: user has already started to log in"
	RpcErrorMsg_BPP_GenerateRandomSalt                  = "BPP: GenerateRandomSalt error"
	RpcErrorMsgF_RCS_CreateCaptcha                      = "RCS: CreateCaptcha error: %s" // Template.
	RpcErrorMsg_LogUserIn_RequestIdGenerator            = "LogUserIn: RIDG error"
	RpcErrorMsg_LogUserIn_UserHasNotStartedToLogIn      = "LogUserIn: user has not started to log in"
	RpcErrorMsg_LogUserIn_UserPreSessionIsNotFound      = "LogUserIn: user's preliminary session is not found"
	RpcErrorMsgF_RCS_CheckCaptcha                       = "RCS: CheckCaptcha error: %s" // Template.
	RpcErrorMsg_LogUserIn_CaptchaAnswerIsWrong          = "LogUserIn: captcha answer is wrong"
	RpcErrorMsg_LogUserIn_PasswordIsWrong               = "LogUserIn: password is wrong"
	RpcErrorMsgF_LogUserIn_PasswordCheckError           = "LogUserIn: password check error: %s" // Template.
	RpcErrorMsg_LogUserIn_VerificationCodeGenerator     = "LogUserIn: VCG error"
	RpcErrorMsg_LogUserIn_VerificationCodeIsNotSet      = "LogUserIn: verification code is not set"
	RpcErrorMsg_LogUserIn_VerificationCodeIsWrong       = "LogUserIn: verification code is wrong"
	RpcErrorMsgF_LogUserIn_JWTCreation                  = "LogUserIn: JWT creation error: %s" // Template.
	RpcErrorMsg_ThisActionIsNotForLoggedUsers           = "this action is not for logged in users"
	RpcErrorMsg_LogUserIn_EmailAddresIsNotValid         = "LogUserIn: e-mail address is not valid"
	RpcErrorMsg_LogUserIn_RequestIdIsNotSet             = "LogUserIn: request ID is not set"
	RpcErrorMsg_LogUserIn_AuthChallengeResponseIsNotSet = "LogUserIn: authorization challenge response is not set"
	RpcErrorMsg_UserIdIsNotSet                          = "user ID is not set"
	RpcErrorMsg_UserIsNotFound                          = "user is not found"
)
