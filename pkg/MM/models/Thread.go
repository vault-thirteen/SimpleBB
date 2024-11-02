package models

import (
	"database/sql"
	"errors"

	"github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/models/interfaces"
)

type Thread struct {
	// Identifier of this thread.
	Id cmb.Id `json:"id"`

	// Identifier of a forum containing this thread.
	ForumId cmb.Id `json:"forumId"`

	// Name of this thread.
	Name cm.Name `json:"name"`

	// List of identifiers of messages of this thread.
	Messages *ul.UidList `json:"messages"`

	// Thread meta-data.
	EventData
}

func NewThread() (t *Thread) {
	return &Thread{
		EventData: EventData{
			Creator: &EventParameters{},
			Editor:  &OptionalEventParameters{},
		},
	}
}

func NewThreadFromScannableSource(src cmi.IScannable) (t *Thread, err error) {
	t = NewThread()
	var x = ul.New()

	err = src.Scan(
		&t.Id,
		&t.ForumId,
		&t.Name,
		x, //&t.Messages,
		&t.Creator.UserId,
		&t.Creator.Time,
		&t.Editor.UserId,
		&t.Editor.Time,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	t.Messages = x
	return t, nil
}

func NewThreadArrayFromRows(rows cmi.IScannableSequence) (ts []Thread, err error) {
	ts = []Thread{}
	var t *Thread

	for rows.Next() {
		t, err = NewThreadFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		ts = append(ts, *t)
	}

	return ts, nil
}
