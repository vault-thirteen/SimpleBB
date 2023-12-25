package app

import (
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	ver "github.com/vault-thirteen/auxie/Versioneer"
)

func NewSettingsFromFile[T cm.ISettings](classSelector T, filePath string, versionInfo *ver.Versioneer) (stn cm.ISettings, err error) {
	return classSelector.UseConstructor(filePath, versionInfo)
}
