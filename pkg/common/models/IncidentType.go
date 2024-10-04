package models

const (
	IncidentType_IllegalAccessAttempt            = 1
	IncidentType_FakeToken                       = 2
	IncidentType_VerificationCodeMismatch        = 3
	IncidentType_DoubleLogInAttempt              = 4
	IncidentType_PreSessionHacking               = 5
	IncidentType_CaptchaAnswerMismatch           = 6
	IncidentType_PasswordMismatch                = 7
	IncidentType_PasswordChangeHacking           = 8
	IncidentType_EmailChangeHacking              = 9
	IncidentType_FakeIPA                         = 10
	IncidentType_ReadingNotificationOfOtherUsers = 11
)

const (
	IncidentTypesCount = 11
)

type IncidentType byte

func (it IncidentType) IsValid() (ok bool) {
	if (it == 0) || (it > IncidentTypesCount) {
		return false
	}

	return true
}
