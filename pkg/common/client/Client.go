package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"sync/atomic"

	jc "github.com/ybbus/jsonrpc/v3"
)

const (
	ErrFUnknownUrlScheme = "unknown URL scheme: %s"
)

const (
	UrlSchemeHttp  = "http"
	UrlSchemeHttps = "https"
)

type Client struct {
	jc            jc.RPCClient
	lastRequestId atomic.Int64
}

func NewClient(dsn string, enableSelfSignedCertificate bool) (client *Client, err error) {
	var dsnUrl *url.URL
	dsnUrl, err = url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	if dsnUrl.Scheme == UrlSchemeHttp {
		return newStandardClient(dsn), nil
	}

	if dsnUrl.Scheme == UrlSchemeHttps {
		if !enableSelfSignedCertificate {
			return newStandardClient(dsn), nil
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

		return newClientWithOpts(dsn, opts), nil

	}

	return nil, fmt.Errorf(ErrFUnknownUrlScheme, dsnUrl.Scheme)
}

func newStandardClient(dsn string) (client *Client) {
	return &Client{
		jc:            jc.NewClient(dsn),
		lastRequestId: atomic.Int64{},
	}
}

func newClientWithOpts(dsn string, opts *jc.RPCClientOpts) (client *Client) {
	return &Client{
		jc:            jc.NewClientWithOpts(dsn, opts),
		lastRequestId: atomic.Int64{},
	}
}

func (c *Client) MakeRequest(ctx context.Context, out any, method string, params ...any) error {
	curRequestId := c.lastRequestId.Add(1)

	req := jc.NewRequestWithID(int(curRequestId), method, params...)

	resp, err := c.jc.CallRaw(ctx, req)
	if err != nil {
		return err
	}

	if resp.Error != nil {
		return resp.Error
	}

	return resp.GetObject(out)
}
