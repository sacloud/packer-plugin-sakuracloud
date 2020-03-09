package iaas

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/trace"
	"github.com/sacloud/packer-builder-sakuracloud/version"
)

// Client represents SakuraCloud API Client
type Client struct {
	Caller  sacloud.APICaller
	Archive Archive
	Zone    string
}

// NewClient returns new SakuraCloud API Client
func NewClient(token, secret, zone string) *Client {
	httpClient := http.DefaultClient
	httpClient.Transport = &sacloud.RateLimitRoundTripper{RateLimitPerSec: 3, Transport: httpClient.Transport}

	if traceMode := os.Getenv("SAKURACLOUD_TRACE"); traceMode != "" {
		trace.AddClientFactoryHooks()
	}

	caller := &sacloud.Client{
		AccessToken:       token,
		AccessTokenSecret: secret,
		UserAgent:         fmt.Sprintf("packer_for_sakuracloud:v%s", version.Version),
		HTTPClient:        httpClient,
	}
	return &Client{
		Caller:  caller,
		Archive: newArchiveClient(caller, zone),
		Zone:    zone,
	}
}
