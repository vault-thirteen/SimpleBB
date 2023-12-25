package settings

// SystemSettings are system settings.
type SystemSettings struct {
	IsDebugMode bool `json:"isDebugMode"`
}

func (s SystemSettings) Check() (err error) {
	return nil
}
