package sakuracloud

import (
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

func testUI() packer.Ui {
	return new(packer.NoopUi)
}

func testStateBag() multistep.StateBag {
	return new(multistep.BasicStateBag)
}

func testConfig() Config {
	values := map[string]interface{}{
		"access_token":        "aaaa",
		"access_token_secret": "bbbb",
		"zone":                "is1a",
		"os_type":             "centos",
	}

	conf, _, _ := NewConfig(values)
	return *conf
}
