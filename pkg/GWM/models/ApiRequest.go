package models

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

// ApiRequest is a mixture of client data with a data set received from the
// client. This object is used for API function calls between services.
type ApiRequest struct {
	Action        string
	Parameters    any
	Authorisation *models.Auth
}
