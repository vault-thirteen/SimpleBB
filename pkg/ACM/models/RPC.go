package models

import (
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

// Ping.

type PingParams = cmr.PingParams
type PingResult = cmr.PingResult

// User registration.

type RegisterUserParams struct {
	cmr.CommonParams

	// Step number.
	StepN cm.StepNumber `json:"stepN"`

	// E-mail address. Is used on steps 1, 2 and 3.
	Email cm.Email `json:"email"`

	// Verification code. Is used on steps 2 and 3.
	VerificationCode cm.VerificationCode `json:"verificationCode"`

	// Name. Is used on step 3.
	Name cm.Name `json:"name"`

	// Password. Is used on step 3.
	Password cm.Password `json:"password"`
}

type RegisterUserResult struct {
	cmr.CommonResult

	// Next required step. If set to zero, no further step is required.
	NextStep cm.StepNumber `json:"nextStep"`
}

type GetListOfRegistrationsReadyForApprovalParams struct {
	cmr.CommonParams
	Page cmb.Count `json:"page"`
}

type GetListOfRegistrationsReadyForApprovalResult struct {
	cmr.CommonResult
	RRFA     []RegistrationReadyForApproval `json:"rrfa"`
	PageData *cmr.PageData                  `json:"pageData,omitempty"`
}

type RejectRegistrationRequestParams struct {
	cmr.CommonParams
	RegistrationRequestId cmb.Id `json:"registrationRequestId"`
}
type RejectRegistrationRequestResult = cmr.CommonResultWithSuccess

type ApproveAndRegisterUserParams struct {
	cmr.CommonParams

	// E-mail address of an approved user.
	Email cm.Email `json:"email"`
}
type ApproveAndRegisterUserResult = cmr.CommonResultWithSuccess

// Logging in and out.

type LogUserInParams struct {
	cmr.CommonParams

	// Step number.
	StepN cm.StepNumber `json:"stepN"`

	// E-mail address.
	// This is the main identifier of a user.
	// It is used on all steps.
	Email cm.Email `json:"email"`

	// Request ID.
	// It protects preliminary sessions from being hi-jacked.
	// Is used on steps 2 and 3.
	RequestId cm.RequestId `json:"requestId"`

	// Captcha answer.
	// This field is optional and may be used on step 2.
	CaptchaAnswer cm.CaptchaAnswer `json:"captchaAnswer"`

	// Authentication data provided for the challenge.
	// Is used on step 2.
	AuthChallengeResponse cmr.AuthChallengeData `json:"authChallengeResponse"`

	// Verification Code.
	// Is used on step 3.
	VerificationCode cm.VerificationCode `json:"verificationCode"`
}

type LogUserInResult struct {
	cmr.CommonResult

	// Next required step. If set to zero, no further step is required.
	NextStep cm.StepNumber `json:"nextStep"`

	RequestId     cm.RequestId          `json:"requestId"`
	AuthDataBytes cmr.AuthChallengeData `json:"authDataBytes"`

	// Captcha parameters.
	IsCaptchaNeeded cmb.Flag      `json:"isCaptchaNeeded"`
	CaptchaId       *cm.CaptchaId `json:"captchaId"`

	// JWT key maker.
	IsWebTokenSet  cmb.Flag          `json:"isWebTokenSet"`
	WebTokenString cm.WebTokenString `json:"wts,omitempty"`
}

type LogUserOutParams struct {
	cmr.CommonParams
}
type LogUserOutResult = cmr.CommonResultWithSuccess

type LogUserOutAParams struct {
	cmr.CommonParams
	UserId cmb.Id `json:"userId"`
}
type LogUserOutAResult = cmr.CommonResultWithSuccess

type GetListOfLoggedUsersParams struct {
	cmr.CommonParams
}

type GetListOfLoggedUsersResult struct {
	cmr.CommonResult
	LoggedUserIds []cmb.Id `json:"loggedUserIds"`
}

type GetListOfAllUsersParams struct {
	cmr.CommonParams
	Page cmb.Count `json:"page"`
}

type GetListOfAllUsersResult struct {
	cmr.CommonResult
	UserIds  []cmb.Id      `json:"userIds"`
	PageData *cmr.PageData `json:"pageData,omitempty"`
}

type IsUserLoggedInParams struct {
	cmr.CommonParams
	UserId cmb.Id `json:"userId"`
}

type IsUserLoggedInResult struct {
	cmr.CommonResult
	UserId         cmb.Id   `json:"userId"`
	IsUserLoggedIn cmb.Flag `json:"isUserLoggedIn"`
}

// Various actions.

type ChangePasswordParams struct {
	cmr.CommonParams

	// Step number.
	StepN cm.StepNumber `json:"stepN"`

	// New password.
	// Is used on step 1.
	NewPassword cm.Password `json:"newPassword"`

	// Request ID.
	// It protects password changes from being hi-jacked.
	// Is used on step 2.
	RequestId cm.RequestId `json:"requestId"`

	// Authentication data provided for the challenge.
	// Is used on step 2.
	AuthChallengeResponse cmr.AuthChallengeData `json:"authChallengeResponse"`

	// Verification Code.
	// Is used on step 2.
	VerificationCode cm.VerificationCode `json:"verificationCode"`

	// Captcha answer.
	// This field is optional and may be used on step 2.
	CaptchaAnswer cm.CaptchaAnswer `json:"captchaAnswer"`
}

type ChangePasswordResult struct {
	cmr.CommonResult
	cmr.Success

	// Next required step. If set to zero, no further step is required.
	NextStep cm.StepNumber `json:"nextStep"`

	RequestId     cm.RequestId          `json:"requestId"`
	AuthDataBytes cmr.AuthChallengeData `json:"authDataBytes"`

	// Captcha parameters.
	IsCaptchaNeeded cmb.Flag     `json:"isCaptchaNeeded"`
	CaptchaId       cm.CaptchaId `json:"captchaId"`
}

type ChangeEmailParams struct {
	cmr.CommonParams

	// Step number.
	StepN cm.StepNumber `json:"stepN"`

	// New e-mail address.
	// Is used on step 1.
	NewEmail cm.Email `json:"newEmail"`

	// Request ID.
	// It protects e-mail changes from being hi-jacked.
	// Is used on step 2.
	RequestId cm.RequestId `json:"requestId"`

	// Authentication data provided for the challenge.
	// Is used on step 2.
	AuthChallengeResponse cmr.AuthChallengeData `json:"authChallengeResponse"`

	// Verification Code for the old e-mail.
	// Is used on step 2.
	VerificationCodeOld cm.VerificationCode `json:"verificationCodeOld"`

	// Verification Code for the new e-mail.
	// Is used on step 2.
	VerificationCodeNew cm.VerificationCode `json:"verificationCodeNew"`

	// Captcha answer.
	// This field is optional and may be used on step 2.
	CaptchaAnswer cm.CaptchaAnswer `json:"captchaAnswer"`
}

type ChangeEmailResult struct {
	cmr.CommonResult
	cmr.Success

	// Next required step. If set to zero, no further step is required.
	NextStep cm.StepNumber `json:"nextStep"`

	RequestId     cm.RequestId          `json:"requestId"`
	AuthDataBytes cmr.AuthChallengeData `json:"authDataBytes"`

	// Captcha parameters.
	IsCaptchaNeeded cmb.Flag     `json:"isCaptchaNeeded"`
	CaptchaId       cm.CaptchaId `json:"captchaId"`
}

type GetUserSessionParams struct {
	cmr.CommonParams
	UserId cmb.Id `json:"userId"`
}

type GetUserSessionResult struct {
	cmr.CommonResult
	User    *User    `json:"user"`
	Session *Session `json:"session"`
}

// User properties.

type GetUserNameParams struct {
	cmr.CommonParams
	UserId cmb.Id `json:"userId"`
}

type GetUserNameResult struct {
	cmr.CommonResult
	User *User `json:"user"`
}

type GetUserRolesParams struct {
	cmr.CommonParams
	UserId cmb.Id `json:"userId"`
}

type GetUserRolesResult struct {
	cmr.CommonResult
	User *User `json:"user"`
}

type ViewUserParametersParams struct {
	cmr.CommonParams
	UserId cmb.Id `json:"userId"`
}

type ViewUserParametersResult struct {
	cmr.CommonResult
	User *User `json:"user"`
}

type SetUserRoleAuthorParams = SetUserRoleCommonParams
type SetUserRoleAuthorResult = SetUserRoleCommonResult
type SetUserRoleWriterParams = SetUserRoleCommonParams
type SetUserRoleWriterResult = SetUserRoleCommonResult
type SetUserRoleReaderParams = SetUserRoleCommonParams
type SetUserRoleReaderResult = SetUserRoleCommonResult

type GetSelfRolesParams struct {
	cmr.CommonParams
}

type GetSelfRolesResult struct {
	cmr.CommonResult
	User *User `json:"user"`
}

// User banning.

type BanUserParams struct {
	cmr.CommonParams
	UserId cmb.Id `json:"userId"`
}
type BanUserResult = cmr.CommonResultWithSuccess

type UnbanUserParams struct {
	cmr.CommonParams
	UserId cmb.Id `json:"userId"`
}
type UnbanUserResult = cmr.CommonResultWithSuccess

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

// Common models.

type SetUserRoleCommonParams struct {
	cmr.CommonParams
	UserId        cmb.Id   `json:"userId"`
	IsRoleEnabled cmb.Flag `json:"isRoleEnabled"`
}
type SetUserRoleCommonResult = cmr.CommonResultWithSuccess
