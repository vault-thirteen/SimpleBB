package im

import (
	"errors"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vault-thirteen/SimpleBB/pkg/ACM/dbo"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
)

const (
	TaskChannelSize = 4
)

const (
	ErrAlreadyStarted = "incident manager is already started"
	ErrAlreadyStopped = "incident manager is already stopped"
)

type IncidentManager struct {
	isWorking              atomic.Bool
	wg                     *sync.WaitGroup
	tasks                  chan *am.Incident
	isTableOfIncidentsUsed bool
	dbo                    *dbo.DatabaseObject
}

func NewIncidentManager(isTableOfIncidentsUsed bool, dbo *dbo.DatabaseObject) (im *IncidentManager) {
	im = &IncidentManager{
		wg:                     new(sync.WaitGroup),
		tasks:                  make(chan *am.Incident, TaskChannelSize),
		isTableOfIncidentsUsed: isTableOfIncidentsUsed,
		dbo:                    dbo,
	}

	im.isWorking.Store(false)

	return im
}

// Start starts the incident manager.
func (im *IncidentManager) Start() (err error) {
	if im.isWorking.Load() {
		return errors.New(ErrAlreadyStarted)
	}

	im.wg.Add(1)
	go im.run()
	im.isWorking.Store(true)

	return nil
}

// run is the main work loop of the incident manager.
func (im *IncidentManager) run() {
	defer im.wg.Done()

	var err error
	for inc := range im.tasks {
		if im.isTableOfIncidentsUsed {
			err = im.saveIncident(inc)
			if err != nil {
				log.Println(err)
			}
		}
	}

	log.Println(c.MsgIncidentManagerHasStopped)
}

// Stop stops the incident manager.
func (im *IncidentManager) Stop() (err error) {
	if !im.isWorking.Load() {
		return errors.New(ErrAlreadyStopped)
	}

	close(im.tasks)
	im.wg.Wait()
	im.isWorking.Store(false)

	return nil
}

func (im *IncidentManager) ReportIncident(itype am.IncidentType, email string, userIPA net.IP) {
	incident := &am.Incident{
		Time:    time.Now(),
		Type:    itype,
		Email:   email,
		UserIPA: userIPA,
	}

	im.tasks <- incident
}

func (im *IncidentManager) saveIncident(inc *am.Incident) (err error) {
	if inc == nil {
		return nil
	}

	if inc.UserIPA == nil {
		err = im.dbo.SaveIncidentWithoutUserIPA(inc.Type, inc.Email)
	} else {
		err = im.dbo.SaveIncident(inc.Type, inc.Email, inc.UserIPA)
	}
	if err != nil {
		return err
	}

	return nil
}
