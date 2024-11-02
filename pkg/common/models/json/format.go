package json

import (
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/models/interfaces"
)

func ToJson(x cmi.IToString) []byte {
	return []byte(x.ToString())
}
