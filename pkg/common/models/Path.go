package models

import (
	"path/filepath"

	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type Path = cmb.Text

func NormalisePath(path Path) Path {
	return Path(filepath.FromSlash(filepath.ToSlash(path.ToString())))
}
