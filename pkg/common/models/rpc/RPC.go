package rpc

import "net"

type CommonParams struct {
	// Authorisation data.
	Auth *Auth `json:"auth"`
}

type CommonResult struct {
	// Time taken to perform the request, in milliseconds.
	TimeSpent int64 `json:"timeSpent,omitempty"`
}

// Auth is authorisation data.
// This field must always be set for all requests except 'Ping'.
type Auth struct {
	// User's IP address, in a common textual format.
	// This field must always be set for all requests.
	UserIPA string `json:"userIPA"`

	// Parsed IP address of a user.	B = Byte array.
	// This field is used only in internal communications.
	UserIPAB net.IP `json:"-"`

	// Authorisation token.
	// This field is set for all requests of a logged user.
	Token string `json:"token"`
}

type PingParams struct{}

type PingResult struct {
	OK bool `json:"ok"`
}