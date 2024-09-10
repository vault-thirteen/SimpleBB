package models

import (
	"path/filepath"

	af "github.com/vault-thirteen/auxie/file"
)

// FrontEndFileData is auxiliary data for a front end static file.
type FrontEndFileData struct {
	UrlPath     string
	FilePath    string
	ContentType string
	CachedFile  []byte
}

func NewFrontEndFileData(frontEndPath string, fileName string, contentType string, frontendAssetsFolder string) (fefd FrontEndFileData, err error) {
	fefd = FrontEndFileData{
		UrlPath:     frontEndPath + fileName,
		FilePath:    filepath.Join(frontendAssetsFolder, fileName),
		ContentType: contentType,
	}

	fefd.CachedFile, err = af.GetFileContents(fefd.FilePath)
	if err != nil {
		return fefd, err
	}

	return fefd, nil
}
