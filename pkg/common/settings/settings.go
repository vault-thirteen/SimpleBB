package s

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/term"
)

const (
	ErrSettingsFileIsNotSet = "settings file is not set"
)

const (
	RpcDurationFieldName  = "dur"
	RpcRequestIdFieldName = "rid"
)

func GetPasswordFromStdin(hint string) (pwd string, err error) {
	fmt.Println(hint)

	var buf []byte
	buf, err = term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func CheckSettingsFilePath(sfp string) (err error) {
	if len(sfp) == 0 {
		return errors.New(ErrSettingsFileIsNotSet)
	}

	return nil
}
