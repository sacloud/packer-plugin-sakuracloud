package sakuracloud

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/stretchr/testify/assert"
)

func TestStepPrepareISO(t *testing.T) {

	ctx := context.Background()

	t.Run("with read error", func(t *testing.T) {
		step := &stepPrepareISO{}
		state := initStepPrepareISOState(t)

		client := state.Get("isoImageClient").(*dummyISOImageClient)
		client.readFunc = func(id int64) (*sacloud.CDROM, error) {
			return nil, errors.New("error")
		}

		action := step.Run(ctx, state)
		err := state.Get("error").(error)

		assert.Equal(t, multistep.ActionHalt, action)
		assert.Error(t, err)
	})

	t.Run("with unavailable iso-image", func(t *testing.T) {
		step := &stepPrepareISO{}
		state := initStepPrepareISOState(t)

		client := state.Get("isoImageClient").(*dummyISOImageClient)
		client.readFunc = func(id int64) (*sacloud.CDROM, error) {
			return createDummyISOImage(id, sacloud.EAFailed), nil
		}

		action := step.Run(ctx, state)
		err := state.Get("error").(error)

		assert.Equal(t, multistep.ActionHalt, action)
		assert.Error(t, err)
	})

	t.Run("normal case", func(t *testing.T) {
		step := &stepPrepareISO{}
		state := initStepPrepareISOState(t)

		client := state.Get("isoImageClient").(*dummyISOImageClient)
		client.readFunc = func(id int64) (*sacloud.CDROM, error) {
			return createDummyISOImage(id, sacloud.EAAvailable), nil
		}

		action := step.Run(ctx, state)
		_, hasError := state.GetOk("error")

		assert.Equal(t, multistep.ActionContinue, action)
		assert.False(t, hasError)
	})
}

func initStepPrepareISOState(t *testing.T) multistep.StateBag {
	state := dummyMinimumStateBag(nil)
	state.Put("iso_id", dummyISOImageID)

	state.Put("isoImageClient", &dummyISOImageClient{
		readFunc: func(int64) (*sacloud.CDROM, error) {
			t.Fail()
			return nil, nil
		},
	})

	return state
}

func createDummyISOImage(id int64, availability sacloud.EAvailability) *sacloud.CDROM {
	isoImage := &sacloud.CDROM{Resource: sacloud.NewResource(id)}
	isoImage.Availability = availability
	return isoImage
}
