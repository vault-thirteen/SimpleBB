package models

import (
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

type User struct {
	Password string `json:"-"`
	cm.UserParameters
}
