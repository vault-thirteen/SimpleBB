package s

import (
	"errors"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
)

// SystemSettings are system settings.
type SystemSettings struct {
	DKeySize    uint `json:"dKeySize"`
	IsDebugMode bool `json:"isDebugMode"`
}

func (s SystemSettings) Check() (err error) {
	if s.DKeySize == 0 {
		return errors.New(c.MsgSystemSettingError)
	}

	return nil
}
