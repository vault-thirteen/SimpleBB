package settings

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

	HttpSettings    `json:"http"`
	CaptchaSettings `json:"captcha"`
}

// HttpSettings are settings of an HTTP server for incoming requests.
type HttpSettings = cs.HttpSettings

// CaptchaSettings are parameters of the captcha.
type CaptchaSettings struct {
	// StoreImages flag configures the output of created images. When it is set
	// to True, Captcha manager stores an image on a storage device and
	// responds only with the task ID. When set to False, Captcha manager does
	// not store an image on a storage device and responds with both image
	// binary data and task ID.
	StoreImages bool `json:"storeImages"`

	// ImagesFolder sets the storage folder for created images. This setting is
	// used together with the 'StoreImages' flag.
	ImagesFolder string `json:"imagesFolder"`

	// Dimensions of created images.
	ImageWidth  uint `json:"imageWidth"`
	ImageHeight uint `json:"imageHeight"`

	// Image's time to live, in seconds. Each image is deleted when this time
	// passes after its creation.
	ImageTTLSec uint `json:"imageTTLSec"`

	ClearImagesFolderAtStart bool `json:"clearImagesFolderAtStart"`

	// This setting allows to start an HTTP server for serving captcha saved
	// image files. This setting is used together with the 'StoreImages' flag.
	UseHttpServerForImages bool `json:"useHttpServerForImages"`

	// Parameters of the HTTP server serving saved image files.
	HttpServerHost string `json:"httpServerHost"`
	HttpServerPort uint16 `json:"httpServerPort"`
	HttpServerName string `json:"httpServerName"`
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

	return stn, nil
}

func (stn *Settings) Check() (err error) {
	err = cs.CheckSettingsFilePath(stn.FilePath)
	if err != nil {
		return err
	}

	// HTTP.
	err = cs.CheckHttpSettings(stn.HttpSettings)
	if err != nil {
		return err
	}

	// Captcha.
	if stn.CaptchaSettings.StoreImages == true {
		if len(stn.CaptchaSettings.ImagesFolder) == 0 {
			return errors.New(c.MsgCaptchaSettingError)
		}
	}
	if (stn.CaptchaSettings.ImageWidth == 0) ||
		(stn.CaptchaSettings.ImageHeight == 0) ||
		(stn.CaptchaSettings.ImageTTLSec == 0) {
		return errors.New(c.MsgCaptchaSettingError)
	}

	if stn.CaptchaSettings.UseHttpServerForImages == true {
		if (len(stn.CaptchaSettings.HttpServerHost) == 0) ||
			(stn.CaptchaSettings.HttpServerPort == 0) ||
			(len(stn.CaptchaSettings.HttpServerName) == 0) {
			return errors.New(c.MsgCaptchaSettingError)
		}
	}

	return nil
}
