package models

import (
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

type RequestId = uint

type PingParams struct{}

type PingResult struct {
	OK bool `json:"ok"`
}

type RegisterUserParams struct {
	cm.CommonParams

	// Step number.
	StepN byte `json:"stepN"`

	// E-mail address.
	// Is used on steps 1, 2 and 3.
	Email string `json:"email"`

	// Verification code.
	// Is used on steps 2 and 3.
	VerificationCode string `json:"verificationCode"`

	// Name.
	// Is used on step 3.
	Name string `json:"name"`

	// Password.
	// Is used on step 3.
	Password string `json:"password"`
}

type RegisterUserResult struct {
	cm.CommonResult

	// Next required step.
	// If set to zero, no further step is required.
	NextStep byte `json:"nextStep"`
}

type ApproveAndRegisterUserParams struct {
	cm.CommonParams

	// E-mail address of an approved user.
	Email string `json:"email"`
}

type ApproveAndRegisterUserResult struct {
	cm.CommonResult

	OK bool `json:"ok"`
}

type LogUserInParams struct {
	cm.CommonParams

	// Step number.
	StepN byte `json:"stepN"`

	// E-mail address.
	// This is the main identifier of a user.
	// It is used on all steps.
	Email string `json:"email"`

	// Request ID.
	// It protects preliminary sessions from being hi-jacked.
	// Is used on steps 2 and 3.
	RequestId string `json:"requestId"`

	// Captcha answer.
	// This field is optional and may be used on step 2.
	CaptchaAnswer string `json:"captchaAnswer"`

	// Authentication data provided for the challenge.
	// Is used on step 2.
	AuthChallengeResponse []byte `json:"authChallengeResponse"`

	// Verification Code.
	// Is used on step 3.
	VerificationCode string `json:"verificationCode"`
}

type LogUserInResult struct {
	cm.CommonResult

	// Next required step.
	// If set to zero, no further step is required.
	NextStep byte `json:"nextStep"`

	RequestId     string `json:"requestId"`
	AuthDataBytes []byte `json:"authDataBytes"`

	// Captcha parameters.
	IsCaptchaNeeded bool   `json:"isCaptchaNeeded"`
	CaptchaId       string `json:"captchaId"`

	// JWT key maker.
	IsWebTokenSet bool   `json:"isWebTokenSet"`
	WTS           string `json:"wts"` // Web token string.
}

type LogUserOutParams struct {
	cm.CommonParams
}

type LogUserOutResult struct {
	cm.CommonResult

	OK bool `json:"ok"`
}

type GetListOfLoggedUsersParams struct {
	cm.CommonParams
}

type GetListOfLoggedUsersResult struct {
	cm.CommonResult

	LoggedUserIds []uint `json:"loggedUserIds"`
}

type IsUserLoggedInParams struct {
	cm.CommonParams

	UserId uint `json:"userId"`
}

type IsUserLoggedInResult struct {
	cm.CommonResult

	UserId         uint `json:"userId"`
	IsUserLoggedIn bool `json:"isUserLoggedIn"`
}

type GetUserRolesParams struct {
	cm.CommonParams

	UserId uint `json:"userId"`
}

type GetUserRolesResult struct {
	cm.CommonResult

	UserId uint `json:"userId"`
	cm.UserRoles
}

type ViewUserParametersParams struct {
	cm.CommonParams

	UserId uint `json:"userId"`
}

type ViewUserParametersResult struct {
	cm.CommonResult

	UserId uint `json:"userId"`
	cm.UserParameters
}

type SetUserRoleAuthorParams struct {
	cm.CommonParams

	UserId        uint `json:"userId"`
	IsRoleEnabled bool `json:"isRoleEnabled"`
}

type SetUserRoleAuthorResult struct {
	cm.CommonResult

	OK bool `json:"ok"`
}

type SetUserRoleWriterParams struct {
	cm.CommonParams

	UserId        uint `json:"userId"`
	IsRoleEnabled bool `json:"isRoleEnabled"`
}

type SetUserRoleWriterResult struct {
	cm.CommonResult

	OK bool `json:"ok"`
}

type SetUserRoleReaderParams struct {
	cm.CommonParams

	UserId        uint `json:"userId"`
	IsRoleEnabled bool `json:"isRoleEnabled"`
}

type SetUserRoleReaderResult struct {
	cm.CommonResult

	OK bool `json:"ok"`
}

type BanUserParams struct {
	cm.CommonParams

	UserId uint `json:"userId"`
}

type BanUserResult struct {
	cm.CommonResult

	OK bool `json:"ok"`
}

type UnbanUserParams struct {
	cm.CommonParams

	UserId uint `json:"userId"`
}

type UnbanUserResult struct {
	cm.CommonResult

	OK bool `json:"ok"`
}

type GetSelfRolesParams struct {
	cm.CommonParams
}

type GetSelfRolesResult struct {
	cm.CommonResult

	UserId uint `json:"userId"`
	cm.UserRoles
}

type ShowDiagnosticDataParams struct{}

type ShowDiagnosticDataResult struct {
	cm.CommonResult

	TotalRequestsCount      uint64 `json:"totalRequestsCount"`
	SuccessfulRequestsCount uint64 `json:"successfulRequestsCount"`
}

type TestParams struct{}

type TestResult struct {
	cm.CommonResult
}
