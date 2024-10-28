package models

const (
	NotificationTemplate_R    = NotificationTemplate("R")
	NotificationTemplate_FT   = NotificationTemplate("FT")
	NotificationTemplate_FRT  = NotificationTemplate("FRT")
	NotificationTemplate_FUT  = NotificationTemplate("FUT")
	NotificationTemplate_FUMT = NotificationTemplate("FUMT")
	NotificationTemplate_FMTU = NotificationTemplate("FMTU")
)

// NotificationTemplate is a template of a notification.
// A template describes contents of a notification and stores links to other
// sources of information.
type NotificationTemplate string

func NewNotificationTemplate(s string) (nt *NotificationTemplate) {
	x := NotificationTemplate(s)
	if !x.IsValid() {
		return nil
	}
	return &x
}

func (nt NotificationTemplate) IsValid() bool {
	switch nt {
	case NotificationTemplate_R,
		NotificationTemplate_FT,
		NotificationTemplate_FRT,
		NotificationTemplate_FUT,
		NotificationTemplate_FUMT,
		NotificationTemplate_FMTU:
		return true

	default:
		return false
	}
}

func (nt NotificationTemplate) ToString() string {
	return string(nt)
}
