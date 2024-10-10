package main

import (
	"log"

	s "github.com/vault-thirteen/SimpleBB/pkg/SM"
	ss "github.com/vault-thirteen/SimpleBB/pkg/SM/settings"
	"github.com/vault-thirteen/SimpleBB/pkg/common/app"
)

func main() {
	theApp, err := app.NewApplication[*ss.Settings, *s.Server](&ss.Settings{}, &s.Server{}, app.ServiceName_SM, app.ConfigurationFilePathDefault_SM)
	mustBeNoError(err)

	err = theApp.Use()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
