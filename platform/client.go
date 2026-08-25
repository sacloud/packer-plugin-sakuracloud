package platform

import (
	"fmt"
	"os"

	"github.com/sacloud/packer-plugin-sakuracloud/version"
	"github.com/sacloud/sacloud-sdk-go/api/iaas"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
)

// Client represents SakuraCloud API Client
type Client struct {
	Caller  iaas.APICaller
	Archive Archive
	Zone    string
}

// NewClient returns new SakuraCloud API Client
func NewClient(token, secret, zone string) (*Client, error) {
	var sa saclient.Client

	// saclient.Client は環境変数、コマンドライン引数、プロファイル(~/.usacloud) などから
	// 標準的に設定を読み込む。ここではまず os.Environ() を読み込ませ、Packer 設定値が
	// 明示的に指定されていれば WithBasicAuth で上書きする。
	if err := sa.SetEnviron(os.Environ()); err != nil {
		return nil, err
	}
	if err := sa.SetWith(saclient.WithUserAgent(fmt.Sprintf("packer-plugin-sakuracloud:v%s", version.Version))); err != nil {
		return nil, err
	}
	if token != "" && secret != "" {
		if err := sa.SetWith(saclient.WithBasicAuth(token, secret)); err != nil {
			return nil, err
		}
	}
	// 旧実装と同等のレートリミットを維持（ユーザーが環境変数で上書き可能）
	if os.Getenv("SAKURA_RATE_LIMIT") == "" && os.Getenv("SAKURACLOUD_API_REQUEST_RATE_LIMIT") == "" {
		if err := sa.SetWith(saclient.WithAPIRequestRateLimit(3)); err != nil {
			return nil, err
		}
	}
	if err := sa.Populate(); err != nil {
		return nil, err
	}

	caller := iaas.NewClientFromSaclient(&sa)

	return &Client{
		Caller:  caller,
		Archive: newArchiveClient(caller, zone),
		Zone:    zone,
	}, nil
}
