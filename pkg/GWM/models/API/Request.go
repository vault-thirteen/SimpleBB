package api

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

// Request is an API request model. It is a mixture of client data with a data
// set received from the client. This object is used for API function calls
// between services.
type Request struct {
	Action        string
	Parameters    any
	Authorisation *models.Auth
}
