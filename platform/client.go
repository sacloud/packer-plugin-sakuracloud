package platform

import (
	"crypto/tls"
	"fmt"

	client "github.com/sacloud/api-client-go"
	"github.com/sacloud/ftps"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/api"
	"github.com/sacloud/packer-plugin-sakuracloud/version"
)

// Client represents SakuraCloud API Client
type Client struct {
	Caller  iaas.APICaller
	Archive Archive
	FTPS    FTPSClient
	Zone    string
}

// NewClient returns new SakuraCloud API Client
func NewClient(token, secret, zone string) (*Client, error) {
	options, err := api.DefaultOption()
	if err != nil {
		return nil, err
	}
	options = api.MergeOptions(options, &api.CallerOptions{
		Options: &client.Options{
			AccessToken:          token,
			AccessTokenSecret:    secret,
			UserAgent:            fmt.Sprintf("packer-plugin-sakuracloud:v%s", version.Version),
			HttpRequestRateLimit: 3,
		},
	})
	caller := api.NewCallerWithOptions(options)

	// FTPS Client
	ftpsClient := &ftps.FTPS{
		TLSConfig: tls.Config{
			InsecureSkipVerify: true, //nolint:gosec
			MaxVersion:         tls.VersionTLS12,
		},
	}

	return &Client{
		Caller:  caller,
		Archive: newArchiveClient(caller, zone),
		FTPS:    FTPSClient(ftpsClient),
		Zone:    zone,
	}, nil
}
