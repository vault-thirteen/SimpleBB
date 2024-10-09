package models

import (
	"database/sql"
	"errors"

	"github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

type Thread struct {
	// Identifier of this thread.
	Id uint `json:"id"`

	// Identifier of a forum containing this thread.
	ForumId uint `json:"forumId"`

	// Name of this thread.
	Name string `json:"name"`

	// List of identifiers of messages of this thread.
	Messages *ul.UidList `json:"messages"`

	// Thread meta-data.
	EventData
}

func NewThread() (thr *Thread) {
	return &Thread{
		EventData: EventData{
			Creator: &EventParameters{},
			Editor:  &OptionalEventParameters{},
		},
	}
}

func NewThreadFromScannableSource(src cm.IScannable) (thr *Thread, err error) {
	thr = NewThread()
	var x = ul.New()

	err = src.Scan(
		&thr.Id,
		&thr.ForumId,
		&thr.Name,
		x, //&thr.Messages,
		&thr.Creator.UserId,
		&thr.Creator.Time,
		&thr.Editor.UserId,
		&thr.Editor.Time,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	thr.Messages = x
	return thr, nil
}
