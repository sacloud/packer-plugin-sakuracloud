package sakuracloud

import (
	"fmt"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/builder"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/libsacloud/sacloud/ostype"
)

type serverBuilder interface {
	// 構築
	Build() (*builder.ServerBuildResult, error)
}

type stepCreateServer struct {
	Debug    bool
	serverID int64
	diskIDs  []int64
}

func (s *stepCreateServer) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)

	// create Server
	ui.Say("Creating server...")

	b := s.createServerBuilder(state)
	createResult, err := b.Build()

	if err != nil {
		err := fmt.Errorf("Error creating server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// We use this in cleanup
	s.serverID = createResult.Server.ID
	s.diskIDs = s.getDiskIDs(createResult)

	// Store the server and disk id for later
	state.Put("server_id", s.serverID)
	state.Put("disk_id", s.diskIDs[0])

	return multistep.ActionContinue
}

func (s *stepCreateServer) Cleanup(state multistep.StateBag) {
	// If the serverID isn't there, we probably never created it
	if s.serverID == 0 && len(s.diskIDs) == 0 {
		return
	}

	client := state.Get("client").(*api.Client)
	client.TraceMode = true
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)

	// Destroy the server we just created
	ui.Say("Destroying server...")

	// force shutdown
	_, err := client.Server.Stop(s.serverID)
	if err != nil {
		ui.Error(fmt.Sprintf("Error destroying server. Please destroy it manually: %s", err))
		return
	}
	// wait for down
	err = client.Server.SleepUntilDown(s.serverID, c.StateTimeout)
	if err != nil {
		ui.Error(fmt.Sprintf("Error destroying server. Please destroy it manually: %s", err))
		return
	}

	// delete server with disks
	if len(s.diskIDs) == 0 {
		_, err = client.Server.Delete(s.serverID)
	} else {
		_, err = client.Server.DeleteWithDisk(s.serverID, s.diskIDs)
	}
	if err != nil {
		ui.Error(fmt.Sprintf(
			"Error destroying server. Please destroy it manually: %s", err))
	}
}

func (s *stepCreateServer) createServerBuilder(state multistep.StateBag) serverBuilder {

	client := state.Get("client").(*api.Client)
	c := state.Get("config").(Config)

	serverName := "packer_builder_sakuracloud"

	switch c.OSType {
	case "centos", "ubuntu", "debian", "coreos", "kusanagi", "vyos":
		b := builder.ServerPublicArchiveUnix(client, s.getOSTypeFromString(c.OSType), serverName, c.Password).
			SetCore(c.Core).
			SetMemory(c.MemorySize).
			SetUseVirtIONetPCI(!c.DisableVirtIONetPCI).
			AddPublicNWConnectedNIC().
			SetDiskSize(c.DiskSize).
			SetDiskConnection(s.getDiskConnection(c)).
			SetDiskPlanID(s.getDiskPlanID(c)).
			AddSSHKey(state.Get("ssh_public_key").(string)).
			SetHostName(serverName)
		if c.ISOImageID > 0 {
			b.SetISOImageID(c.ISOImageID)
		}
		if c.UseUSKeyboard {
			b.AppendTag(sacloud.TagKeyboardUS)
		}
		return b
	case "custom":
		var b *builder.CommonServerBuilder
		if c.SourceArchive > 0 {
			b = builder.ServerFromArchive(client, serverName, c.SourceArchive)
		} else {
			b = builder.ServerFromDisk(client, serverName, c.SourceDisk)
		}
		if c.ISOImageID > 0 {
			b.SetISOImageID(c.ISOImageID)
		}
		if c.UseUSKeyboard {
			b.AppendTag(sacloud.TagKeyboardUS)
		}

		return b.SetCore(c.Core).
			SetMemory(c.MemorySize).
			SetUseVirtIONetPCI(!c.DisableVirtIONetPCI).
			AddPublicNWConnectedNIC().
			SetDiskSize(c.DiskSize).
			SetDiskConnection(s.getDiskConnection(c)).
			SetDiskPlanID(s.getDiskPlanID(c)).
			AddSSHKey(state.Get("ssh_public_key").(string)).
			SetHostName(serverName)

	case "windows":
		b := builder.ServerPublicArchiveWindows(client, serverName, c.SourceArchive).
			SetCore(c.Core).
			SetMemory(c.MemorySize).
			SetUseVirtIONetPCI(!c.DisableVirtIONetPCI).
			AddPublicNWConnectedNIC().
			SetDiskSize(c.DiskSize).
			SetDiskConnection(s.getDiskConnection(c)).
			SetDiskPlanID(s.getDiskPlanID(c))
		if c.ISOImageID > 0 {
			b.SetISOImageID(c.ISOImageID)
		}
		if c.UseUSKeyboard {
			b.AppendTag(sacloud.TagKeyboardUS)
		}
		return b
	case "iso":
		var b *builder.BlankDiskServerBuilder
		b = builder.ServerBlankDisk(client, serverName)
		if c.ISOImageID > 0 {
			b.SetISOImageID(c.ISOImageID)
		}
		if c.UseUSKeyboard {
			b.AppendTag(sacloud.TagKeyboardUS)
		}
		b.SetCore(c.Core).
			SetMemory(c.MemorySize).
			SetUseVirtIONetPCI(!c.DisableVirtIONetPCI).
			AddPublicNWConnectedNIC().
			SetDiskSize(c.DiskSize).
			SetDiskConnection(s.getDiskConnection(c)).
			SetDiskPlanID(s.getDiskPlanID(c))
		return b
	}

	return nil
}

func (s *stepCreateServer) getOSTypeFromString(os string) ostype.ArchiveOSTypes {
	switch os {
	case "centos":
		return ostype.CentOS
	case "ubuntu":
		return ostype.Ubuntu
	case "debian":
		return ostype.Debian
	case "vyos":
		return ostype.VyOS
	case "coreos":
		return ostype.CoreOS
	case "kusanagi":
		return ostype.Kusanagi
	case "custom":
		return ostype.Custom
	case "windows":
		return ostype.Custom
	}
	panic(fmt.Errorf("invalid ostype [%s]", os))
}

func (s *stepCreateServer) getDiskConnection(config Config) sacloud.EDiskConnection {
	switch config.DiskConnection {
	case "ide":
		return sacloud.DiskConnectionIDE
	case "virtio":
		return sacloud.DiskConnectionVirtio
	}

	panic(fmt.Errorf("invalid config: disk_connection[%s]", config.DiskConnection))
}

func (s *stepCreateServer) getDiskPlanID(config Config) sacloud.DiskPlanID {
	switch config.DiskPlan {
	case "ssd":
		return sacloud.DiskPlanSSDID
	case "hdd":
		return sacloud.DiskPlanHDDID
	}

	panic(fmt.Errorf("invalid config: disk_plan[%s]", config.DiskPlan))
}

func (s *stepCreateServer) getDiskIDs(buildResult *builder.ServerBuildResult) []int64 {
	res := []int64{}
	for _, disk := range buildResult.Disks {
		res = append(res, disk.Disk.ID)
	}
	return res
}
