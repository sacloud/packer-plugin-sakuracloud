module github.com/sacloud/packer-plugin-sakuracloud

require (
	github.com/hashicorp/hcl/v2 v2.12.0
	github.com/hashicorp/packer-plugin-sdk v0.2.13
	github.com/mitchellh/go-vnc v0.0.0-20150629162542-723ed9867aed
	github.com/mitchellh/mapstructure v1.4.1
	github.com/sacloud/ftps v1.1.0
	github.com/sacloud/libsacloud/v2 v2.15.1
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/zclconf/go-cty v1.10.0
	golang.org/x/crypto v0.0.0-20220427172511-eb4f295cb31f
)

replace github.com/zclconf/go-cty => github.com/azr/go-cty v1.1.1-0.20200203143058-28fcda2fe0cc

go 1.16
