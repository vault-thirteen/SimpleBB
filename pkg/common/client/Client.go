package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cs "github.com/vault-thirteen/SimpleBB/pkg/common/settings"
	jc "github.com/ybbus/jsonrpc/v3"
)

const (
	ErrFUnknownUrlScheme             = "unknown URL scheme: %s"
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
	shortName     string
	jc            jc.RPCClient
	lastRequestId atomic.Int64
}

func NewClient(dsn string, enableSelfSignedCertificate bool, shortName string) (client *Client, err error) {
	if len(shortName) == 0 {
		return nil, errors.New(ErrShortNameIsNotSet)
	}

	var dsnUrl *url.URL
	dsnUrl, err = url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	if dsnUrl.Scheme == UrlSchemeHttp {
		return newStandardClient(dsn, shortName), nil
	}

	if dsnUrl.Scheme == UrlSchemeHttps {
		if !enableSelfSignedCertificate {
			return newStandardClient(dsn, shortName), nil
		}

		httpTransport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}

		httpClient := &http.Client{
			Transport: httpTransport,
		}

		opts := &jc.RPCClientOpts{
			HTTPClient: httpClient,
		}

		return newClientWithOpts(dsn, shortName, opts), nil

	}

	return nil, fmt.Errorf(ErrFUnknownUrlScheme, dsnUrl.Scheme)
}

func NewClientWithSCS(scs *cs.ServiceClientSettings, shortName string) (serviceClient *Client, err error) {
	if scs == nil {
		return nil, errors.New(ErrServerClientSettingsAreNotSet)
	}

	dsn := fmt.Sprintf("%s://%s:%d%s", scs.Schema, scs.Host, scs.Port, scs.Path)

	serviceClient, err = NewClient(dsn, scs.EnableSelfSignedCertificate, shortName)
	if err != nil {
		return nil, err
	}

	return serviceClient, nil
}

func newStandardClient(dsn string, shortName string) (client *Client) {
	return &Client{
		shortName:     shortName,
		jc:            jc.NewClient(dsn),
		lastRequestId: atomic.Int64{},
	}
}

func newClientWithOpts(dsn string, shortName string, opts *jc.RPCClientOpts) (client *Client) {
	return &Client{
		shortName:     shortName,
		jc:            jc.NewClientWithOpts(dsn, opts),
		lastRequestId: atomic.Int64{},
	}
}

func (cli *Client) MakeRequest(ctx context.Context, out any, method string, params ...any) error {
	curRequestId := cli.lastRequestId.Add(1)

	req := jc.NewRequestWithID(int(curRequestId), method, params...)

	resp, err := cli.jc.CallRaw(ctx, req)
	if err != nil {
		return err
	}

	if resp.Error != nil {
		return resp.Error
	}

	return resp.GetObject(out)
}

func (cli *Client) Ping(verbose bool) (err error) {
	// While several services are inter-dependent, we make several attempts to
	// ping the module.

	if verbose {
		fmt.Print(fmt.Sprintf(c.MsgFPingingModule, cli.shortName))
	}

	var params = cm.PingParams{}

	var result = new(cm.PingResult)
	for i := 1; i <= c.ServicePingAttemptsCount; i++ {
		err = cli.MakeRequest(context.Background(), result, FuncPing, params)
		if err == nil {
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

	if !result.OK {
		return errors.New(fmt.Sprintf(c.MsgFModuleIsBroken, cli.shortName))
	}

	if verbose {
		fmt.Println(c.MsgOK)
	}

	return nil
}
