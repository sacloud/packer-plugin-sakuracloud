package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/power"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/packer-plugin-sakuracloud/platform"
)

type stepPowerOff struct {
	Debug bool
}

func (s *stepPowerOff) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	c := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)

	caller := state.Get("iaasClient").(*platform.Client).Caller
	serverOp := iaas.NewServerOp(caller)
	serverID := state.Get("server_id").(types.ID)

	stepStartMsg(ui, s.Debug, "PowerOff")

	server, err := serverOp.Read(ctx, c.Zone, serverID)
	if err != nil {
		err := fmt.Errorf("Error checking server state: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if server.InstanceStatus.IsDown() {
		// Server is already off, don't do anything
		stepEndMsg(ui, s.Debug, "PowerOff")
		return multistep.ActionContinue
	}

	ui.Say("\tShutting down the server...")

	if err := power.ShutdownServer(ctx, serverOp, c.Zone, serverID, c.ForceShutdown); err != nil {
		err := fmt.Errorf("Error powering off server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	stepEndMsg(ui, s.Debug, "PowerOff")
	return multistep.ActionContinue
}

func (s *stepPowerOff) Cleanup(state multistep.StateBag) {
	// no cleanup
}
