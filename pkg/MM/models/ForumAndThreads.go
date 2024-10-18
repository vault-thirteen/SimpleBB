package models

import (
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

type ForumAndThreads struct {
	Forum    *Forum        `json:"forum"`
	Threads  []Thread      `json:"threads"`
	PageData *cmr.PageData `json:"pageData,omitempty"`
}
