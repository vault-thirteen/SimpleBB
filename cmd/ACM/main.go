package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	a "github.com/vault-thirteen/SimpleBB/pkg/ACM"
	as "github.com/vault-thirteen/SimpleBB/pkg/ACM/settings"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	"github.com/vault-thirteen/Versioneer"
)

const (
	ServiceName                  = "Access Control Module"
	ConfigurationFilePathDefault = "ACM.json"
)

func main() {
	showIntro()

	cla, err := readCLA()
	mustBeNoError(err)
	if cla.IsDefaultFile() {
		log.Println(c.MsgUsingDefaultConfigurationFile)
	}

	var stn *as.Settings
	stn, err = as.NewSettingsFromFile(cla.ConfigurationFilePath)
	mustBeNoError(err)

	log.Println(c.MsgServerIsStarting)
	var srv *a.Server
	srv, err = a.NewServer(stn)
	mustBeNoError(err)

	err = srv.Start()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c.MsgHttpsServer + srv.GetListenDsn())

	serverMustBeStopped := srv.GetStopChannel()
	waitForQuitSignalFromOS(serverMustBeStopped)
	<-*serverMustBeStopped

	log.Println(c.MsgServerIsStopping)
	err = srv.Stop()
	if err != nil {
		log.Println(err)
	}
	log.Println(c.MsgServerIsStopped)
	time.Sleep(time.Second)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func showIntro() {
	versioneer, err := ver.New()
	mustBeNoError(err)
	versioneer.ShowIntroText(ServiceName)
	versioneer.ShowComponentsInfoText()
	fmt.Println()
}

func waitForQuitSignalFromOS(serverMustBeStopped *chan bool) {
	osSignals := make(chan os.Signal, 16)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for sig := range osSignals {
			switch sig {
			case syscall.SIGINT,
				syscall.SIGTERM:
				log.Println(c.MsgQuitSignalIsReceived, sig)
				*serverMustBeStopped <- true
			}
		}
	}()
}
