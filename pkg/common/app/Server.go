package app

import (
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/models/interfaces"
)

func NewServer[T cmi.IServer](classSelector T, settings cmi.ISettings) (srv cmi.IServer, err error) {
	return classSelector.UseConstructor(settings)
}
