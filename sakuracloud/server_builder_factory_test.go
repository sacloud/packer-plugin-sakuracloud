package sakuracloud

import (
	"testing"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/sacloud/libsacloud/builder"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/libsacloud/sacloud/ostype"
	"github.com/sacloud/packer-builder-sakuracloud/iaas"
	"github.com/stretchr/testify/assert"
)

type dummyBuilderAPIResults struct {
	getBySpecResult           *sacloud.ProductServer
	getBySpecError            error
	archiveFindByOSTypeResult *sacloud.Archive
	archiveFindByOSTypeError  error
}

func (t *dummyBuilderAPIResults) init() {
	t.getBySpecResult = nil
	t.getBySpecError = nil
	t.archiveFindByOSTypeResult = nil
	t.archiveFindByOSTypeError = nil
}

var dummyBuilderAPIResult = &dummyBuilderAPIResults{}
var dummyTestBuilder = &iaas.Builder{
	APIClient: &dummyBuilderAPIClient{
		getBySpecFunc: func(core int, memGB int, gen sacloud.PlanGenerations) (*sacloud.ProductServer, error) {
			return dummyBuilderAPIResult.getBySpecResult, dummyBuilderAPIResult.getBySpecError
		},
		archiveFindByOSTypeFunc: func(os ostype.ArchiveOSTypes) (*sacloud.Archive, error) {
			return dummyBuilderAPIResult.archiveFindByOSTypeResult, dummyBuilderAPIResult.archiveFindByOSTypeError
		},
	},
}

