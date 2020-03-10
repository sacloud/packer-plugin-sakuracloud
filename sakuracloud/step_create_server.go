package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/ostype"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	diskBuilders "github.com/sacloud/libsacloud/v2/utils/builder/disk"
	serverBuilders "github.com/sacloud/libsacloud/v2/utils/builder/server"
	"github.com/sacloud/libsacloud/v2/utils/power"
	"github.com/sacloud/packer-builder-sakuracloud/iaas"
	"github.com/sacloud/packer-builder-sakuracloud/sakuracloud/constants"
)

type stepCreateServer struct {
	Debug    bool
	serverID types.ID
	diskIDs  []types.ID
}

func (s *stepCreateServer) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	c := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)

	stepStartMsg(ui, s.Debug, "CreateServer")

	// create Server
	ui.Say("\tCreating server...")

	builder := s.createServerBuilder(state)
	created, err := builder.Build(ctx, c.Zone)

	if err != nil {
		err := fmt.Errorf("Error creating server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// We use this in cleanup
	s.serverID = created.ServerID
	s.diskIDs = created.DiskIDs

	// Store the server and disk id for later
	state.Put("server_id", s.serverID)
	state.Put("disk_id", s.diskIDs[0])

	stepEndMsg(ui, s.Debug, "CreateServer")
	return multistep.ActionContinue
}

func (s *stepCreateServer) Cleanup(state multistep.StateBag) {
	// If the serverID isn't there, we probably never created it
	if s.serverID.IsEmpty() && len(s.diskIDs) == 0 {
		return
	}

	c := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)
	caller := state.Get("iaasClient").(iaas.Client).Caller
	serverOp := sacloud.NewServerOp(caller)
	ctx := context.Background()

	// Destroy the server we just created
	ui.Say("\tDestroying server...")

	// force shutdown
	if err := power.ShutdownServer(ctx, serverOp, c.Zone, s.serverID, true); err != nil {
		ui.Error(fmt.Sprintf("Error destroying server. Please destroy it manually: %s", err))
		return
	}

	// delete server with disks
	var err error
	if len(s.diskIDs) == 0 {
		err = serverOp.Delete(ctx, c.Zone, s.serverID)
	} else {
		err = serverOp.DeleteWithDisks(ctx, c.Zone, s.serverID, &sacloud.ServerDeleteWithDisksRequest{IDs: s.diskIDs})
	}
	if err != nil {
		ui.Error(fmt.Sprintf(
			"Error destroying server. Please destroy it manually: %s", err))
	}
}

func (s *stepCreateServer) createServerBuilder(state multistep.StateBag) *serverBuilders.Builder {
	c := state.Get("config").(Config)
	caller := state.Get("iaasClient").(iaas.Client).Caller

	interfaceDriver := types.InterfaceDrivers.VirtIO
	if c.DisableVirtIONetPCI {
		interfaceDriver = types.InterfaceDrivers.E1000
	}

	builder := &serverBuilders.Builder{
		Name:            constants.ServerName,
		CPU:             c.Core,
		MemoryGB:        c.MemorySize,
		Commitment:      types.Commitments.Standard, // TODO
		Generation:      types.PlanGenerations.Default,
		InterfaceDriver: interfaceDriver,
		BootAfterCreate: true,
		CDROMID:         c.ISOImageID,
		NIC:             &serverBuilders.SharedNICSetting{}, // TODO 共有セグメントのみサポート
		DiskBuilders:    s.createDiskBuilder(state),
		Client:          serverBuilders.NewBuildersAPIClient(caller),
		ForceShutdown:   false,
	}
	return builder
}

func (s *stepCreateServer) createDiskBuilder(state multistep.StateBag) []diskBuilders.Builder {
	c := state.Get("config").(Config)
	caller := state.Get("iaasClient").(iaas.Client).Caller

	director := diskBuilders.Director{
		OSType:          ostype.StrToOSType(c.OSType),
		Name:            constants.ServerName,
		SizeGB:          c.DiskSize,
		DistantFrom:     nil,
		PlanID:          types.DiskPlanIDMap[c.DiskPlan],
		Connection:      types.DiskConnectionMap[c.DiskConnection],
		SourceDiskID:    c.SourceDisk,
		SourceArchiveID: c.SourceArchive,
		EditParameter: &diskBuilders.EditRequest{
			HostName: constants.ServerName,
			Password: c.Password,
		},
		Client: diskBuilders.NewBuildersAPIClient(caller),
	}

	if keys, ok := state.GetOk("publicKeys"); ok {
		director.EditParameter.SSHKeys = keys.([]string)
	}

	return []diskBuilders.Builder{director.Builder()}
}
