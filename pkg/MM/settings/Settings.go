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

	HttpsSettings  `json:"https"`
	DbSettings     `json:"db"`
	SystemSettings `json:"system"`
	ACMSettings    `json:"acm"`
}

// HttpsSettings are settings of an HTTPS server for incoming requests.
type HttpsSettings = cs.HttpsSettings

// DbSettings are parameters of the Database.
type DbSettings = cs.DbSettings

// SystemSettings are system settings.
type SystemSettings struct {
	MessageEditTime uint `json:"messageEditTime"`
	PageSize        uint `json:"pageSize"`

	// NewThreadsAtTop parameter controls how new and updated threads are
	// placed inside forums. If set to 'True', then following will happen:
	// 1. New threads will be added to the start (top) of the list of forum's
	// threads instead of being added to the end (bottom) of the list;
	// 2. New messages added to threads will update the thread moving it to the
	// start (top) position of the list of forum's threads.
	// If set to 'False', then new threads are added to the end (bottom) of the
	// list and thread's new messages do not update thread's position in the
	// list.
	NewThreadsAtTop bool `json:"newThreadsAtTop"`

	IsDebugMode bool `json:"isDebugMode"`
}

// ACMSettings are settings of an Access Control Module.
type ACMSettings struct {
	Host                        string `json:"host"`
	Port                        uint16 `json:"port"`
	Path                        string `json:"path"`
	EnableSelfSignedCertificate bool   `json:"enableSelfSignedCertificate"`
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

	// ACM.
	if (len(stn.ACMSettings.Host) == 0) || (stn.ACMSettings.Port == 0) {
		return errors.New(c.MsgACMSettingError)
	}
	if len(stn.ACMSettings.Path) == 0 {
		return errors.New(c.MsgACMSettingError)
	}

	return nil
}
