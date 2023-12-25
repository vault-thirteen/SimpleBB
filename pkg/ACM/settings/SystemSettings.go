package s

import (
	"errors"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
)

// SystemSettings are system settings.
type SystemSettings struct {
	SiteName                     string `json:"siteName"`
	SiteDomain                   string `json:"siteDomain"`
	VerificationCodeLength       int    `json:"verificationCodeLength"`
	UserNameMaxLenInBytes        uint   `json:"userNameMaxLenInBytes"`
	UserPasswordMaxLenInBytes    uint   `json:"userPasswordMaxLenInBytes"`
	PreRegUserExpirationTime     uint   `json:"preRegUserExpirationTime"`
	IsAdminApprovalRequired      bool   `json:"isAdminApprovalRequired"`
	LogInRequestIdLength         int    `json:"logInRequestIdLength"`
	LogInTryTimeout              uint   `json:"logInTryTimeout"`
	PreSessionExpirationTime     uint   `json:"preSessionExpirationTime"`
	SessionMaxDuration           uint   `json:"sessionMaxDuration"`
	PasswordChangeExpirationTime uint   `json:"passwordChangeExpirationTime"`
	EmailChangeExpirationTime    uint   `json:"emailChangeExpirationTime"`
	ActionTryTimeout             uint   `json:"actionTryTimeout"`

	// This setting must be synchronised with settings of the Gateway module.
	IsTableOfIncidentsUsed bool `json:"isTableOfIncidentsUsed"`

	// This setting is used only when a table of incidents is enabled.
	BlockTimePerIncident BlockTimePerIncident `json:"blockTimePerIncident"`

	IsDebugMode bool `json:"isDebugMode"`
}

// BlockTimePerIncident is block time in seconds for each type of incident.
type BlockTimePerIncident struct {
	IllegalAccessAttempt     uint `json:"illegalAccessAttempt"`     // 1.
	FakeToken                uint `json:"fakeToken"`                // 2.
	VerificationCodeMismatch uint `json:"verificationCodeMismatch"` // 3.
	DoubleLogInAttempt       uint `json:"doubleLogInAttempt"`       // 4.
	PreSessionHacking        uint `json:"preSessionHacking"`        // 5.
	CaptchaAnswerMismatch    uint `json:"captchaAnswerMismatch"`    // 6.
	PasswordMismatch         uint `json:"passwordMismatch"`         // 7.
	PasswordChangeHacking    uint `json:"passwordChangeHacking"`    // 8.
	EmailChangeHacking       uint `json:"emailChangeHacking"`       // 9.
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
		(s.ActionTryTimeout == 0) {
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
			(s.BlockTimePerIncident.EmailChangeHacking == 0) {
			return errors.New(c.MsgSystemSettingError)
		}
	}

	return nil
}
