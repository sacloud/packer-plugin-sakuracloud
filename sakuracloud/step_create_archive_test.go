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

func testReadArchive(id int64) *sacloud.Archive {
	archive := &sacloud.Archive{}
	archive.Name = "test"
	archive.Resource = sacloud.NewResource(id)
	archive.Tags = dummyParentArchiveTags
	return archive
}

func TestStepCreateArchive_NormalCase(t *testing.T) {
	state := initStepCreateArchiveState()
	step := &stepCreateArchive{}
	ctx := context.Background()

	action := step.Run(ctx, state)

	assert.Equal(t, multistep.ActionContinue, action)
	assert.Equal(t, dummyCreatedArchiveID, state.Get("archive_id").(int64))
	assert.Equal(t, dummyArchiveName, state.Get("archive_name").(string))
	assert.Equal(t, "is1a", state.Get("zone").(string))
}

func TestStepCreateArchiveState_CreateArchive(t *testing.T) {
	t.Run("exists archive tag parameter", func(t *testing.T) {
		state := initStepCreateArchiveState()
		step := &stepCreateArchive{}

		archiveReq := step.createArchiveReq(state)

		assert.Equal(t, dummyArchiveName, archiveReq.Name)
		assert.Equal(t, dummyDiskID, archiveReq.SourceDisk.ID)
		assert.EqualValues(t, dummyArchiveTags, archiveReq.Tags)
		assert.Equal(t, dummyDescription, archiveReq.Description)
	})

	t.Run("empty archive tag parameter", func(t *testing.T) {
		state := initStepCreateArchiveState()
		step := &stepCreateArchive{}
		config := state.Get("config").(Config)
		config.ArchiveTags = []string{}
		state.Put("config", config)

		archiveReq := step.createArchiveReq(state)

		assert.Equal(t, dummyArchiveName, archiveReq.Name)
		assert.Equal(t, dummyDiskID, archiveReq.SourceDisk.ID)
		assert.EqualValues(t, dummyParentArchiveTags, archiveReq.Tags)
		assert.Equal(t, dummyDescription, archiveReq.Description)
	})
}

func TestStepCreateArchiveState_ErrorCreateArchive(t *testing.T) {

	state := initStepCreateArchiveState()
	step := &stepCreateArchive{}
	ctx := context.Background()

	archiveClient := state.Get("archiveClient").(*dummyArchiveClient)
	archiveClient.createFunc = func(*sacloud.Archive) (*sacloud.Archive, error) {
		return nil, errors.New("error")
	}

	action := step.Run(ctx, state)
	err, ok := state.GetOk("error")

	assert.Equal(t, multistep.ActionHalt, action)
	assert.True(t, ok)
	assert.NotNil(t, err)
}

func TestStepCreateArchiveState_ErrorSleepWhileCopying(t *testing.T) {

	state := initStepCreateArchiveState()
	step := &stepCreateArchive{}
	ctx := context.Background()

	archiveClient := state.Get("archiveClient").(*dummyArchiveClient)
	archiveClient.sleepWhileCopyingFunc = func(int64, time.Duration) error {
		return errors.New("error")
	}

	action := step.Run(ctx, state)
	err, ok := state.GetOk("error")

	assert.Equal(t, multistep.ActionHalt, action)
	assert.True(t, ok)
	assert.NotNil(t, err)
}

func initStepCreateArchiveState() multistep.StateBag {

	config := dummyConfig()
	config.ArchiveName = dummyArchiveName
	config.ArchiveTags = dummyArchiveTags
	config.ArchiveDescription = dummyDescription
	state := dummyMinimumStateBag(&config)
	state.Put("disk_id", dummyDiskID)

	state.Put("archiveClient", &dummyArchiveClient{
		createFunc: func(param *sacloud.Archive) (*sacloud.Archive, error) {
			param.Resource = sacloud.NewResource(dummyCreatedArchiveID)
			return param, nil
		},
		readFunc: func(id int64) (*sacloud.Archive, error) {
			archive := testReadArchive(id)
			return archive, nil
		},
	})
	state.Put("diskClient", &dummyDiskClient{
		getPublicArchiveIDFromAncestorsFunc: func(id int64) (int64, bool) {
			return dummyReadArchiveID, true
		},
	})
	state.Put("basicClient", &dummyBasicClient{})

	return state
}
