package net

import (
	"errors"
	"net"
)

const (
	ErrParseError = "IP address parse error"
)

// ParseIPA parses a string into an IP address and returns an error on error.
// Golang's built-in parser does not return an error! What a shame.
func ParseIPA(s string) (ipa net.IP, err error) {
	ipa = net.ParseIP(s)

	if ipa == nil {
		return nil, errors.New(ErrParseError)
	}

	return ipa, nil
}
