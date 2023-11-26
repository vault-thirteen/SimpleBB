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

	IntHttpSettings  `json:"intHttp"`
	ExtHttpsSettings `json:"extHttps"`
	DbSettings       `json:"db"`
	SystemSettings   `json:"system"`
}

// IntHttpSettings are settings of an HTTP server for internal requests.
type IntHttpSettings = cs.HttpSettings

// ExtHttpsSettings are settings of an HTTPS server for external requests.
type ExtHttpsSettings = cs.HttpsSettings

// DbSettings are parameters of the Database.
type DbSettings = cs.DbSettings

// SystemSettings are system settings.
type SystemSettings struct {
	SiteName   string `json:"siteName"`
	SiteDomain string `json:"siteDomain"`

	// Gateway module uses this path for HTTP handlers.
	// This setting must be synchronized with settings of the ACM module.
	EmailVerificationUrlPath string `json:"emailVerificationUrlPath"`

	// This setting must be synchronized with settings of the ACM module.
	IsFirewallUsed bool `json:"isFirewallUsed"`
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

	// Int. HTTP.
	err = cs.CheckHttpSettings(stn.IntHttpSettings)
	if err != nil {
		return err
	}

	// Ext. HTTPS.
	err = cs.CheckHttpsSettings(stn.ExtHttpsSettings)
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
	if len(stn.SystemSettings.EmailVerificationUrlPath) == 0 {
		return errors.New(c.MsgSystemSettingError)
	}

	return nil
}
