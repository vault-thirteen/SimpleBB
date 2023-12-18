package models

type IScannable interface {
	Scan(dest ...any) error
}
