module github.com/sacloud/packer-builder-sakuracloud

require (
	github.com/hashicorp/hcl/v2 v2.3.0
	github.com/hashicorp/packer v1.5.4
	github.com/mitchellh/go-vnc v0.0.0-20150629162542-723ed9867aed
	github.com/mitchellh/mapstructure v1.1.2
	github.com/sacloud/ftps v0.0.0-20171205062625-42fc0f9886fe
	github.com/sacloud/libsacloud/v2 v2.3.0
	github.com/stretchr/testify v1.4.0
	github.com/zclconf/go-cty v1.2.1
	golang.org/x/crypto v0.0.0-20200117160349-530e935923ad
)

replace github.com/zclconf/go-cty => github.com/azr/go-cty v1.1.1-0.20200203143058-28fcda2fe0cc

go 1.13
