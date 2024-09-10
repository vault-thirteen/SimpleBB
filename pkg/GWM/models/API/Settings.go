package api

type Settings struct {
	Version            uint   `json:"version"`
	ProductVersion     string `json:"productVersion"`
	SiteName           string `json:"siteName"`
	SiteDomain         string `json:"siteDomain"`
	IsFrontEndEnabled  bool   `json:"isFrontEndEnabled"`
	FrontEndPath       string `json:"frontEndPath"`
	ApiPath            string `json:"apiPath"`
	CaptchaPath        string `json:"captchaPath"`
	SessionMaxDuration uint   `json:"sessionMaxDuration"`
	MessageEditTime    uint   `json:"messageEditTime"`
	PageSize           uint   `json:"pageSize"`
}
