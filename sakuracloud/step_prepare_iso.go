package sakuracloud

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"github.com/sacloud/libsacloud/api"
)

type stepPrepareISO struct{}

func (s *stepPrepareISO) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*api.Client)
	config := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)

	isoID := config.ISOImageID
	if isoID == 0 {
		isoID = state.Get("iso_id").(int64)
	}

	config.ISOImageID = isoID

	iso, err := client.CDROM.Read(isoID)
	if err != nil {
		err := fmt.Errorf("Error invalid ISO image ID: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	if !iso.IsAvailable() {
		err := fmt.Errorf("Error invalid ISO image Status: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("config", config)
	return multistep.ActionContinue
}

func (s *stepPrepareISO) Cleanup(state multistep.StateBag) {
}
