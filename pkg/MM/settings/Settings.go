package s

import (
	"encoding/json"
	"os"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	"github.com/vault-thirteen/SimpleBB/pkg/common/app"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cs "github.com/vault-thirteen/SimpleBB/pkg/common/settings"
	ver "github.com/vault-thirteen/auxie/Versioneer"
)

// Settings is Server's settings.
type Settings struct {
	// Path to the file with these settings.
	FilePath string `json:"-"`

	// Program versioning information.
	VersionInfo *ver.Versioneer `json:"-"`

	HttpsSettings  `json:"https"`
	DbSettings     `json:"db"`
	SystemSettings `json:"system"`

	// External services.
	AcmSettings cs.ServiceClientSettings `json:"acm"`
	NmSettings  cs.ServiceClientSettings `json:"nm"`
	SmSettings  cs.ServiceClientSettings `json:"sm"`
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

	// External services.
	err = stn.AcmSettings.Check()
	if err != nil {
		return cs.DetailedScsError(app.ServiceShortName_ACM, err)
	}
	err = stn.NmSettings.Check()
	if err != nil {
		return cs.DetailedScsError(app.ServiceShortName_NM, err)
	}
	err = stn.SmSettings.Check()
	if err != nil {
		return cs.DetailedScsError(app.ServiceShortName_SM, err)
	}

	return nil
}

func (stn *Settings) UseConstructor(filePath string, versionInfo *ver.Versioneer) (cm.ISettings, error) {
	return NewSettingsFromFile(filePath, versionInfo)
}
