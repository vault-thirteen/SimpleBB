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
}

// HttpsSettings are settings of an HTTPS server for incoming requests.
type HttpsSettings = cs.HttpsSettings

// DbSettings are parameters of the Database.
type DbSettings = cs.DbSettings

// SystemSettings are system settings.
type SystemSettings struct {
	SiteName               string `json:"siteName"`
	SiteDomain             string `json:"siteDomain"`
	VerificationCodeLength int    `json:"verificationCodeLength"`

	// ACM module uses this path for composition of e-mail messages.
	// This setting must be synchronized with settings of the Gateway module.
	EmailVerificationUrlPath string `json:"emailVerificationUrlPath"`

	UserNameMaxLenInBytes     uint `json:"userNameMaxLenInBytes"`
	UserPasswordMaxLenInBytes uint `json:"userPasswordMaxLenInBytes"`
	PreRegUserExpirationTime  uint `json:"preRegUserExpirationTime"`
	IsAdminApprovalRequired   bool `json:"isAdminApprovalRequired"`
	LogInRequestIdLength      int  `json:"logInRequestIdLength"`
	LogInTryTimeout           uint `json:"logInTryTimeout"`
	PreSessionExpirationTime  uint `json:"preSessionExpirationTime"`
	SessionMaxDuration        uint `json:"sessionMaxDuration"`

	// This setting must be synchronized with settings of the Gateway module.
	IsTableOfIncidentsUsed bool `json:"isTableOfIncidentsUsed"`

	IsDebugMode bool `json:"isDebugMode"`
}

// SmtpModuleSettings are settings of an SMTP Module.
type SmtpModuleSettings struct {
	Host                        string `json:"host"`
	Port                        uint16 `json:"port"`
	Path                        string `json:"path"`
	MessageSubjectTemplate      string `json:"messageSubjectTemplate"`
	MessageBodyTemplateForReg   string `json:"messageBodyTemplateForReg"`
	MessageBodyTemplateForLogIn string `json:"messageBodyTemplateForLogIn"`
}

// CaptchaServiceSettings are settings of a Captcha Service.
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

// JWTSettings are settings for JSON Web Tokens.
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
	if (len(stn.SystemSettings.SiteName) == 0) || (len(stn.SystemSettings.SiteDomain) == 0) {
		return errors.New(c.MsgSystemSettingError)
	}
	if stn.SystemSettings.VerificationCodeLength < 8 {
		return errors.New(c.MsgSystemSettingError)
	}
	if len(stn.SystemSettings.EmailVerificationUrlPath) == 0 {
		return errors.New(c.MsgSystemSettingError)
	}

	// SMTP.
	if (len(stn.SmtpModuleSettings.Host) == 0) || (stn.SmtpModuleSettings.Port == 0) {
		return errors.New(c.MsgSmtpSettingError)
	}
	if len(stn.SmtpModuleSettings.Path) == 0 {
		return errors.New(c.MsgSmtpSettingError)
	}
	if len(stn.SmtpModuleSettings.MessageSubjectTemplate) == 0 {
		return errors.New(c.MsgSmtpSettingError)
	}
	if len(stn.SmtpModuleSettings.MessageBodyTemplateForReg) == 0 {
		return errors.New(c.MsgSmtpSettingError)
	}

	// Captcha.
	if (len(stn.CaptchaServiceSettings.RpcHost) == 0) || (stn.CaptchaServiceSettings.RpcPort == 0) {
		return errors.New(c.MsgCaptchaSettingError)
	}
	if len(stn.CaptchaServiceSettings.RpcPath) == 0 {
		return errors.New(c.MsgCaptchaSettingError)
	}
	if (len(stn.CaptchaServiceSettings.ImgHost) == 0) || (stn.CaptchaServiceSettings.ImgPort == 0) {
		return errors.New(c.MsgCaptchaSettingError)
	}
	if len(stn.CaptchaServiceSettings.ImgPath) == 0 {
		return errors.New(c.MsgCaptchaSettingError)
	}

	return nil
}
