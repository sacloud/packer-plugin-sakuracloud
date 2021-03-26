module github.com/sacloud/packer-plugin-sakuracloud

require (
	github.com/hashicorp/hcl/v2 v2.8.0
	github.com/hashicorp/packer-plugin-sdk v0.1.0
	github.com/mitchellh/go-vnc v0.0.0-20150629162542-723ed9867aed
	github.com/mitchellh/mapstructure v1.4.0
	github.com/sacloud/ftps v1.1.0
	github.com/sacloud/libsacloud/v2 v2.15.1
	github.com/stretchr/testify v1.6.1
	github.com/zclconf/go-cty v1.7.0
	golang.org/x/crypto v0.0.0-20201208171446-5f87f3452ae9
)

replace github.com/zclconf/go-cty => github.com/azr/go-cty v1.1.1-0.20200203143058-28fcda2fe0cc

go 1.16