func TestCreateServerBuilder(t *testing.T) {

	dummyBuilderAPIResult.init()
	dummyBuilderAPIResult.archiveFindByOSTypeResult = &sacloud.Archive{Resource: sacloud.NewResource(dummyArchiveID)}

	expects := []struct {
		caseName       string
		config         map[string]interface{}
		builderExpects createServerBuilderExpects
	}{
		{
			caseName: "with iso-image",
			config: map[string]interface{}{
				"access_token":        "aaaa",
				"access_token_secret": "bbbb",
				"zone":                "is1a",
				"os_type":             "iso",
				"password":            dummyServerPassword,
				"iso_id":              dummyISOImageID,
			},
			builderExpects: createServerBuilderExpects{
				hasNetworkInterfaceProp: true,
				hasDiskProp:             true,
				hasDiskSourceProp:       false,
				hasHasDiskEditProp:      false,
				isoImageID:              dummyISOImageID,
				core:                    defaultCore,
				memory:                  defaultMemory,
				diskConnection:          sacloud.DiskConnectionVirtio,
				diskPlan:                sacloud.DiskPlanSSDID,
				diskSize:                defaultDiskSize,
			},
		},
		{
			caseName: "with source-archive-id",
			config: map[string]interface{}{
				"access_token":        "aaaa",
				"access_token_secret": "bbbb",
				"zone":                "is1a",
				"os_type":             "custom",
				"source_archive":      dummyArchiveID,
			},
			builderExpects: createServerBuilderExpects{
				hasNetworkInterfaceProp: true,
				hasDiskProp:             true,
				hasDiskSourceProp:       true,
				hasHasDiskEditProp:      true,
				sourceArchiveID:         dummyArchiveID,
				core:                    defaultCore,
				memory:                  defaultMemory,
				diskConnection:          sacloud.DiskConnectionVirtio,
				diskPlan:                sacloud.DiskPlanSSDID,
				diskSize:                defaultDiskSize,
				sshKeys:                 []string{dummySSHKeyBody},
				hostName:                defaultHostName,
			},
		},
		{
			caseName: "with source-disk-id",
			config: map[string]interface{}{
				"access_token":        "aaaa",
				"access_token_secret": "bbbb",
				"zone":                "is1a",
				"os_type":             "custom",
				"source_disk":         dummyDiskID,
			},
			builderExpects: createServerBuilderExpects{
				hasNetworkInterfaceProp: true,
				hasDiskProp:             true,
				hasDiskSourceProp:       true,
				hasHasDiskEditProp:      true,
				sourceDiskID:            dummyDiskID,
				core:                    defaultCore,
				memory:                  defaultMemory,
				diskConnection:          sacloud.DiskConnectionVirtio,
				diskPlan:                sacloud.DiskPlanSSDID,
				diskSize:                defaultDiskSize,
				sshKeys:                 []string{dummySSHKeyBody},
				hostName:                defaultHostName,
			},
		},
		{
			caseName: "with windows public-archive",
			config: map[string]interface{}{
				"access_token":        "aaaa",
				"access_token_secret": "bbbb",
				"zone":                "is1a",
				"os_type":             "windows2016",
			},
			builderExpects: createServerBuilderExpects{
				hasNetworkInterfaceProp: true,
				hasDiskProp:             true,
				hasDiskSourceProp:       true,
				hasHasDiskEditProp:      false,
				sourceArchiveID:         dummyArchiveID,
				core:                    defaultCore,
				memory:                  defaultMemory,
				diskConnection:          sacloud.DiskConnectionVirtio,
				diskPlan:                sacloud.DiskPlanSSDID,
				diskSize:                defaultDiskSize,
			},
		},
		{
			caseName: "with fixed unix archive",
			config: map[string]interface{}{
				"access_token":        "aaaa",
				"access_token_secret": "bbbb",
				"zone":                "is1a",
				"os_type":             "sophos-utm",
			},
			builderExpects: createServerBuilderExpects{
				hasNetworkInterfaceProp: true,
				hasDiskProp:             true,
				hasDiskSourceProp:       true,
				hasHasDiskEditProp:      false,
				sourceArchiveID:         dummyArchiveID,
				core:                    defaultCore,
				memory:                  defaultMemory,
				diskConnection:          sacloud.DiskConnectionVirtio,
				diskPlan:                sacloud.DiskPlanSSDID,
				diskSize:                defaultDiskSize,
			},
		},
		{
			caseName: "with fixed unix archive",
			config: map[string]interface{}{
				"access_token":        "aaaa",
				"access_token_secret": "bbbb",
				"zone":                "is1a",
				"os_type":             "centos",
				"password":            dummyServerPassword,
			},
			builderExpects: createServerBuilderExpects{
				hasNetworkInterfaceProp: true,
				hasDiskProp:             true,
				hasDiskSourceProp:       true,
				hasHasDiskEditProp:      true,
				sourceArchiveID:         dummyArchiveID,
				core:                    defaultCore,
				memory:                  defaultMemory,
				diskConnection:          sacloud.DiskConnectionVirtio,
				diskPlan:                sacloud.DiskPlanSSDID,
				diskSize:                defaultDiskSize,
				sshKeys:                 []string{dummySSHKeyBody},
				hostName:                defaultHostName,
				password:                dummyServerPassword,
			},
		},
		{
			caseName: "all parameters",
			config: map[string]interface{}{
				"access_token":        "aaaa",
				"access_token_secret": "bbbb",
				"zone":                "is1a",
				"os_type":             "centos",
				"core":                dummyServerCore,
				"memory_size":         dummyServerMemory,
				"disable_virtio_net":  true,
				"us_keyboard":         true,
				"disk_size":           dummyDiskSize,
				"disk_connection":     "ide",
				"disk_plan":           "hdd",
				"password":            dummyServerPassword,
			},
			builderExpects: createServerBuilderExpects{
				hasNetworkInterfaceProp: true,
				hasDiskProp:             true,
				hasDiskSourceProp:       true,
				hasHasDiskEditProp:      true,
				sourceArchiveID:         dummyArchiveID,
				core:                    dummyServerCore,
				memory:                  dummyServerMemory,
				disableVirtIONetPCI:     true,
				useUSKeyboard:           true,
				diskSize:                dummyDiskSize,
				diskConnection:          sacloud.DiskConnectionIDE,
				diskPlan:                sacloud.DiskPlanHDDID,
				sshKeys:                 []string{dummySSHKeyBody},
				hostName:                defaultHostName,
				password:                dummyServerPassword,
			},
		},
	}
	for _, expect := range expects {
		t.Run(expect.caseName, func(t *testing.T) {

			config := dummyConfigWithValues(expect.config)
			state := initCreateServerBuilderState(&config)
			factory := &defaultServerBuilderFactory{}
			actualBuilder := factory.createServerBuilder(state)

			be := expect.builderExpects

			// check builder's interface
			b, ok := actualBuilder.(builder.Builder)
			assert.True(t, ok)
			assert.Equal(t, be.hasNetworkInterfaceProp, b.HasNetworkInterfaceProperty(), "unexpected value: HasNetworkInterfaceProperty")
			assert.Equal(t, be.hasDiskProp, b.HasDiskProperty(), "unexpected value: HasDiskProperty")
			assert.Equal(t, be.hasDiskSourceProp, b.HasDiskSourceProperty(), "unexpected value: HasDiskSourceProperty")
			assert.Equal(t, be.hasHasDiskEditProp, b.HasDiskEditProperty(), "unexpected value: HasDiskEditProperty")

			// check sources
			dsv := b.(builder.DiskSourceProperty)
			assert.Equal(t, be.sourceDiskID, dsv.GetSourceDiskID())
			assert.Equal(t, be.sourceArchiveID, dsv.GetSourceArchiveID())

			// check common props
			cv := b.(builder.CommonProperty)
			assert.Equal(t, be.isoImageID, cv.GetISOImageID())
			assert.Equal(t, be.core, cv.GetCore())
			assert.Equal(t, be.memory, cv.GetMemory())

			nicConn := sacloud.InterfaceDriverVirtIO
			if be.disableVirtIONetPCI {
				nicConn = sacloud.InterfaceDriverE1000
			}
			assert.Equal(t, nicConn, cv.GetInterfaceDriver())
			assert.Equal(t, be.useUSKeyboard, hasUSKeyboardTag(cv.GetTags()))

			// check disk props
			dv := b.(builder.DiskProperty)
			assert.Equal(t, be.diskSize, dv.GetDiskSize())
			assert.Equal(t, be.diskConnection, dv.GetDiskConnection())
			assert.Equal(t, be.diskPlan, dv.GetDiskPlanID())

			// disk edit values
			dev := b.(builder.DiskEditProperty)
			assert.Equal(t, be.sshKeys, dev.GetSSHKeys())
			assert.Equal(t, be.hostName, dev.GetHostName())
			assert.Equal(t, be.password, dev.GetPassword())
		})
	}
}

func hasUSKeyboardTag(tags []string) bool {
	for _, v := range tags {
		if v == sacloud.TagKeyboardUS {
			return true
		}
	}
	return false
}

type createServerBuilderExpects struct {
	hasNetworkInterfaceProp bool
	hasDiskProp             bool
	hasDiskSourceProp       bool
	hasHasDiskEditProp      bool
	isoImageID              int64
	sourceDiskID            int64
	sourceArchiveID         int64
	core                    int
	memory                  int
	disableVirtIONetPCI     bool
	useUSKeyboard           bool
	diskSize                int
	diskConnection          sacloud.EDiskConnection
	diskPlan                sacloud.DiskPlanID
	sshKeys                 []string
	hostName                string
	password                string
}

func initCreateServerBuilderState(config *Config) multistep.StateBag {
	state := dummyMinimumStateBag(config)
	state.Put("publicKeys", []string{dummySSHKeyBody})

	state.Put("builder", dummyTestBuilder)

	return state
}
