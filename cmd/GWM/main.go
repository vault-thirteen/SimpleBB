package main

import (
	"log"

	g "github.com/vault-thirteen/SimpleBB/pkg/GWM"
	gs "github.com/vault-thirteen/SimpleBB/pkg/GWM/settings"
	"github.com/vault-thirteen/SimpleBB/pkg/common/app"
)

func main() {
	theApp, err := app.NewApplication[*gs.Settings, *g.Server](&gs.Settings{}, &g.Server{}, app.ServiceName_GWM, app.ConfigurationFilePathDefault_GWM)
	mustBeNoError(err)

	err = theApp.Use()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
