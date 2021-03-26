package iaas

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sacloud/ftps"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/trace"
	"github.com/sacloud/packer-plugin-sakuracloud/version"
)

// Client represents SakuraCloud API Client
type Client struct {
	Caller  sacloud.APICaller
	Archive Archive
	FTPS    FTPSClient
	Zone    string
}

// NewClient returns new SakuraCloud API Client
func NewClient(token, secret, zone string) *Client {
	// HTTP Client
	httpClient := http.DefaultClient
	httpClient.Transport = &sacloud.RateLimitRoundTripper{RateLimitPerSec: 3, Transport: httpClient.Transport}

	// Sacloud API Client
	if traceMode := os.Getenv("SAKURACLOUD_TRACE"); traceMode != "" {
		trace.AddClientFactoryHooks()
	}

	caller := &sacloud.Client{
		AccessToken:       token,
		AccessTokenSecret: secret,
		UserAgent:         fmt.Sprintf("packer_for_sakuracloud:v%s", version.Version),
		HTTPClient:        httpClient,
	}

	// FTPS Client
	ftpsClient := &ftps.FTPS{}
	ftpsClient.TLSConfig.InsecureSkipVerify = true

	return &Client{
		Caller:  caller,
		Archive: newArchiveClient(caller, zone),
		FTPS:    FTPSClient(ftpsClient),
		Zone:    zone,
	}
}
