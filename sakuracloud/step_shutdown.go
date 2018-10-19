package sakuracloud

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/sacloud/libsacloud/api"
)

type stepShutdown struct {
	Debug bool
}

func (s *stepShutdown) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*api.Client)
	ui := state.Get("ui").(packer.Ui)
	serverID := state.Get("server_id").(int64)

	stepStartMsg(ui, s.Debug, "Shutdown Server")

	ui.Say("\tGracefully shutting down server...")

	_, err := client.Server.Shutdown(serverID)
	if err != nil {
		err := fmt.Errorf("Error shutting down server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	client.Server.SleepUntilDown(serverID, 1*time.Minute)

	stepEndMsg(ui, s.Debug, "Shutdown Server")
	return multistep.ActionContinue
}

func (s *stepShutdown) Cleanup(state multistep.StateBag) {
	// no cleanup
}
