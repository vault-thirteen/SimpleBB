package s

import (
	"errors"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

// JWTSettings are settings for JSON web tokens.
type JWTSettings struct {
	PrivateKeyFilePath cm.Path `json:"privateKeyFilePath"`
	PublicKeyFilePath  cm.Path `json:"publicKeyFilePath"`
	SigningMethod      string  `json:"signingMethod"`
}

func (s JWTSettings) Check() (err error) {
	if (len(s.PrivateKeyFilePath) == 0) ||
		(len(s.PublicKeyFilePath) == 0) ||
		(len(s.SigningMethod) == 0) {
		return errors.New(c.MsgJwtSettingError)
	}

	return nil
}
