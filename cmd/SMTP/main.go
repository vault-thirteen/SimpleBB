package main

import (
	"log"

	s "github.com/vault-thirteen/SimpleBB/pkg/SMTP"
	ss "github.com/vault-thirteen/SimpleBB/pkg/SMTP/settings"
	"github.com/vault-thirteen/SimpleBB/pkg/common/app"
)

func main() {
	theApp, err := app.NewApplication[*ss.Settings, *s.Server](&ss.Settings{}, &s.Server{}, app.ServiceName_SMTP, app.ConfigurationFilePathDefault_SMTP)
	mustBeNoError(err)

	err = theApp.Use()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
