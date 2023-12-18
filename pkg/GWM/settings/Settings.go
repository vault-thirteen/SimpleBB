package s

import (
	"encoding/json"
	"errors"
	"os"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cs "github.com/vault-thirteen/SimpleBB/pkg/common/settings"
)

const (
	ClientIPAddressSource_Direct       = 1
	ClientIPAddressSource_CustomHeader = 2
)

const (
	ErrUnknownClientIPAddressSource = "unknown client IP address source"
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

	// This setting must be synchronised with settings of the ACM module.
	IsFirewallUsed bool `json:"isFirewallUsed"`

	// ClientIPAddressSource setting selects where to search for client's IP
	// address. '1' means that IP address is taken directly from the client's
	// address of the HTTP request; '2' means that IP address is taken from the
	// custom HTTP header which is configured by the ClientIPAddressHeader
	// setting. One of the most common examples of a custom header may be the
	// 'X-Forwarded-For' HTTP header. For most users the first variant ('1') is
	// the most suitable. The second variant ('2') may be used if you are
	// proxying requests of your clients somewhere inside your own network
	// infrastructure, such as via a load balancer or with a reverse proxy.
	ClientIPAddressSource byte   `json:"clientIPAddressSource"`
	ClientIPAddressHeader string `json:"clientIPAddressHeader"`

	IsFrontEndEnabled bool   `json:"isFrontEndEnabled"`
	FrontEndPath      string `json:"frontEndPath"`
	ApiPath           string `json:"apiPath"`
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
	if (len(stn.SystemSettings.SiteName) == 0) ||
		(len(stn.SystemSettings.SiteDomain) == 0) ||
		(stn.SystemSettings.ClientIPAddressSource < ClientIPAddressSource_Direct) ||
		(stn.SystemSettings.ClientIPAddressSource > ClientIPAddressSource_CustomHeader) ||
		(len(stn.SystemSettings.FrontEndPath) == 0) ||
		(len(stn.SystemSettings.ApiPath) == 0) ||
		(stn.SystemSettings.FrontEndPath == stn.SystemSettings.ApiPath) {
		return errors.New(c.MsgSystemSettingError)
	}

	if stn.SystemSettings.ClientIPAddressSource == ClientIPAddressSource_CustomHeader {
		if len(stn.SystemSettings.ClientIPAddressHeader) == 0 {
			return errors.New(c.MsgSystemSettingError)
		}
	}

	return nil
}
