package s

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	"golang.org/x/term"
)

const (
	ErrFileIsNotSet   = "file is not set"
	ErrHttpSetting    = "error in HTTP server setting"
	ErrDbSetting      = "error in database server setting"
	ErrSystemSetting  = "error in system setting"
	ErrSmtpSetting    = "error in SMTP module setting"
	ErrCaptchaSetting = "error in captcha service setting"
)

// Settings is Server's settings.
type Settings struct {
	// Path to the file with these settings.
	FilePath string `json:"-"`

	HttpSettings           `json:"http"`
	DbSettings             `json:"db"`
	SystemSettings         `json:"system"`
	SmtpModuleSettings     `json:"smtp"`
	CaptchaServiceSettings `json:"captcha"`
	JWTSettings            `json:"jwt"`
	UserRoleSettings       `json:"role"`
}

// HttpSettings are settings of an HTTP server for incoming requests.
type HttpSettings struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

// DbSettings are parameters of the Database. When a password is not set, it is
// taken from the stdin.
type DbSettings struct {
	// Access settings.
	DriverName string `json:"driverName"`
	Net        string `json:"net"`
	Host       string `json:"host"`
	Port       uint16 `json:"port"`
	DBName     string `json:"dbName"`
	User       string `json:"user"`
	Password   string `json:"password"`

	// Database structure and initialization settings.
	TableNamePrefix        string   `json:"tableNamePrefix"`
	TablesToInit           []string `json:"tablesToInit"`
	TableInitScriptsFolder string   `json:"tableInitScriptsFolder"`

	// Various specific MySQL settings.
	AllowNativePasswords bool              `json:"allowNativePasswords"`
	CheckConnLiveness    bool              `json:"checkConnLiveness"`
	MaxAllowedPacket     int               `json:"maxAllowedPacket"`
	Params               map[string]string `json:"params"`
}

// SystemSettings are system settings.
type SystemSettings struct {
	SiteName               string `json:"siteName"`
	SiteDomain             string `json:"siteDomain"`
	VerificationCodeLength int    `json:"verificationCodeLength"`

	// This path is used only for composition of e-mail messages.
	// It must be synchronized with settings of the Gateway Module.
	EmailVerificationUrlPath string `json:"emailVerificationUrlPath"`

	UserNameMaxLenInBytes     uint `json:"userNameMaxLenInBytes"`
	UserPasswordMaxLenInBytes uint `json:"userPasswordMaxLenInBytes"`
	PreRegUserExpirationTime  uint `json:"preRegUserExpirationTime"`
	IsAdminApprovalRequired   bool `json:"isAdminApprovalRequired"`
	LogInRequestIdLength      int  `json:"logInRequestIdLength"`
	LogInTryTimeout           uint `json:"logInTryTimeout"`
	PreSessionExpirationTime  uint `json:"preSessionExpirationTime"`
	SessionMaxDuration        uint `json:"sessionMaxDuration"`
	IsTableOfIncidentsUsed    bool `json:"isTableOfIncidentsUsed"`
	IsDebugMode               bool `json:"isDebugMode"`
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
		stn.DbSettings.Password, err = getPasswordFromStdin()
		if err != nil {
			return stn, err
		}
	}

	return stn, nil
}

func (stn *Settings) Check() (err error) {
	if len(stn.FilePath) == 0 {
		return errors.New(ErrFileIsNotSet)
	}

	// HTTP.
	if (len(stn.HttpSettings.Host) == 0) || (stn.HttpSettings.Port == 0) {
		return errors.New(ErrHttpSetting)
	}
	if (len(stn.HttpSettings.CertFile) == 0) || (len(stn.HttpSettings.KeyFile) == 0) {
		return errors.New(ErrHttpSetting)
	}

	// DB.
	if (len(stn.DbSettings.DriverName) == 0) || (len(stn.DbSettings.Net) == 0) {
		return errors.New(ErrDbSetting)
	}
	if (len(stn.DbSettings.Host) == 0) || (stn.DbSettings.Port == 0) {
		return errors.New(ErrDbSetting)
	}
	if (len(stn.DbSettings.DBName) == 0) || (len(stn.DbSettings.User) == 0) {
		return errors.New(ErrDbSetting)
	}

	// System.
	if (len(stn.SystemSettings.SiteName) == 0) || (len(stn.SystemSettings.SiteDomain) == 0) {
		return errors.New(ErrSystemSetting)
	}
	if stn.SystemSettings.VerificationCodeLength < 8 {
		return errors.New(ErrSystemSetting)
	}
	if len(stn.SystemSettings.EmailVerificationUrlPath) == 0 {
		return errors.New(ErrSystemSetting)
	}

	// SMTP.
	if (len(stn.SmtpModuleSettings.Host) == 0) || (stn.SmtpModuleSettings.Port == 0) {
		return errors.New(ErrSmtpSetting)
	}
	if len(stn.SmtpModuleSettings.Path) == 0 {
		return errors.New(ErrSmtpSetting)
	}
	if len(stn.SmtpModuleSettings.MessageSubjectTemplate) == 0 {
		return errors.New(ErrSmtpSetting)
	}
	if len(stn.SmtpModuleSettings.MessageBodyTemplateForReg) == 0 {
		return errors.New(ErrSmtpSetting)
	}

	// Captcha.
	if (len(stn.CaptchaServiceSettings.RpcHost) == 0) || (stn.CaptchaServiceSettings.RpcPort == 0) {
		return errors.New(ErrCaptchaSetting)
	}
	if len(stn.CaptchaServiceSettings.RpcPath) == 0 {
		return errors.New(ErrCaptchaSetting)
	}
	if (len(stn.CaptchaServiceSettings.ImgHost) == 0) || (stn.CaptchaServiceSettings.ImgPort == 0) {
		return errors.New(ErrCaptchaSetting)
	}
	if len(stn.CaptchaServiceSettings.ImgPath) == 0 {
		return errors.New(ErrCaptchaSetting)
	}

	return nil
}

func getPasswordFromStdin() (pwd string, err error) {
	fmt.Println(c.MsgEnterDatabasePassword)

	var buf []byte
	buf, err = term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}

	return string(buf), nil
}
