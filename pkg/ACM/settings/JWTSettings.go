package s

import (
	"errors"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
)

// JWTSettings are settings for JSON web tokens.
type JWTSettings struct {
	PrivateKeyFilePath string `json:"privateKeyFilePath"`
	PublicKeyFilePath  string `json:"publicKeyFilePath"`
	SigningMethod      string `json:"signingMethod"`
}

func (s JWTSettings) Check() (err error) {
	if (len(s.PrivateKeyFilePath) == 0) ||
		(len(s.PublicKeyFilePath) == 0) ||
		(len(s.SigningMethod) == 0) {
		return errors.New(c.MsgJwtSettingError)
	}

	return nil
}
