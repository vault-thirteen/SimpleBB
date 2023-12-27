package s

import (
	"errors"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
)

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
	IsDebugMode       bool   `json:"isDebugMode"`
}

func (s SystemSettings) Check() (err error) {
	if (len(s.SiteName) == 0) ||
		(len(s.SiteDomain) == 0) ||
		(s.ClientIPAddressSource < ClientIPAddressSource_Direct) ||
		(s.ClientIPAddressSource > ClientIPAddressSource_CustomHeader) ||
		(len(s.FrontEndPath) == 0) ||
		(len(s.ApiPath) == 0) ||
		(s.FrontEndPath == s.ApiPath) {
		return errors.New(c.MsgSystemSettingError)
	}

	if s.ClientIPAddressSource == ClientIPAddressSource_CustomHeader {
		if len(s.ClientIPAddressHeader) == 0 {
			return errors.New(c.MsgSystemSettingError)
		}
	}

	return nil
}
