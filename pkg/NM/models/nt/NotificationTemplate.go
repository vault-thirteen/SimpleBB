package nt

import (
	nm "github.com/vault-thirteen/SimpleBB/pkg/NM/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type NotificationTemplate struct {
	Id           cmb.Id                         `json:"id"`
	Name         cmb.Text                       `json:"name"`
	FormatString *nm.FormatString               `json:"formatString"`
	Arguments    *NotificationTemplateArguments `json:"arguments"`
}
