package main

import (
	"log"

	r "github.com/vault-thirteen/SimpleBB/pkg/RCS"
	rs "github.com/vault-thirteen/SimpleBB/pkg/RCS/settings"
	"github.com/vault-thirteen/SimpleBB/pkg/common/app"
)

func main() {
	theApp, err := app.NewApplication[*rs.Settings, *r.Server](&rs.Settings{}, &r.Server{}, app.ServiceName_RCS, app.ConfigurationFilePathDefault_RCS)
	mustBeNoError(err)

	err = theApp.Use()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
