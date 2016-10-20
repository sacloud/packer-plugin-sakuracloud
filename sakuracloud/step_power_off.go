package sakuracloud

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"github.com/sacloud/libsacloud/api"
)

type stepPowerOff struct {
	Debug bool
}

func (s *stepPowerOff) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*api.Client)
	c := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)
	serverID := state.Get("server_id").(int64)

	server, err := client.Server.Read(serverID)
	if err != nil {
		err := fmt.Errorf("Error checking server state: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if server.Instance.IsDown() {
		// Server is already off, don't do anything
		return multistep.ActionContinue
	}

	// Pull the plug on the Droplet
	ui.Say("Forcefully shutting down Droplet...")
	_, err = client.Server.Stop(serverID)
	if err != nil {
		err := fmt.Errorf("Error powering off server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	err = client.Server.SleepUntilDown(serverID, c.StateTimeout)
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *stepPowerOff) Cleanup(state multistep.StateBag) {
	// no cleanup
}
