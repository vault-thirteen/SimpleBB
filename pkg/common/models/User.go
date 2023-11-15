package models

import "time"

type UserParameters struct {
	Id           uint      `json:"id"`
	PreRegTime   time.Time `json:"preRegTime"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	ApprovalTime time.Time `json:"approvalTime"`
	RegTime      time.Time `json:"regTime"`
	UserRoles
	LastBadLogInTime *time.Time `json:"lastBadLogInTime"`
	BanTime          *time.Time `json:"banTime"`
}

type UserRoles struct {
	IsAdministrator bool `json:"isAdministrator"`
	IsModerator     bool `json:"isModerator"`
	IsAuthor        bool `json:"isAuthor"`
	IsWriter        bool `json:"isWriter"`
	IsReader        bool `json:"isReader"`
	CanLogIn        bool `json:"canLogIn"`
}
