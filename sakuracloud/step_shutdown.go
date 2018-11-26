package sakuracloud

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/sacloud/packer-builder-sakuracloud/iaas"
)

type stepShutdown struct {
	Debug bool
}

func (s *stepShutdown) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	serverClient := state.Get("serverClient").(iaas.ServerClient)
	ui := state.Get("ui").(packer.Ui)
	serverID := state.Get("server_id").(int64)

	stepStartMsg(ui, s.Debug, "Shutdown Server")

	ui.Say("\tGracefully shutting down server...")

	_, err := serverClient.Shutdown(serverID)
	if err != nil {
		err := fmt.Errorf("Error shutting down server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// if occur timeout, it be force shutdown by step_power_off.
	// so ignore any errors to fallback early.
	if err := serverClient.SleepUntilDown(serverID, 1*time.Minute); err != nil {
		ui.Message("Graceful shutdown is timed out")
	}

	stepEndMsg(ui, s.Debug, "Shutdown Server")
	return multistep.ActionContinue
}

func (s *stepShutdown) Cleanup(state multistep.StateBag) {
	// no cleanup
}
