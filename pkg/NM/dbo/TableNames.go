package dbo

const (
	TableIncidents     = "Incidents"
	TableNotifications = "Notifications"
	TableSystemEvents  = "SystemEvents"
)

type TableNames struct {
	Incidents     string
	Notifications string
	SystemEvents  string
}
