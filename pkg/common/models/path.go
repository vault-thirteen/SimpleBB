package models

import "path/filepath"

func NormalisePath(path string) string {
	return filepath.FromSlash(filepath.ToSlash(path))
}
