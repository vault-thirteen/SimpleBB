package models

import (
	"database/sql"
	"errors"

	"github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

// ThreadLink is a short variant of a thread which stores only IDs.
type ThreadLink struct {
	// Identifier of this thread.
	Id cmb.Id `json:"id"`

	// Identifier of a forum containing this thread.
	ForumId cmb.Id `json:"forumId"`

	// List of identifiers of messages of this thread.
	Messages *ul.UidList `json:"messages"`
}

func NewThreadLink() (tl *ThreadLink) {
	return &ThreadLink{}
}

func NewThreadLinkFromScannableSource(src cm.IScannable) (tl *ThreadLink, err error) {
	tl = NewThreadLink()
	var x = ul.New()

	err = src.Scan(
		&tl.Id,
		&tl.ForumId,
		x, //&tl.Messages,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	tl.Messages = x
	return tl, nil
}
