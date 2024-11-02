package s

import (
	"encoding/json"
	"os"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	"github.com/vault-thirteen/SimpleBB/pkg/common/app"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/models/interfaces"
	cs "github.com/vault-thirteen/SimpleBB/pkg/common/settings"
	ver "github.com/vault-thirteen/auxie/Versioneer"
)

// Settings is Server's settings.
type Settings struct {
	// Path to the file with these settings.
	FilePath cm.Path `json:"-"`

	// Program versioning information.
	VersionInfo *ver.Versioneer `json:"-"`

	HttpsSettings    `json:"https"`
	DbSettings       `json:"db"`
	SystemSettings   `json:"system"`
	JWTSettings      `json:"jwt"`
	UserRoleSettings `json:"role"`
	MessageSettings  `json:"message"`

	// Settings of host storing captcha images.
	CaptchaImageSettings `json:"captcha"`

	// External services.
	GwmSettings  cs.ServiceClientSettings `json:"gwm"`
	SmtpSettings cs.ServiceClientSettings `json:"smtp"`
	RcsSettings  cs.ServiceClientSettings `json:"rcs"`
}

func NewSettingsFromFile(filePath string, versionInfo *ver.Versioneer) (stn *Settings, err error) {
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

	stn.FilePath = cm.Path(filePath)

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

	stn.VersionInfo = versionInfo

	return stn, nil
}

func (stn *Settings) Check() (err error) {
	err = cs.CheckSettingsFilePath(stn.FilePath)
	if err != nil {
		return err
	}

	// HTTPS.
	err = stn.HttpsSettings.Check()
	if err != nil {
		return err
	}

	// DB.
	err = stn.DbSettings.Check()
	if err != nil {
		return err
	}

	// System.
	err = stn.SystemSettings.Check()
	if err != nil {
		return err
	}

	// JWT.
	err = stn.JWTSettings.Check()
	if err != nil {
		return err
	}

	// User Roles.
	err = stn.UserRoleSettings.Check()
	if err != nil {
		return err
	}

	// Message.
	err = stn.MessageSettings.Check()
	if err != nil {
		return err
	}

	// Captcha image.
	err = stn.CaptchaImageSettings.Check()
	if err != nil {
		return err
	}

	// External services.
	err = stn.GwmSettings.Check()
	if err != nil {
		return cs.DetailedScsError(app.ServiceShortName_GWM, err)
	}

	err = stn.SmtpSettings.Check()
	if err != nil {
		return cs.DetailedScsError(app.ServiceShortName_SMTP, err)
	}

	err = stn.RcsSettings.Check()
	if err != nil {
		return cs.DetailedScsError(app.ServiceShortName_RCS, err)
	}

	return nil
}

func (stn *Settings) UseConstructor(filePath string, versionInfo *ver.Versioneer) (cmi.ISettings, error) {
	return NewSettingsFromFile(filePath, versionInfo)
}
