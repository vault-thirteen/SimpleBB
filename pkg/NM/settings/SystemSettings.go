package s

import (
	"errors"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

// SystemSettings are system settings.
type SystemSettings struct {
	NotificationTtl        cmb.Count `json:"notificationTtl"`
	NotificationCountLimit cmb.Count `json:"notificationCountLimit"`
	PageSize               cmb.Count `json:"pageSize"`
	DKeySize               cmb.Count `json:"dKeySize"`

	// This setting must be synchronised with settings of the Gateway module.
	IsTableOfIncidentsUsed cmb.Flag `json:"isTableOfIncidentsUsed"`

	// This setting is used only when a table of incidents is enabled.
	BlockTimePerIncident BlockTimePerIncident `json:"blockTimePerIncident"`

	IsDebugMode cmb.Flag `json:"isDebugMode"`
}

// BlockTimePerIncident is block time in seconds for each type of incident.
type BlockTimePerIncident struct {
	IllegalAccessAttempt            cmb.Count `json:"illegalAccessAttempt"`            // 1.
	ReadingNotificationOfOtherUsers cmb.Count `json:"readingNotificationOfOtherUsers"` // 2.
	WrongDKey                       cmb.Count `json:"wrongDKey"`                       // 3.
}

func (s SystemSettings) Check() (err error) {
	if (s.NotificationTtl == 0) ||
		(s.NotificationCountLimit == 0) ||
		(s.PageSize == 0) ||
		(s.DKeySize == 0) {
		return errors.New(c.MsgSystemSettingError)
	}

	// Incidents.
	if s.IsTableOfIncidentsUsed {
		if (s.BlockTimePerIncident.IllegalAccessAttempt == 0) ||
			(s.BlockTimePerIncident.ReadingNotificationOfOtherUsers == 0) ||
			(s.BlockTimePerIncident.WrongDKey == 0) {
			return errors.New(c.MsgSystemSettingError)
		}
	}

	return nil
}
