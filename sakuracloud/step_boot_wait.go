package sakuracloud

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"time"
)

// stepBootWait waits the configured time period.
type stepBootWait struct {
	Debug bool
}

func (s *stepBootWait) Run(state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)

	stepStartMsg(ui, s.Debug, "BootWait")

	if int64(config.BootWait) > 0 {
		ui.Say(fmt.Sprintf("\tWaiting %s for boot...", config.BootWait))
		time.Sleep(config.BootWait)
	}

	stepEndMsg(ui, s.Debug, "BootWait")
	return multistep.ActionContinue
}

func (s *stepBootWait) Cleanup(state multistep.StateBag) {}
