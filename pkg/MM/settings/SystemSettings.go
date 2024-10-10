package s

import (
	"errors"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
)

// SystemSettings are system settings.
type SystemSettings struct {
	DKeySize        uint `json:"dKeySize"`
	MessageEditTime uint `json:"messageEditTime"`
	PageSize        uint `json:"pageSize"`

	// NewThreadsAtTop parameter controls how new and updated threads are
	// placed inside forums. If set to 'True', then following will happen:
	// 1. New threads will be added to the start (top) of the list of forum's
	// threads instead of being added to the end (bottom) of the list;
	// 2. New messages added to threads will update the thread moving it to the
	// start (top) position of the list of forum's threads.
	// If set to 'False', then new threads are added to the end (bottom) of the
	// list and thread's new messages do not update thread's position in the
	// list.
	NewThreadsAtTop bool `json:"newThreadsAtTop"`

	IsDebugMode bool `json:"isDebugMode"`
}

func (s SystemSettings) Check() (err error) {
	if (s.DKeySize == 0) ||
		(s.MessageEditTime == 0) ||
		(s.PageSize == 0) {
		return errors.New(c.MsgSystemSettingError)
	}

	return nil
}
