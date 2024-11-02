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

const (
	ErrUnknownClientIPAddressSource = "unknown client IP address source"
)

// Settings is Server's settings.
type Settings struct {
	// Path to the file with these settings.
	FilePath cm.Path `json:"-"`

	// Program versioning information.
	VersionInfo *ver.Versioneer `json:"-"`

	IntHttpSettings  `json:"intHttp"`
	ExtHttpsSettings `json:"extHttps"`
	DbSettings       `json:"db"`
	SystemSettings   `json:"system"`

	// External services.
	AcmSettings cs.ServiceClientSettings `json:"acm"`
	MmSettings  cs.ServiceClientSettings `json:"mm"`
	NmSettings  cs.ServiceClientSettings `json:"nm"`
	SmSettings  cs.ServiceClientSettings `json:"sm"`
}

func NewSettings() *Settings {
	return &Settings{
		SystemSettings: SystemSettings{
			ClientIPAddressSource: *cm.NewClientIPAddressSource(),
		},
	}
}

func NewSettingsFromFile(filePath string, versionInfo *ver.Versioneer) (stn *Settings, err error) {
	var buf []byte
	buf, err = os.ReadFile(filePath)
	if err != nil {
		return stn, err
	}

	stn = NewSettings()
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

	// Int. HTTP.
	err = stn.IntHttpSettings.Check()
	if err != nil {
		return err
	}

	// Ext. HTTPS.
	err = stn.ExtHttpsSettings.Check()
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

	err = stn.MmSettings.Check()
	if err != nil {
		return cs.DetailedScsError(app.ServiceShortName_MM, err)
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

func (stn *Settings) UseConstructor(filePath string, versionInfo *ver.Versioneer) (cmi.ISettings, error) {
	return NewSettingsFromFile(filePath, versionInfo)
}
