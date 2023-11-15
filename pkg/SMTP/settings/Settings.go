package settings

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
	ErrSmtpSetting  = "error in SMTP server setting"
)

// Settings is Server's settings.
type Settings struct {
	// Path to the file with these settings.
	FilePath string `json:"-"`

	HttpSettings `json:"http"`
	SmtpSettings `json:"smtp"`
}

// HttpSettings are settings of an HTTP server for incoming requests.
type HttpSettings struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

// SmtpSettings are parameters of the SMTP server used for sending e-mail
// messages. When a password is not set, it is taken from the stdin.
type SmtpSettings struct {
	Host      string `json:"host"`
	Port      uint16 `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	UserAgent string `json:"userAgent"`
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
		stn.SmtpSettings.Password, err = getPasswordFromStdin()
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

	if (len(stn.HttpSettings.Host) == 0) || (stn.HttpSettings.Port == 0) {
		return errors.New(ErrHttpSetting)
	}

	if (len(stn.SmtpSettings.Host) == 0) || (stn.SmtpSettings.Port == 0) {
		return errors.New(ErrSmtpSetting)
	}
	if (len(stn.SmtpSettings.User) == 0) || (len(stn.SmtpSettings.UserAgent) == 0) {
		return errors.New(ErrSmtpSetting)
	}

	return nil
}

func getPasswordFromStdin() (pwd string, err error) {
	fmt.Println(c.MsgEnterSmtpPassword)

	var buf []byte
	buf, err = term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}

	return string(buf), nil
}
