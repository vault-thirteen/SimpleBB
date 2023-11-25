package dbo

const (
	TablePreRegisteredUsers = "PreRegisteredUsers"
	TableUsers              = "Users"
	TablePreSessions        = "PreSessions"
	TableSessions           = "Sessions"
	TableIncidents          = "Incidents"
)

type TableNames struct {
	PreRegisteredUsers string
	Users              string
	PreSessions        string
	Sessions           string
	Incidents          string
}
