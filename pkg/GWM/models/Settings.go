package models

import (
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type Settings struct {
	Version                   cmb.Count `json:"version"`
	ProductVersion            cmb.Text  `json:"productVersion"`
	SiteName                  cmb.Text  `json:"siteName"`
	SiteDomain                cmb.Text  `json:"siteDomain"`
	CaptchaFolder             cm.Path   `json:"captchaFolder"`
	SessionMaxDuration        cmb.Count `json:"sessionMaxDuration"`
	MessageEditTime           cmb.Count `json:"messageEditTime"`
	PageSize                  cmb.Count `json:"pageSize"`
	ApiFolder                 cm.Path   `json:"apiFolder"`
	PublicSettingsFileName    cm.Path   `json:"publicSettingsFileName"`
	IsFrontEndEnabled         cmb.Flag  `json:"isFrontEndEnabled"`
	FrontEndStaticFilesFolder cm.Path   `json:"frontEndStaticFilesFolder"`
	NotificationCountLimit    cmb.Count `json:"notificationCountLimit"`
}
