package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/packer-plugin-sakuracloud/iaas"
)

type stepPrepareISO struct {
	Debug bool
}

func (s *stepPrepareISO) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	c := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)

	caller := state.Get("iaasClient").(*iaas.Client).Caller
	isoImageOp := sacloud.NewCDROMOp(caller)

	stepStartMsg(ui, s.Debug, "PrepareISO")

	if c.ISOImageID.IsEmpty() {
		c.ISOImageID = state.Get("iso_id").(types.ID)
	}

	image, err := isoImageOp.Read(ctx, c.Zone, c.ISOImageID)
	if err != nil {
		err := fmt.Errorf("Error invalid ISO image ID: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if !image.Availability.IsAvailable() {
		err := fmt.Errorf("Error invalid ISO image Status: %s", image.Availability)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("config", c)
	stepEndMsg(ui, s.Debug, "PrepareISO")
	return multistep.ActionContinue
}

func (s *stepPrepareISO) Cleanup(state multistep.StateBag) {
}
