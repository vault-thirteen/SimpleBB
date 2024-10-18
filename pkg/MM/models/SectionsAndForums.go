package models

type SectionsAndForums struct {
	Sections []Section `json:"sections"`
	Forums   []Forum   `json:"forums"`
}
