package models

import (
	"encoding/json"

	"github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

// ApiRequestWithOnlyAction is a data set for an API request received from a
// client in JSON format. This data is later mixed with additional client data.
type ApiRequestWithOnlyAction struct {
	// Action is a name of the function which a client wants to perform.
	Action string `json:"action"`

	// Parameters are function parameters.
	Parameters json.RawMessage `json:"parameters"`
}

// ApiRequest is a mixture of client data with a data set received from the
// client. This object is used for API function calls between services.
type ApiRequest struct {
	Action        string
	Parameters    any
	Authorisation *models.Auth
}
