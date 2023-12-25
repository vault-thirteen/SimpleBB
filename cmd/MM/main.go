package main

import (
	"log"

	m "github.com/vault-thirteen/SimpleBB/pkg/MM"
	ms "github.com/vault-thirteen/SimpleBB/pkg/MM/settings"
	"github.com/vault-thirteen/SimpleBB/pkg/common/app"
)

func main() {
	theApp, err := app.NewApplication[*ms.Settings, *m.Server](&ms.Settings{}, &m.Server{}, app.ServiceName_MM, app.ConfigurationFilePathDefault_MM)
	mustBeNoError(err)

	err = theApp.Use()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
