package app

import (
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

func NewServer[T cm.IServer](classSelector T, settings cm.ISettings) (srv cm.IServer, err error) {
	return classSelector.UseConstructor(settings)
}
