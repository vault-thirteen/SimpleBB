package models

type UserData struct {
	User    *User
	Session *Session
}

func NewUserData() (ud *UserData) {
	return &UserData{
		User:    &User{},
		Session: &Session{},
	}
}
