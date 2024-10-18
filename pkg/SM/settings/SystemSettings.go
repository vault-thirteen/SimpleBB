package s

import (
	"errors"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

// SystemSettings are system settings.
type SystemSettings struct {
	PageSize    cmb.Count `json:"pageSize"`
	DKeySize    cmb.Count `json:"dKeySize"`
	IsDebugMode cmb.Flag  `json:"isDebugMode"`
}

func (s SystemSettings) Check() (err error) {
	if (s.PageSize == 0) ||
		(s.DKeySize == 0) {
		return errors.New(c.MsgSystemSettingError)
	}

	return nil
}
