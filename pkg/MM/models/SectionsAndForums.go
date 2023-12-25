package models

type SectionsAndForums struct {
	Sections []Section `json:"sections"`
	Forums   []Forum   `json:"forums"`
}

func NewSectionsAndForums() (saf *SectionsAndForums) {
	saf = &SectionsAndForums{
		Sections: make([]Section, 0),
		Forums:   make([]Forum, 0),
	}

	return saf
}
