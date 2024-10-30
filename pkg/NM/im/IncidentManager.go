package im

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	gc "github.com/vault-thirteen/SimpleBB/pkg/GWM/client"
	gm "github.com/vault-thirteen/SimpleBB/pkg/GWM/models"
	"github.com/vault-thirteen/SimpleBB/pkg/NM/dbo"
	s "github.com/vault-thirteen/SimpleBB/pkg/NM/settings"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	"github.com/vault-thirteen/SimpleBB/pkg/common/app"
	"github.com/vault-thirteen/SimpleBB/pkg/common/avm"
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/client"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

const (
	TaskChannelSize = 4
)

type IncidentManager struct {
	ssp                    *avm.SSP
	wg                     *sync.WaitGroup
	tasks                  chan *cm.Incident
	isTableOfIncidentsUsed bool
	dbo                    *dbo.DatabaseObject
	gwmClient              *cc.Client

	// Block time in seconds for each incident type.
	blockTimePerIncidentType [cm.IncidentTypeMax + 1]cmb.Count
}

func NewIncidentManager(
	isTableOfIncidentsUsed bool,
	dbo *dbo.DatabaseObject,
	gwmClient *cc.Client,
	blockTimePerIncident *s.BlockTimePerIncident,
) (im *IncidentManager) {
	im = &IncidentManager{
		ssp:                      avm.NewSSP(),
		wg:                       new(sync.WaitGroup),
		tasks:                    make(chan *cm.Incident, TaskChannelSize),
		isTableOfIncidentsUsed:   isTableOfIncidentsUsed,
		dbo:                      dbo,
		gwmClient:                gwmClient,
		blockTimePerIncidentType: initBlockTimePerIncidentType(blockTimePerIncident),
	}

	return im
}

func initBlockTimePerIncidentType(blockTimePerIncident *s.BlockTimePerIncident) (blockTimePerIncidentType [cm.IncidentTypeMax + 1]cmb.Count) {
	// The "zero"-indexed element is empty because it is not used.
	blockTimePerIncidentType[cm.IncidentType_IllegalAccessAttempt] = blockTimePerIncident.IllegalAccessAttempt
	blockTimePerIncidentType[cm.IncidentType_ReadingNotificationOfOtherUsers] = blockTimePerIncident.ReadingNotificationOfOtherUsers
	blockTimePerIncidentType[cm.IncidentType_WrongDKey] = blockTimePerIncident.WrongDKey

	return blockTimePerIncidentType
}

// Start starts the incident manager.
func (im *IncidentManager) Start() (err error) {
	im.ssp.Lock()
	defer im.ssp.Unlock()

	err = im.ssp.BeginStart()
	if err != nil {
		return err
	}

	im.wg.Add(1)
	go im.run()

	im.ssp.CompleteStart()

	return nil
}

// run is the main work loop of the incident manager.
func (im *IncidentManager) run() {
	defer im.wg.Done()

	var err error
	var re *jrm1.RpcError
	for inc := range im.tasks {
		if im.isTableOfIncidentsUsed {
			err = cm.CheckIncident(inc)
			if err != nil {
				log.Println(err)
				continue
			}

			err = im.saveIncident(inc)
			im.logError(err)

			re = im.informGateway(inc)
			// This is why Go language is a complete Scheiße (utter trash):
			// https://github.com/golang/go/issues/40442
			if re != nil {
				err = re.AsError()
			} else {
				err = nil
			}
			im.logError(err)
		}
	}

	log.Println(c.MsgIncidentManagerHasStopped)
}

// Stop stops the incident manager.
func (im *IncidentManager) Stop() (err error) {
	im.ssp.Lock()
	defer im.ssp.Unlock()

	err = im.ssp.BeginStop()
	if err != nil {
		return err
	}

	close(im.tasks)
	im.wg.Wait()

	im.ssp.CompleteStop()

	return nil
}

func (im *IncidentManager) ReportIncident(itype cmb.EnumValue, email cm.Email, userIPA net.IP) {
	incident := &cm.Incident{
		Time:    time.Now(),
		Type:    cm.NewIncidentTypeWithValue(itype),
		Email:   email,
		UserIPA: userIPA,
	}

	im.tasks <- incident
}

func (im *IncidentManager) logError(err error) {
	if err == nil {
		return
	}

	log.Println(err)
}

func (im *IncidentManager) saveIncident(inc *cm.Incident) (err error) {
	if inc.UserIPA == nil {
		err = im.dbo.SaveIncidentWithoutUserIPA(app.ModuleId_NM, inc.Type, inc.Email)
	} else {
		err = im.dbo.SaveIncident(app.ModuleId_NM, inc.Type, inc.Email, inc.UserIPA)
	}
	if err != nil {
		return err
	}

	return nil
}

func (im *IncidentManager) informGateway(inc *cm.Incident) (re *jrm1.RpcError) {
	blockTime := im.blockTimePerIncidentType[inc.Type.GetValue()]

	// Some incidents are only statistical.
	if blockTime == 0 {
		return nil
	}

	// Some incidents may have an empty IP address.
	// By the way, Go language does not even check anything and returns the
	// `<nil>` string if the IP address is empty. This is a very serious bug in
	// the language, but developers are too stupid to understand this.
	// https://github.com/golang/go/issues/39516
	if inc.UserIPA == nil {
		return nil
	}

	// Other incidents must be directed to the Gateway module.
	var params = gm.BlockIPAddressParams{
		UserIPA:      cm.IPAS(inc.UserIPA.String()),
		BlockTimeSec: blockTime,
	}

	var result = new(gm.BlockIPAddressResult)
	var err error
	re, err = im.gwmClient.MakeRequest(context.Background(), gc.FuncBlockIPAddress, params, result)
	if err != nil {
		im.logError(err)
		return jrm1.NewRpcErrorByUser(c.RpcErrorCode_RPCCall, c.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return re
	}
	if !result.OK {
		err = errors.New(fmt.Sprintf(c.MsgFModuleIsBroken, app.ServiceShortName_GWM))
		im.logError(err)
		return jrm1.NewRpcErrorByUser(c.RpcErrorCode_RPCCall, c.RpcErrorMsg_RPCCall, nil)
	}

	return nil
}
