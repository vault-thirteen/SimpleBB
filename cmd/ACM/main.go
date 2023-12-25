package main

import (
	"log"

	a "github.com/vault-thirteen/SimpleBB/pkg/ACM"
	as "github.com/vault-thirteen/SimpleBB/pkg/ACM/settings"
	"github.com/vault-thirteen/SimpleBB/pkg/common/app"
)

func main() {
	theApp, err := app.NewApplication[*as.Settings, *a.Server](&as.Settings{}, &a.Server{}, app.ServiceName_ACM, app.ConfigurationFilePathDefault_ACM)
	mustBeNoError(err)

	err = theApp.Use()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
