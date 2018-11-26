package sakuracloud

import (
	"fmt"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/sacloud/libsacloud/builder"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/libsacloud/sacloud/ostype"
	"github.com/sacloud/packer-builder-sakuracloud/iaas"
)

type serverBuilderFactory interface {
	createServerBuilder(multistep.StateBag) serverBuilder
}

type defaultServerBuilderFactory struct{}

func (s *defaultServerBuilderFactory) createServerBuilder(state multistep.StateBag) (b serverBuilder) {

	builderClient := state.Get("builder").(iaas.ServerBuilder)
	c := state.Get("config").(Config)

	serverName := defaultHostName
	os := s.getOSTypeFromString(c.OSType)

	switch {
	case c.OSType == "iso":
		b = builderClient.FromBlankDisk(serverName)
	case os == ostype.Custom:
		if c.SourceArchive > 0 {
			b = builderClient.FromArchive(serverName, c.SourceArchive)
		} else {
			b = builderClient.FromDisk(serverName, c.SourceDisk)
		}
	case os.IsWindows():
		b = builderClient.FromPublicArchiveWindows(os, serverName)
	case os == ostype.Netwiser, os == ostype.SophosUTM, os == ostype.OPNsense:
		b = builderClient.FromPublicArchiveFixedUnix(os, serverName)
	default:
		b = builderClient.FromPublicArchiveUnix(os, serverName, c.Password)
	}

	bu := b.(builder.Builder)
	if bu.HasCommonProperty() { // always true in this context
		b := b.(builder.CommonProperty)
		b.SetCore(c.Core)
		b.SetMemory(c.MemorySize)
		if c.DisableVirtIONetPCI {
			b.SetInterfaceDriver(sacloud.InterfaceDriverE1000)
		}
		if c.UseUSKeyboard {
			b.SetTags(append(b.GetTags(), sacloud.TagKeyboardUS))
		}
		if c.ISOImageID > 0 {
			b.SetISOImageID(c.ISOImageID)
		}
	}
	if bu.HasNetworkInterfaceProperty() { // always true in this context
		b := b.(builder.NetworkInterfaceProperty)
		b.AddPublicNWConnectedNIC()
	}
	if bu.HasDiskProperty() {
		b := b.(builder.DiskProperty)
		b.SetDiskSize(c.DiskSize)
		b.SetDiskConnection(s.getDiskConnection(c))
		b.SetDiskPlanID(s.getDiskPlanID(c))
	}
	if bu.HasDiskEditProperty() {
		b := b.(builder.DiskEditProperty)
		if keys, ok := state.GetOk("publicKeys"); ok {
			for _, key := range keys.([]string) {
				b.AddSSHKey(key)
			}
		}
		b.SetHostName(serverName)
	}

	return b
}

func (s *defaultServerBuilderFactory) getOSTypeFromString(os string) ostype.ArchiveOSTypes {
	return ostype.StrToOSType(os)
}

func (s *defaultServerBuilderFactory) getDiskConnection(config Config) sacloud.EDiskConnection {
	switch config.DiskConnection {
	case "ide":
		return sacloud.DiskConnectionIDE
	case "virtio":
		return sacloud.DiskConnectionVirtio
	}

	panic(fmt.Errorf("invalid config: disk_connection[%s]", config.DiskConnection))
}

func (s *defaultServerBuilderFactory) getDiskPlanID(config Config) sacloud.DiskPlanID {
	switch config.DiskPlan {
	case "ssd":
		return sacloud.DiskPlanSSDID
	case "hdd":
		return sacloud.DiskPlanHDDID
	}

	panic(fmt.Errorf("invalid config: disk_plan[%s]", config.DiskPlan))
}
