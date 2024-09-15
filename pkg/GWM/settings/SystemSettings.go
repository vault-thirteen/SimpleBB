package s

import (
	"errors"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
)

// SystemSettings are system settings.
// Many of these settings must be synchronised with other modules.
type SystemSettings struct {
	SettingsVersion uint   `json:"settingsVersion"`
	SiteName        string `json:"siteName"`
	SiteDomain      string `json:"siteDomain"`

	// Firewall.
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

	// Captcha.
	CaptchaImgServerHost string `json:"captchaImgServerHost"`
	CaptchaImgServerPort uint16 `json:"captchaImgServerPort"`
	CaptchaFolder        string `json:"captchaFolder"`

	// Sessions and messages.
	SessionMaxDuration uint `json:"sessionMaxDuration"`
	MessageEditTime    uint `json:"messageEditTime"`
	PageSize           uint `json:"pageSize"`

	// URL paths.
	ApiFolder              string `json:"apiFolder"`
	PublicSettingsFileName string `json:"publicSettingsFileName"`

	// Front end.
	IsFrontEndEnabled         bool   `json:"isFrontEndEnabled"`
	FrontEndStaticFilesFolder string `json:"frontEndStaticFilesFolder"`
	FrontEndAssetsFolder      string `json:"frontEndAssetsFolder"`

	// Development settings.
	IsDebugMode                               bool   `json:"isDebugMode"`
	IsDeveloperMode                           bool   `json:"isDeveloperMode"`
	DevModeHttpHeaderAccessControlAllowOrigin string `json:"devModeHttpHeaderAccessControlAllowOrigin"`
}

func (s SystemSettings) Check() (err error) {
	if (s.SettingsVersion == 0) ||
		(len(s.SiteName) == 0) ||
		(len(s.SiteDomain) == 0) ||
		(s.ClientIPAddressSource < ClientIPAddressSource_Direct) ||
		(s.ClientIPAddressSource > ClientIPAddressSource_CustomHeader) ||
		(len(s.CaptchaImgServerHost) == 0) ||
		(s.CaptchaImgServerPort == 0) ||
		(len(s.CaptchaFolder) == 0) ||
		(s.SessionMaxDuration == 0) ||
		(s.MessageEditTime == 0) ||
		(s.PageSize == 0) ||
		(len(s.ApiFolder) == 0) ||
		(len(s.PublicSettingsFileName) == 0) {
		return errors.New(c.MsgSystemSettingError)
	}

	if s.IsFrontEndEnabled {
		if (len(s.FrontEndStaticFilesFolder) == 0) ||
			(len(s.FrontEndAssetsFolder) == 0) {
			return errors.New(c.MsgSystemSettingError)
		}
	}

	if s.ClientIPAddressSource == ClientIPAddressSource_CustomHeader {
		if len(s.ClientIPAddressHeader) == 0 {
			return errors.New(c.MsgSystemSettingError)
		}
	}

	if s.IsDeveloperMode {
		if len(s.DevModeHttpHeaderAccessControlAllowOrigin) == 0 {
			return errors.New(c.MsgSystemSettingError)
		}
	}

	return nil
}
