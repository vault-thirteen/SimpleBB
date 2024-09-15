package models

type Settings struct {
	Version                   uint   `json:"version"`
	ProductVersion            string `json:"productVersion"`
	SiteName                  string `json:"siteName"`
	SiteDomain                string `json:"siteDomain"`
	CaptchaFolder             string `json:"captchaFolder"`
	SessionMaxDuration        uint   `json:"sessionMaxDuration"`
	MessageEditTime           uint   `json:"messageEditTime"`
	PageSize                  uint   `json:"pageSize"`
	ApiFolder                 string `json:"apiFolder"`
	PublicSettingsFileName    string `json:"publicSettingsFileName"`
	IsFrontEndEnabled         bool   `json:"isFrontEndEnabled"`
	FrontEndStaticFilesFolder string `json:"frontEndStaticFilesFolder"`
}
