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
	"github.com/vault-thirteen/SimpleBB/pkg/ACM/dbo"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	s "github.com/vault-thirteen/SimpleBB/pkg/ACM/settings"
	gc "github.com/vault-thirteen/SimpleBB/pkg/GWM/client"
	gm "github.com/vault-thirteen/SimpleBB/pkg/GWM/models"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	"github.com/vault-thirteen/SimpleBB/pkg/common/app"
	"github.com/vault-thirteen/SimpleBB/pkg/common/avm"
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/client"
)

const (
	TaskChannelSize = 4
)

type IncidentManager struct {
	ssp                    *avm.SSP
	wg                     *sync.WaitGroup
	tasks                  chan *am.Incident
	isTableOfIncidentsUsed bool
	dbo                    *dbo.DatabaseObject
	gwmClient              *cc.Client

	// Block time in seconds for each incident type.
	blockTimePerIncidentType [am.IncidentTypesCount + 1]uint
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
		tasks:                    make(chan *am.Incident, TaskChannelSize),
		isTableOfIncidentsUsed:   isTableOfIncidentsUsed,
		dbo:                      dbo,
		gwmClient:                gwmClient,
		blockTimePerIncidentType: initBlockTimePerIncidentType(blockTimePerIncident),
	}

	return im
}

func initBlockTimePerIncidentType(blockTimePerIncident *s.BlockTimePerIncident) (blockTimePerIncidentType [am.IncidentTypesCount + 1]uint) {
	// The "zero"-indexed element is empty because it is not used.
	blockTimePerIncidentType[am.IncidentType_IllegalAccessAttempt] = blockTimePerIncident.IllegalAccessAttempt
	blockTimePerIncidentType[am.IncidentType_FakeToken] = blockTimePerIncident.FakeToken
	blockTimePerIncidentType[am.IncidentType_VerificationCodeMismatch] = blockTimePerIncident.VerificationCodeMismatch
	blockTimePerIncidentType[am.IncidentType_DoubleLogInAttempt] = blockTimePerIncident.DoubleLogInAttempt
	blockTimePerIncidentType[am.IncidentType_PreSessionHacking] = blockTimePerIncident.PreSessionHacking
	blockTimePerIncidentType[am.IncidentType_CaptchaAnswerMismatch] = blockTimePerIncident.CaptchaAnswerMismatch
	blockTimePerIncidentType[am.IncidentType_PasswordMismatch] = blockTimePerIncident.PasswordMismatch
	blockTimePerIncidentType[am.IncidentType_PasswordChangeHacking] = blockTimePerIncident.PasswordChangeHacking
	blockTimePerIncidentType[am.IncidentType_EmailChangeHacking] = blockTimePerIncident.EmailChangeHacking

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
			err = am.CheckIncident(inc)
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

func (im *IncidentManager) ReportIncident(itype am.IncidentType, email string, userIPA net.IP) {
	incident := &am.Incident{
		Time:    time.Now(),
		Type:    itype,
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

func (im *IncidentManager) saveIncident(inc *am.Incident) (err error) {
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

func (im *IncidentManager) informGateway(inc *am.Incident) (re *jrm1.RpcError) {
	blockTime := im.blockTimePerIncidentType[inc.Type]

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
		UserIPA:      inc.UserIPA.String(),
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
