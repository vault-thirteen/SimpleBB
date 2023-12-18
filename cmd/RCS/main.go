package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	r "github.com/vault-thirteen/SimpleBB/pkg/RCS"
	rs "github.com/vault-thirteen/SimpleBB/pkg/RCS/settings"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	"github.com/vault-thirteen/Versioneer"
)

const (
	ServiceName                  = "Captcha Module"
	ConfigurationFilePathDefault = "RCS.json"
)

func main() {
	showIntro()

	cla, err := readCLA()
	mustBeNoError(err)
	if cla.IsDefaultFile() {
		log.Println(c.MsgUsingDefaultConfigurationFile)
	}

	var stn *rs.Settings
	stn, err = rs.NewSettingsFromFile(cla.ConfigurationFilePath)
	mustBeNoError(err)

	log.Println(c.MsgServerIsStarting)
	var srv *r.Server
	srv, err = r.NewServer(stn)
	mustBeNoError(err)

	err = srv.Start()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c.MsgRpcHttpServer + srv.GetListenDsn())
	if stn.CaptchaSettings.UseHttpServerForImages {
		fmt.Println(c.MsgImagesHttpServer + srv.GetCaptchaManagerListenDsn())
	}

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
