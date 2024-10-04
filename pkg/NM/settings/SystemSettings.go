package s

import (
	"errors"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
)

// SystemSettings are system settings.
type SystemSettings struct {
	NotificationTtl        uint `json:"notificationTtl"`
	NotificationCountLimit uint `json:"notificationCountLimit"`

	// This setting must be synchronised with settings of the Gateway module.
	IsTableOfIncidentsUsed bool `json:"isTableOfIncidentsUsed"`

	// This setting is used only when a table of incidents is enabled.
	BlockTimePerIncident BlockTimePerIncident `json:"blockTimePerIncident"`

	IsDebugMode bool `json:"isDebugMode"`
}

// BlockTimePerIncident is block time in seconds for each type of incident.
type BlockTimePerIncident struct {
	IllegalAccessAttempt            uint `json:"illegalAccessAttempt"`            // 1.
	ReadingNotificationOfOtherUsers uint `json:"readingNotificationOfOtherUsers"` // 2.
}

func (s SystemSettings) Check() (err error) {
	if (s.NotificationTtl == 0) ||
		(s.NotificationCountLimit == 0) {
		return errors.New(c.MsgSystemSettingError)
	}

	// Incidents.
	if s.IsTableOfIncidentsUsed {
		if (s.BlockTimePerIncident.IllegalAccessAttempt == 0) ||
			(s.BlockTimePerIncident.ReadingNotificationOfOtherUsers == 0) {
			return errors.New(c.MsgSystemSettingError)
		}
	}

	return nil
}
