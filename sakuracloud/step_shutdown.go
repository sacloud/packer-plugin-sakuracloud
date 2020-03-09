package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/power"
)

type stepShutdown struct {
	Debug bool
}

func (s *stepShutdown) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	c := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)

	caller := state.Get("sacloudAPICaller").(sacloud.APICaller)
	serverOp := sacloud.NewServerOp(caller)
	serverID := state.Get("server_id").(types.ID)

	stepStartMsg(ui, s.Debug, "Shutdown Server")

	ui.Say("\tGracefully shutting down server...")

	if err := power.ShutdownServer(ctx, serverOp, c.Zone, serverID, c.ForceShutdown); err != nil {
		err := fmt.Errorf("Error shutting down server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	stepEndMsg(ui, s.Debug, "Shutdown Server")
	return multistep.ActionContinue
}

func (s *stepShutdown) Cleanup(state multistep.StateBag) {
	// no cleanup
}
