package sakuracloud

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/pkg/errors"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/stretchr/testify/assert"
)

func TestStepShutdown(t *testing.T) {

	ctx := context.Background()

	t.Run("with shutdown error", func(t *testing.T) {
		step := &stepShutdown{}
		state := initStepShutdownState(t)

		serverClient := state.Get("serverClient").(*dummyServerClient)
		serverClient.shutdownFunc = func(int64) (bool, error) {
			return false, errors.New("error")
		}

		action := step.Run(ctx, state)
		err, hasError := state.GetOk("error")

		assert.Equal(t, multistep.ActionHalt, action)
		assert.True(t, hasError)
		assert.Error(t, err.(error))
	})

	t.Run("ignore sleep error", func(t *testing.T) {
		step := &stepShutdown{}
		state := initStepShutdownState(t)

		serverClient := state.Get("serverClient").(*dummyServerClient)
		serverClient.shutdownFunc = func(int64) (bool, error) {
			return true, nil
		}
		serverClient.sleepUntilDownFunc = func(int64, time.Duration) error {
			return errors.New("error")
		}

		action := step.Run(ctx, state)
		_, hasError := state.GetOk("error")

		assert.Equal(t, multistep.ActionContinue, action)
		assert.False(t, hasError)
	})

	t.Run("normal case", func(t *testing.T) {
		step := &stepShutdown{}
		state := initStepShutdownState(t)

		serverClient := state.Get("serverClient").(*dummyServerClient)
		serverClient.shutdownFunc = func(int64) (bool, error) {
			return true, nil
		}
		serverClient.sleepUntilDownFunc = func(int64, time.Duration) error {
			return nil
		}

		action := step.Run(ctx, state)
		_, hasError := state.GetOk("error")

		assert.Equal(t, multistep.ActionContinue, action)
		assert.False(t, hasError)
	})
}

func initStepShutdownState(t *testing.T) multistep.StateBag {
	state := dummyMinimumStateBag(nil)
	state.Put("server_id", dummyServerID)

	state.Put("serverClient", &dummyServerClient{
		readFunc: func(int64) (*sacloud.Server, error) {
			t.Fail()
			return nil, nil
		},
		stopFunc: func(int64) (bool, error) {
			t.Fail()
			return false, nil
		},
		sleepUntilDownFunc: func(int64, time.Duration) error {
			t.Fail()
			return nil
		},
	})

	return state
}
