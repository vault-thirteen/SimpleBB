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
	ErrFileIsNotSet = "file is not set"
	ErrHttpSetting  = "error in HTTP server setting"
	ErrDbSetting    = "error in database server setting"
	ErrACMSetting   = "error is ACM module setting"
)

// Settings is Server's settings.
type Settings struct {
	// Path to the file with these settings.
	FilePath string `json:"-"`

	HttpSettings   `json:"http"`
	DbSettings     `json:"db"`
	SystemSettings `json:"system"`
	ACMSettings    `json:"acm"`
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

	// ACM.
	if (len(stn.ACMSettings.Host) == 0) || (stn.ACMSettings.Port == 0) {
		return errors.New(ErrACMSetting)
	}
	if len(stn.ACMSettings.Path) == 0 {
		return errors.New(ErrACMSetting)
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
