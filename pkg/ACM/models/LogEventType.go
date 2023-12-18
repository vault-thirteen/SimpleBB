package models

const (
	LogEventType_LogIn  = 1
	LogEventType_LogOut = 2
)

const (
	LogEventTypesCount = 2
)

type LogEventType byte

func (let LogEventType) IsValid() (ok bool) {
	if (let == 0) || (let > LogEventTypesCount) {
		return false
	}

	return true
}
