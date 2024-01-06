package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	cs "github.com/vault-thirteen/SimpleBB/pkg/common/settings"
	"github.com/vault-thirteen/auxie/number"
)

const (
	ErrServerClientSettingsAreNotSet = "server client settings are not set"
	ErrShortNameIsNotSet             = "short name is not set"
)

const (
	UrlSchemeHttp  = "http"
	UrlSchemeHttps = "https"
)

const (
	FuncPing               = "Ping"
	FuncShowDiagnosticData = "ShowDiagnosticData"
)

type Client struct {
	shortName string
	jc        *jrm1.Client
}

// NewClient is a constructor of an RPC client.
// Port in DSN must be explicitly set.
func NewClient(shortName string, dsn string, enableSelfSignedCertificate bool) (client *Client, err error) {
	if len(shortName) == 0 {
		return nil, errors.New(ErrShortNameIsNotSet)
	}

	var dsnUrl *url.URL
	dsnUrl, err = url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	var customHttpClient *http.Client
	if (dsnUrl.Scheme == UrlSchemeHttps) && enableSelfSignedCertificate {
		customHttpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}

	return newCustomClient(shortName, dsnUrl, customHttpClient)
}

func NewClientWithSCS(scs *cs.ServiceClientSettings, shortName string) (serviceClient *Client, err error) {
	if scs == nil {
		return nil, errors.New(ErrServerClientSettingsAreNotSet)
	}

	dsn := fmt.Sprintf("%s://%s:%d%s", scs.Schema, scs.Host, scs.Port, scs.Path)

	serviceClient, err = NewClient(shortName, dsn, scs.EnableSelfSignedCertificate)
	if err != nil {
		return nil, err
	}

	return serviceClient, nil
}

func newCustomClient(shortName string, dsnUrl *url.URL, customHttpClient *http.Client) (client *Client, err error) {
	var port uint16
	port, err = number.ParseUint16(dsnUrl.Port())
	if err != nil {
		return nil, err
	}

	path := dsnUrl.RequestURI()

	var clientSettings *jrm1.ClientSettings
	clientSettings, err = jrm1.NewClientSettings(dsnUrl.Scheme, dsnUrl.Hostname(), port, path, customHttpClient, nil, true)
	if err != nil {
		return nil, err
	}

	var rpcClient *jrm1.Client
	rpcClient, err = jrm1.NewClient(clientSettings)
	if err != nil {
		return nil, err
	}

	client = &Client{
		shortName: shortName,
		jc:        rpcClient,
	}

	return client, nil
}

func (cli *Client) MakeRequest(ctx context.Context, method string, params any, result any) (re *jrm1.RpcError, err error) {
	return cli.jc.Call(ctx, method, params, result)
}

func (cli *Client) Ping(verbose bool) (err error) {
	// While several services are inter-dependent, we make several attempts to
	// ping the module.

	if verbose {
		fmt.Print(fmt.Sprintf(c.MsgFPingingModule, cli.shortName))
	}

	var params = cmr.PingParams{}

	var result = new(cmr.PingResult)
	var re *jrm1.RpcError
	for i := 1; i <= c.ServicePingAttemptsCount; i++ {
		re, err = cli.MakeRequest(context.Background(), FuncPing, params, result)
		if (err == nil) && (re == nil) {
			break
		}

		if verbose {
			fmt.Print(c.MsgPingAttempt)
		}

		if i < c.ServicePingAttemptsCount {
			time.Sleep(time.Second * time.Duration(c.ServiceNextPingAttemptDelaySec))
		}
	}

	if err != nil {
		return err
	}
	if re != nil {
		return re.AsError()
	}
	if !result.OK {
		return errors.New(fmt.Sprintf(c.MsgFModuleIsBroken, cli.shortName))
	}

	if verbose {
		fmt.Println(c.MsgOK)
	}

	return nil
}
