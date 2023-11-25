package s

import (
	"errors"
	"fmt"
	"os"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	"golang.org/x/term"
)

// HttpSettings are settings of a generic HTTP server.
type HttpSettings struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

// HttpsSettings are settings of a generic HTTPS server.
type HttpsSettings struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

// DbSettings are parameters of a MySQL database together with some other
// settings related to database infrastructure. When a password is not set, it
// is taken from the stdin.
type DbSettings struct {
	// Access settings.
	DriverName string `json:"driverName"`
	Net        string `json:"net"`
	Host       string `json:"host"`
	Port       uint16 `json:"port"`
	DBName     string `json:"dbName"`
	User       string `json:"user"`
	Password   string `json:"password"`

	// Various specific MySQL settings.
	AllowNativePasswords bool              `json:"allowNativePasswords"`
	CheckConnLiveness    bool              `json:"checkConnLiveness"`
	MaxAllowedPacket     int               `json:"maxAllowedPacket"`
	Params               map[string]string `json:"params"`

	// Database structure and initialization settings.
	TableNamePrefix        string   `json:"tableNamePrefix"`
	TablesToInit           []string `json:"tablesToInit"`
	TableInitScriptsFolder string   `json:"tableInitScriptsFolder"`
}

func GetPasswordFromStdin(hint string) (pwd string, err error) {
	fmt.Println(hint)

	var buf []byte
	buf, err = term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func CheckSettingsFilePath(sfp string) (err error) {
	if len(sfp) == 0 {
		return errors.New(c.MsgFileIsNotSet)
	}

	return nil
}

func CheckHttpSettings(hs HttpSettings) (err error) {
	if (len(hs.Host) == 0) || (hs.Port == 0) {
		return errors.New(c.MsgHttpSettingError)
	}

	return nil
}

func CheckHttpsSettings(hss HttpsSettings) (err error) {
	if (len(hss.Host) == 0) || (hss.Port == 0) {
		return errors.New(c.MsgHttpsSettingError)
	}

	if (len(hss.CertFile) == 0) || (len(hss.KeyFile) == 0) {
		return errors.New(c.MsgHttpsSettingError)
	}

	return nil
}

func CheckDatabaseSettings(dbs DbSettings) (err error) {
	if (len(dbs.DriverName) == 0) || (len(dbs.Net) == 0) {
		return errors.New(c.MsgDbSettingError)
	}

	if (len(dbs.Host) == 0) || (dbs.Port == 0) {
		return errors.New(c.MsgDbSettingError)
	}

	if (len(dbs.DBName) == 0) || (len(dbs.User) == 0) {
		return errors.New(c.MsgDbSettingError)
	}

	return nil
}
