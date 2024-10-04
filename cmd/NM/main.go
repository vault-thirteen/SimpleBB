package main

import (
	"log"

	n "github.com/vault-thirteen/SimpleBB/pkg/NM"
	ns "github.com/vault-thirteen/SimpleBB/pkg/NM/settings"
	"github.com/vault-thirteen/SimpleBB/pkg/common/app"
)

func main() {
	theApp, err := app.NewApplication[*ns.Settings, *n.Server](&ns.Settings{}, &n.Server{}, app.ServiceName_NM, app.ConfigurationFilePathDefault_NM)
	mustBeNoError(err)

	err = theApp.Use()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
