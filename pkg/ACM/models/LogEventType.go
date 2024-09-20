package models

const (
	LogEventType_LogIn   = 1
	LogEventType_LogOut  = 2 // Self logging out.
	LogEventType_LogOutA = 3 // Logging out by an administrator.
)

const (
	LogEventTypesCount = 3
)

type LogEventType byte

func (let LogEventType) IsValid() (ok bool) {
	if (let == 0) || (let > LogEventTypesCount) {
		return false
	}

	return true
}
