package s

// UserRoleSettings are settings for special user roles.
type UserRoleSettings struct {
	// List of IDs of users having a moderator role.
	ModeratorIds []uint `json:"moderatorIds"`

	// List of IDs of users having an administrator role.
	AdministratorIds []uint `json:"administratorIds"`
}

func (s UserRoleSettings) Check() (err error) {
	return nil
}
