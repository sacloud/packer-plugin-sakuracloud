package sakuracloud

import (
	"context"
	"testing"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/pkg/errors"
	"github.com/sacloud/libsacloud/builder"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/stretchr/testify/assert"
)

var dummyFactory = &dummyBuilderFactory{
	builder: dummyBuilder,
}

var dummyBuilder = &dummyServerBuilder{}

func TestStepCreateServer_Run(t *testing.T) {

	t.Run("with build error", func(t *testing.T) {
		ctx := context.Background()
		state := initStepCreateServerState()
		step := &stepCreateServer{}

		dummyBuilder.init()
		dummyBuilder.err = errors.New("err")

		action := step.Run(ctx, state)
		err, ok := state.GetOk("error")
		assert.True(t, ok)
		assert.Error(t, err.(error))
		assert.Equal(t, multistep.ActionHalt, action)
	})

	t.Run("normal case", func(t *testing.T) {
		ctx := context.Background()
		state := initStepCreateServerState()
		step := &stepCreateServer{}

		dummyBuilder.init()
		dummyBuilder.result = &builder.ServerBuildResult{
			Server: &sacloud.Server{Resource: sacloud.NewResource(dummyServerID)},
			Disks: []*builder.DiskBuildResult{
				{
					Disk: &sacloud.Disk{Resource: sacloud.NewResource(dummyDiskID)},
				},
			},
		}

		action := step.Run(ctx, state)
		_, errExists := state.GetOk("error")
		serverID, serverIDExists := state.GetOk("server_id")
		diskID, diskIDExists := state.GetOk("disk_id")

		assert.False(t, errExists)
		assert.Equal(t, dummyServerID, serverID.(int64))
		assert.True(t, serverIDExists)
		assert.Equal(t, dummyDiskID, diskID.(int64))
		assert.True(t, diskIDExists)

		assert.Equal(t, multistep.ActionContinue, action)

	})
}

func initStepCreateServerState() multistep.StateBag {
	config := dummyConfig()
	state := dummyMinimumStateBag(&config)
	state.Put("builderFactory", dummyFactory)
	return state
}
