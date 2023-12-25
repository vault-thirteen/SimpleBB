package models

type IServer interface {
	Start() error
	ReportStart()
	Stop() error
	GetStopChannel() *chan bool
	UseConstructor(ISettings) (IServer, error)
}
