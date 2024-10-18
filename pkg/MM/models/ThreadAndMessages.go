package models

import (
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

type ThreadAndMessages struct {
	Thread   *Thread       `json:"thread"`
	Messages []Message     `json:"messages"`
	PageData *cmr.PageData `json:"pageData"`
}
