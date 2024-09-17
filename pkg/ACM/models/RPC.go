package models

import (
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

type RequestId = uint

// Ping.

type PingParams = cmr.PingParams
type PingResult = cmr.PingResult

// User registration.

type RegisterUserParams struct {
	cmr.CommonParams

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
	cmr.CommonResult

	// Next required step.
	// If set to zero, no further step is required.
	NextStep byte `json:"nextStep"`
}

type ApproveAndRegisterUserParams struct {
	cmr.CommonParams

	// E-mail address of an approved user.
	Email string `json:"email"`
}

type ApproveAndRegisterUserResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

// Logging in and out.

type LogUserInParams struct {
	cmr.CommonParams

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
	cmr.CommonResult

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
	WTS           string `json:"wts,omitempty"` // Web token string.
}

type LogUserOutParams struct {
	cmr.CommonParams
}

type LogUserOutResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

type GetListOfLoggedUsersParams struct {
	cmr.CommonParams
}

type GetListOfLoggedUsersResult struct {
	cmr.CommonResult

	LoggedUserIds []uint `json:"loggedUserIds"`
}

type IsUserLoggedInParams struct {
	cmr.CommonParams

	UserId uint `json:"userId"`
}

type IsUserLoggedInResult struct {
	cmr.CommonResult

	UserId         uint `json:"userId"`
	IsUserLoggedIn bool `json:"isUserLoggedIn"`
}

// Various actions.

type ChangePasswordParams struct {
	cmr.CommonParams

	// Step number.
	StepN byte `json:"stepN"`

	// New password.
	// Is used on step 1.
	NewPassword string `json:"newPassword"`

	// Request ID.
	// It protects password changes from being hi-jacked.
	// Is used on step 2.
	RequestId string `json:"requestId"`

	// Authentication data provided for the challenge.
	// Is used on step 2.
	AuthChallengeResponse []byte `json:"authChallengeResponse"`

	// Verification Code.
	// Is used on step 2.
	VerificationCode string `json:"verificationCode"`

	// Captcha answer.
	// This field is optional and may be used on step 2.
	CaptchaAnswer string `json:"captchaAnswer"`
}

type ChangePasswordResult struct {
	cmr.CommonResult

	// Next required step.
	// If set to zero, no further step is required.
	NextStep byte `json:"nextStep"`

	RequestId     string `json:"requestId"`
	AuthDataBytes []byte `json:"authDataBytes"`

	// Captcha parameters.
	IsCaptchaNeeded bool   `json:"isCaptchaNeeded"`
	CaptchaId       string `json:"captchaId"`

	// Success flag.
	OK bool `json:"ok"`
}

type ChangeEmailParams struct {
	cmr.CommonParams

	// Step number.
	StepN byte `json:"stepN"`

	// New e-mail address.
	// Is used on step 1.
	NewEmail string `json:"newEmail"`

	// Request ID.
	// It protects e-mail changes from being hi-jacked.
	// Is used on step 2.
	RequestId string `json:"requestId"`

	// Authentication data provided for the challenge.
	// Is used on step 2.
	AuthChallengeResponse []byte `json:"authChallengeResponse"`

	// Verification Code for the old e-mail.
	// Is used on step 2.
	VerificationCodeOld string `json:"verificationCodeOld"`

	// Verification Code for the new e-mail.
	// Is used on step 2.
	VerificationCodeNew string `json:"verificationCodeNew"`

	// Captcha answer.
	// This field is optional and may be used on step 2.
	CaptchaAnswer string `json:"captchaAnswer"`
}

type ChangeEmailResult struct {
	cmr.CommonResult

	// Next required step.
	// If set to zero, no further step is required.
	NextStep byte `json:"nextStep"`

	RequestId     string `json:"requestId"`
	AuthDataBytes []byte `json:"authDataBytes"`

	// Captcha parameters.
	IsCaptchaNeeded bool   `json:"isCaptchaNeeded"`
	CaptchaId       string `json:"captchaId"`

	// Success flag.
	OK bool `json:"ok"`
}

type GetListOfAllUsersParams struct {
	cmr.CommonParams

	Page uint `json:"page"`
}

type GetListOfAllUsersResult struct {
	cmr.CommonResult

	UserIds    []uint `json:"userIds"`
	Page       uint   `json:"page,omitempty"`
	TotalPages uint   `json:"totalPages,omitempty"`
	TotalUsers uint   `json:"totalUsers,omitempty"`
}

// User properties.

type GetUserRolesParams struct {
	cmr.CommonParams

	UserId uint `json:"userId"`
}

type GetUserRolesResult struct {
	cmr.CommonResult

	UserId uint `json:"userId"`
	cm.UserRoles
}

type ViewUserParametersParams struct {
	cmr.CommonParams

	UserId uint `json:"userId"`
}

type ViewUserParametersResult struct {
	cmr.CommonResult

	cm.UserParameters
}

type SetUserRoleAuthorParams struct {
	cmr.CommonParams

	UserId        uint `json:"userId"`
	IsRoleEnabled bool `json:"isRoleEnabled"`
}

type SetUserRoleAuthorResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

type SetUserRoleWriterParams struct {
	cmr.CommonParams

	UserId        uint `json:"userId"`
	IsRoleEnabled bool `json:"isRoleEnabled"`
}

type SetUserRoleWriterResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

type SetUserRoleReaderParams struct {
	cmr.CommonParams

	UserId        uint `json:"userId"`
	IsRoleEnabled bool `json:"isRoleEnabled"`
}

type SetUserRoleReaderResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

type GetSelfRolesParams struct {
	cmr.CommonParams
}

type GetSelfRolesResult struct {
	cmr.CommonResult

	UserId uint `json:"userId"`
	cm.UserRoles
}

// User banning.

type BanUserParams struct {
	cmr.CommonParams

	UserId uint `json:"userId"`
}

type BanUserResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

type UnbanUserParams struct {
	cmr.CommonParams

	UserId uint `json:"userId"`
}

type UnbanUserResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

// Other.

type ShowDiagnosticDataParams struct{}

type ShowDiagnosticDataResult struct {
	cmr.CommonResult
	cmr.RequestsCount
}

type TestParams struct{}

type TestResult struct {
	cmr.CommonResult
}
