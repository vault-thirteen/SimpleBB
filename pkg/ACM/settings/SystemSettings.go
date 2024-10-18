package s

import (
	"errors"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

// SystemSettings are system settings.
type SystemSettings struct {
	SiteName                     cmb.Text  `json:"siteName"`
	SiteDomain                   cmb.Text  `json:"siteDomain"`
	VerificationCodeLength       cmb.Count `json:"verificationCodeLength"`
	UserNameMaxLenInBytes        cmb.Count `json:"userNameMaxLenInBytes"`
	UserPasswordMaxLenInBytes    cmb.Count `json:"userPasswordMaxLenInBytes"`
	PreRegUserExpirationTime     cmb.Count `json:"preRegUserExpirationTime"`
	IsAdminApprovalRequired      cmb.Flag  `json:"isAdminApprovalRequired"`
	LogInRequestIdLength         cmb.Count `json:"logInRequestIdLength"`
	LogInTryTimeout              cmb.Count `json:"logInTryTimeout"`
	PreSessionExpirationTime     cmb.Count `json:"preSessionExpirationTime"`
	SessionMaxDuration           cmb.Count `json:"sessionMaxDuration"`
	PasswordChangeExpirationTime cmb.Count `json:"passwordChangeExpirationTime"`
	EmailChangeExpirationTime    cmb.Count `json:"emailChangeExpirationTime"`
	ActionTryTimeout             cmb.Count `json:"actionTryTimeout"`
	PageSize                     cmb.Count `json:"pageSize"`

	// This setting must be synchronised with settings of the Gateway module.
	IsTableOfIncidentsUsed cmb.Flag `json:"isTableOfIncidentsUsed"`

	// This setting is used only when a table of incidents is enabled.
	BlockTimePerIncident BlockTimePerIncident `json:"blockTimePerIncident"`

	IsDebugMode bool `json:"isDebugMode"`
}

// BlockTimePerIncident is block time in seconds for each type of incident.
type BlockTimePerIncident struct {
	IllegalAccessAttempt     cmb.Count `json:"illegalAccessAttempt"`     // 1.
	FakeToken                cmb.Count `json:"fakeToken"`                // 2.
	VerificationCodeMismatch cmb.Count `json:"verificationCodeMismatch"` // 3.
	DoubleLogInAttempt       cmb.Count `json:"doubleLogInAttempt"`       // 4.
	PreSessionHacking        cmb.Count `json:"preSessionHacking"`        // 5.
	CaptchaAnswerMismatch    cmb.Count `json:"captchaAnswerMismatch"`    // 6.
	PasswordMismatch         cmb.Count `json:"passwordMismatch"`         // 7.
	PasswordChangeHacking    cmb.Count `json:"passwordChangeHacking"`    // 8.
	EmailChangeHacking       cmb.Count `json:"emailChangeHacking"`       // 9.
	FakeIPA                  cmb.Count `json:"fakeIPA"`                  // 10.
}

func (s SystemSettings) Check() (err error) {
	if (len(s.SiteName) == 0) ||
		(len(s.SiteDomain) == 0) ||
		(s.VerificationCodeLength < 8) ||
		(s.UserNameMaxLenInBytes == 0) ||
		(s.UserPasswordMaxLenInBytes == 0) ||
		(s.PreRegUserExpirationTime == 0) ||
		(s.LogInRequestIdLength == 0) ||
		(s.LogInTryTimeout == 0) ||
		(s.PreSessionExpirationTime == 0) ||
		(s.SessionMaxDuration == 0) ||
		(s.PasswordChangeExpirationTime == 0) ||
		(s.EmailChangeExpirationTime == 0) ||
		(s.ActionTryTimeout == 0) ||
		(s.PageSize == 0) {
		return errors.New(c.MsgSystemSettingError)
	}

	// Firewall.
	if s.IsTableOfIncidentsUsed {
		if (s.BlockTimePerIncident.IllegalAccessAttempt == 0) ||
			(s.BlockTimePerIncident.FakeToken == 0) ||
			(s.BlockTimePerIncident.VerificationCodeMismatch == 0) ||
			(s.BlockTimePerIncident.DoubleLogInAttempt == 0) ||
			(s.BlockTimePerIncident.PreSessionHacking == 0) ||
			(s.BlockTimePerIncident.CaptchaAnswerMismatch == 0) ||
			(s.BlockTimePerIncident.PasswordMismatch == 0) ||
			(s.BlockTimePerIncident.PasswordChangeHacking == 0) ||
			(s.BlockTimePerIncident.EmailChangeHacking == 0) ||
			(s.BlockTimePerIncident.FakeIPA == 0) {
			return errors.New(c.MsgSystemSettingError)
		}
	}

	return nil
}
