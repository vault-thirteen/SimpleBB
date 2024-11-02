package interfaces

type IScannable interface {
	Scan(dest ...any) error
}
