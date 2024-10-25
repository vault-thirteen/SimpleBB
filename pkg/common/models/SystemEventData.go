package models

import (
	"fmt"

	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

const (
	ErrSystemEventType       = "system event type error"
	ErrSystemEventParameters = "system event parameters error"
)

// SystemEventData is a set of parameters of a system event.
type SystemEventData struct {
	// Event type. This field is required. All other fields are optional.
	// The number and values of other fields depend on the event type.
	Type SystemEventType `json:"type"`

	// ID of the thread mentioned in the event.
	// E.g., if a thread was renamed, it's ID is this field.
	ThreadId *cmb.Id `json:"threadId"`

	// ID of the message mentioned in the event.
	// E.g., if a message was added into a thread, it's ID is this field.
	MessageId *cmb.Id `json:"messageId"`

	// ID of the user mentioned in the event.
	// E.g., if a user has changed name of a thread, it's ID is this field.
	UserId *cmb.Id `json:"userId"`

	// Auxiliary fields.

	// ID of a user who initially created the object.
	// Some events deal with two UserId's, where one is the ID of the original
	// creator of the object and another is the user who modified the object.
	// In such cases, the original creator of the object is stored in the
	// 'Creator' field and the user who modified the object is stored in the
	// 'UserId' field. In cases where only a single user is important, this
	// field is not used, as this only user is stored in the 'UserId' field.
	Creator *cmb.Id `json:"creator,omitempty"`
}

func (sed *SystemEventData) CheckType() (ok bool, err error) {
	if !sed.Type.IsValid() {
		return false, fmt.Errorf(ErrSystemEventType)
	}

	return true, nil
}

func (sed *SystemEventData) isThreadIdSet() (ok bool) {
	return (sed.ThreadId != nil) && (*sed.ThreadId > 0)
}

func (sed *SystemEventData) isMessageIdSet() (ok bool) {
	return (sed.MessageId != nil) && (*sed.MessageId > 0)
}

func (sed *SystemEventData) isUserIdSet() (ok bool) {
	return (sed.UserId != nil) && (*sed.UserId > 0)
}

func (sed *SystemEventData) isCreatorSet() (ok bool) {
	return (sed.Creator != nil) && (*sed.Creator > 0)
}

func (sed *SystemEventData) CheckParameters() (ok bool, err error) {
	// Default requirements.
	var req = SystemEventRequirements{
		IsThreadIdRequired: true,
		IsUserIdRequired:   true,
	}

	switch sed.Type {
	case SystemEventType_ThreadParentChange:
		// Default requirements are used (TU).

	case SystemEventType_ThreadNameChange:
		// Default requirements are used (TU).

	case SystemEventType_ThreadDeletion:
		// Default requirements are used (TU).

	case SystemEventType_ThreadNewMessage:
		// TMU.
		req.IsMessageIdRequired = true

	case SystemEventType_ThreadMessageEdit:
		// TMU.
		req.IsMessageIdRequired = true

	case SystemEventType_ThreadMessageDeletion:
		// TMU.
		req.IsMessageIdRequired = true

	case SystemEventType_MessageTextEdit:
		// TMUC.
		req.IsMessageIdRequired = true
		req.IsCreatorRequired = true

	case SystemEventType_MessageParentChange:
		// TMUC.
		req.IsMessageIdRequired = true
		req.IsCreatorRequired = true

	case SystemEventType_MessageDeletion:
		// TMUC.
		req.IsMessageIdRequired = true
		req.IsCreatorRequired = true

	default:
		return false, fmt.Errorf(ErrSystemEventType)
	}

	// Check the required parameters.
	if req.IsThreadIdRequired {
		if !sed.isThreadIdSet() {
			return false, fmt.Errorf(ErrSystemEventParameters)
		}
	}

	if req.IsMessageIdRequired {
		if !sed.isMessageIdSet() {
			return false, fmt.Errorf(ErrSystemEventParameters)
		}
	}

	if req.IsUserIdRequired {
		if !sed.isUserIdSet() {
			return false, fmt.Errorf(ErrSystemEventParameters)
		}
	}

	if req.IsCreatorRequired {
		if !sed.isCreatorSet() {
			return false, fmt.Errorf(ErrSystemEventParameters)
		}
	}

	return true, nil
}
