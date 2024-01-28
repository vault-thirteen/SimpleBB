package api

type Settings struct {
	SiteName           string `json:"siteName"`
	SiteDomain         string `json:"siteDomain"`
	IsFrontEndEnabled  bool   `json:"isFrontEndEnabled"`
	FrontEndPath       string `json:"frontEndPath"`
	SessionMaxDuration uint   `json:"sessionMaxDuration"`
}
