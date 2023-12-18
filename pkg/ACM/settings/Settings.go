package s

import (
	"encoding/json"
	"errors"
	"os"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cs "github.com/vault-thirteen/SimpleBB/pkg/common/settings"
)

// Settings is Server's settings.
type Settings struct {
	// Path to the file with these settings.
	FilePath string `json:"-"`

	HttpsSettings          `json:"https"`
	DbSettings             `json:"db"`
	SystemSettings         `json:"system"`
	SmtpModuleSettings     `json:"smtp"`
	CaptchaServiceSettings `json:"captcha"`
	JWTSettings            `json:"jwt"`
	UserRoleSettings       `json:"role"`
	GatewayModuleSettings  `json:"gw"`
}

// HttpsSettings are settings of an HTTPS server for incoming requests.
type HttpsSettings = cs.HttpsSettings

// DbSettings are parameters of the database.
type DbSettings = cs.DbSettings

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

// SmtpModuleSettings are settings of an SMTP module.
type SmtpModuleSettings struct {
	Host                              string `json:"host"`
	Port                              uint16 `json:"port"`
	Path                              string `json:"path"`
	MessageSubjectTemplate            string `json:"messageSubjectTemplate"`
	MessageBodyTemplateForReg         string `json:"messageBodyTemplateForReg"`
	MessageBodyTemplateForLogIn       string `json:"messageBodyTemplateForLogIn"`
	MessageBodyTemplateForPwdChange   string `json:"messageBodyTemplateForPwdChange"`
	MessageBodyTemplateForEmailChange string `json:"messageBodyTemplateForEmailChange"`
}

// CaptchaServiceSettings are settings of a Captcha service.
type CaptchaServiceSettings struct {
	// RPC Server.
	RpcHost string `json:"rpcHost"`
	RpcPort uint16 `json:"rpcPort"`
	RpcPath string `json:"rpcPath"`

	// Images Server.
	ImgHost string `json:"imgHost"`
	ImgPort uint16 `json:"imgPort"`
	ImgPath string `json:"imgPath"`
}

// JWTSettings are settings for JSON web tokens.
type JWTSettings struct {
	PrivateKeyFilePath string `json:"privateKeyFilePath"`
	PublicKeyFilePath  string `json:"publicKeyFilePath"`
	SigningMethod      string `json:"signingMethod"`
}

// UserRoleSettings are settings for special user roles.
type UserRoleSettings struct {
	// List of IDs of users having a moderator role.
	ModeratorIds []uint `json:"moderatorIds"`

	// List of IDs of users having an administrator role.
	AdministratorIds []uint `json:"administratorIds"`
}

// GatewayModuleSettings are settings of a Gateway module.
type GatewayModuleSettings struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
	Path string `json:"path"`
}

func NewSettingsFromFile(filePath string) (stn *Settings, err error) {
	var buf []byte
	buf, err = os.ReadFile(filePath)
	if err != nil {
		return stn, err
	}

	stn = &Settings{}
	err = json.Unmarshal(buf, stn)
	if err != nil {
		return stn, err
	}

	stn.FilePath = filePath

	err = stn.Check()
	if err != nil {
		return stn, err
	}

	if len(stn.Password) == 0 {
		stn.DbSettings.Password, err = cs.GetPasswordFromStdin(c.MsgEnterDatabasePassword)
		if err != nil {
			return stn, err
		}
	}

	return stn, nil
}

