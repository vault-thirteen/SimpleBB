package s

import (
	"errors"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
)

// SystemSettings are system settings.
type SystemSettings struct {
	SettingsVersion uint   `json:"settingsVersion" :"settings_version"`
	SiteName        string `json:"siteName" :"site_name"`
	SiteDomain      string `json:"siteDomain" :"site_domain"`

	// This setting must be synchronised with settings of the ACM module.
	IsFirewallUsed bool `json:"isFirewallUsed" :"is_firewall_used"`

	// ClientIPAddressSource setting selects where to search for client's IP
	// address. '1' means that IP address is taken directly from the client's
	// address of the HTTP request; '2' means that IP address is taken from the
	// custom HTTP header which is configured by the ClientIPAddressHeader
	// setting. One of the most common examples of a custom header may be the
	// 'X-Forwarded-For' HTTP header. For most users the first variant ('1') is
	// the most suitable. The second variant ('2') may be used if you are
	// proxying requests of your clients somewhere inside your own network
	// infrastructure, such as via a load balancer or with a reverse proxy.
	ClientIPAddressSource byte   `json:"clientIPAddressSource" :"client_ip_address_source"`
	ClientIPAddressHeader string `json:"clientIPAddressHeader" :"client_ip_address_header"`

	// URL.
	IsFrontEndEnabled bool   `json:"isFrontEndEnabled" :"is_front_end_enabled"`
	FrontEndPath      string `json:"frontEndPath" :"front_end_path"`
	ApiPath           string `json:"apiPath" :"api_path"`
	CaptchaPath       string `json:"captchaPath" :"captcha_path"`

	// Captcha.
	// These settings must be synchronised with settings of the RCS module.
	CaptchaImgServerHost   string `json:"captchaImgServerHost" :"captcha_img_server_host"`
	CaptchaImgServerPort   uint16 `json:"captchaImgServerPort" :"captcha_img_server_port"`
	SessionMaxDuration     uint   `json:"sessionMaxDuration" :"session_max_duration"`
	MessageEditTime        uint   `json:"messageEditTime" :"message_edit_time"`
	PageSize               uint   `json:"pageSize" :"page_size"`
	PublicSettingsFileName string `json:"publicSettingsFileName" :"public_settings_file_name"`
	FrontendAssetsFolder   string `json:"frontendAssetsFolder" :"frontend_assets_folder"`

	IsDebugMode bool `json:"isDebugMode" :"is_debug_mode"`
}

func (s SystemSettings) Check() (err error) {
	if (s.SettingsVersion == 0) ||
		(len(s.SiteName) == 0) ||
		(len(s.SiteDomain) == 0) ||
		(s.ClientIPAddressSource < ClientIPAddressSource_Direct) ||
		(s.ClientIPAddressSource > ClientIPAddressSource_CustomHeader) ||
		(len(s.ApiPath) == 0) ||
		(len(s.CaptchaPath) == 0) ||
		(s.ApiPath == s.CaptchaPath) ||
		(len(s.CaptchaImgServerHost) == 0) ||
		(s.CaptchaImgServerPort == 0) ||
		(s.SessionMaxDuration == 0) ||
		(s.MessageEditTime == 0) ||
		(s.PageSize == 0) ||
		(len(s.PublicSettingsFileName) == 0) ||
		(len(s.FrontendAssetsFolder) == 0) {
		return errors.New(c.MsgSystemSettingError)
	}

	if s.IsFrontEndEnabled {
		if (len(s.FrontEndPath) == 0) ||
			(s.FrontEndPath == s.ApiPath) ||
			(s.FrontEndPath == s.CaptchaPath) {
			return errors.New(c.MsgSystemSettingError)
		}
	}

	if s.ClientIPAddressSource == ClientIPAddressSource_CustomHeader {
		if len(s.ClientIPAddressHeader) == 0 {
			return errors.New(c.MsgSystemSettingError)
		}
	}

	return nil
}
