package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/sacloud/libsacloud/builder"
	"github.com/sacloud/packer-builder-sakuracloud/iaas"
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

func (s *stepCreateServer) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)

	stepStartMsg(ui, s.Debug, "CreateServer")

	// create Server
	ui.Say("\tCreating server...")

	factory := state.Get("builderFactory").(serverBuilderFactory)
	b := factory.createServerBuilder(state)
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

	stepEndMsg(ui, s.Debug, "CreateServer")
	return multistep.ActionContinue
}

func (s *stepCreateServer) Cleanup(state multistep.StateBag) {
	// If the serverID isn't there, we probably never created it
	if s.serverID == 0 && len(s.diskIDs) == 0 {
		return
	}

	serverClient := state.Get("serverClient").(iaas.ServerClient)
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)

	// Destroy the server we just created
	ui.Say("\tDestroying server...")

	// force shutdown
	_, err := serverClient.Stop(s.serverID)
	if err != nil {
		ui.Error(fmt.Sprintf("Error destroying server. Please destroy it manually: %s", err))
		return
	}
	// wait for down
	err = serverClient.SleepUntilDown(s.serverID, c.APIClientTimeout)
	if err != nil {
		ui.Error(fmt.Sprintf("Error destroying server. Please destroy it manually: %s", err))
		return
	}

	// delete server with disks
	if len(s.diskIDs) == 0 {
		_, err = serverClient.Delete(s.serverID)
	} else {
		_, err = serverClient.DeleteWithDisk(s.serverID, s.diskIDs)
	}
	if err != nil {
		ui.Error(fmt.Sprintf(
			"Error destroying server. Please destroy it manually: %s", err))
	}
}

func (s *stepCreateServer) getDiskIDs(buildResult *builder.ServerBuildResult) []int64 {
	var res []int64
	for _, disk := range buildResult.Disks {
		res = append(res, disk.Disk.ID)
	}
	return res
}