func (stn *Settings) Check() (err error) {
	err = cs.CheckSettingsFilePath(stn.FilePath)
	if err != nil {
		return err
	}

	// HTTPS.
	err = cs.CheckHttpsSettings(stn.HttpsSettings)
	if err != nil {
		return err
	}

	// DB.
	err = cs.CheckDatabaseSettings(stn.DbSettings)
	if err != nil {
		return err
	}

	// System.
	if (len(stn.SystemSettings.SiteName) == 0) ||
		(len(stn.SystemSettings.SiteDomain) == 0) ||
		(stn.SystemSettings.VerificationCodeLength < 8) ||
		(stn.SystemSettings.UserNameMaxLenInBytes == 0) ||
		(stn.SystemSettings.UserPasswordMaxLenInBytes == 0) ||
		(stn.SystemSettings.PreRegUserExpirationTime == 0) ||
		(stn.SystemSettings.LogInRequestIdLength == 0) ||
		(stn.SystemSettings.LogInTryTimeout == 0) ||
		(stn.SystemSettings.PreSessionExpirationTime == 0) ||
		(stn.SystemSettings.SessionMaxDuration == 0) ||
		(stn.SystemSettings.PasswordChangeExpirationTime == 0) ||
		(stn.SystemSettings.EmailChangeExpirationTime == 0) ||
		(stn.SystemSettings.ActionTryTimeout == 0) {
		return errors.New(c.MsgSystemSettingError)
	}

	// Firewall.
	if stn.SystemSettings.IsTableOfIncidentsUsed {
		if (stn.SystemSettings.BlockTimePerIncident.IllegalAccessAttempt == 0) ||
			(stn.SystemSettings.BlockTimePerIncident.FakeToken == 0) ||
			(stn.SystemSettings.BlockTimePerIncident.VerificationCodeMismatch == 0) ||
			(stn.SystemSettings.BlockTimePerIncident.DoubleLogInAttempt == 0) ||
			(stn.SystemSettings.BlockTimePerIncident.PreSessionHacking == 0) ||
			(stn.SystemSettings.BlockTimePerIncident.CaptchaAnswerMismatch == 0) ||
			(stn.SystemSettings.BlockTimePerIncident.PasswordMismatch == 0) ||
			(stn.SystemSettings.BlockTimePerIncident.PasswordChangeHacking == 0) ||
			(stn.SystemSettings.BlockTimePerIncident.EmailChangeHacking == 0) {
			return errors.New(c.MsgSystemSettingError)
		}
	}

	// SMTP.
	if (len(stn.SmtpModuleSettings.Host) == 0) ||
		(stn.SmtpModuleSettings.Port == 0) ||
		(len(stn.SmtpModuleSettings.Path) == 0) ||
		(len(stn.SmtpModuleSettings.MessageSubjectTemplate) == 0) ||
		(len(stn.SmtpModuleSettings.MessageBodyTemplateForReg) == 0) ||
		(len(stn.SmtpModuleSettings.MessageBodyTemplateForLogIn) == 0) ||
		(len(stn.SmtpModuleSettings.MessageBodyTemplateForPwdChange) == 0) ||
		(len(stn.SmtpModuleSettings.MessageBodyTemplateForEmailChange) == 0) {
		return errors.New(c.MsgSmtpSettingError)
	}

	// Captcha.
	if (len(stn.CaptchaServiceSettings.RpcHost) == 0) ||
		(stn.CaptchaServiceSettings.RpcPort == 0) ||
		(len(stn.CaptchaServiceSettings.RpcPath) == 0) ||
		(len(stn.CaptchaServiceSettings.ImgHost) == 0) ||
		(stn.CaptchaServiceSettings.ImgPort == 0) ||
		(len(stn.CaptchaServiceSettings.ImgPath) == 0) {
		return errors.New(c.MsgCaptchaSettingError)
	}

	// JWT.
	if (len(stn.PrivateKeyFilePath) == 0) ||
		(len(stn.PublicKeyFilePath) == 0) ||
		(len(stn.SigningMethod) == 0) {
		return errors.New(c.MsgJwtSettingError)
	}

	// Gateway module.
	if (len(stn.GatewayModuleSettings.Host) == 0) ||
		(stn.GatewayModuleSettings.Port == 0) ||
		(len(stn.GatewayModuleSettings.Path) == 0) {
		return errors.New(c.MsgGatewaySettingError)
	}

	return nil
}
