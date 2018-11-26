package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/sacloud/packer-builder-sakuracloud/iaas"
)

type stepPowerOff struct {
	Debug bool
}

func (s *stepPowerOff) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	serverClient := state.Get("serverClient").(iaas.ServerClient)
	c := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)
	serverID := state.Get("server_id").(int64)

	stepStartMsg(ui, s.Debug, "PowerOff")

	server, err := serverClient.Read(serverID)
	if err != nil {
		err := fmt.Errorf("Error checking server state: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if server.Instance.IsDown() {
		// Server is already off, don't do anything
		stepEndMsg(ui, s.Debug, "PowerOff")
		return multistep.ActionContinue
	}

	// Pull the plug on the Droplet
	ui.Say("\tForcefully shutting down Droplet...")
	_, err = serverClient.Stop(serverID)
	if err != nil {
		err := fmt.Errorf("Error powering off server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	err = serverClient.SleepUntilDown(serverID, c.APIClientTimeout)
	if err != nil {
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
