package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/packer-builder-sakuracloud/iaas"
)

type stepCreateArchive struct {
	Debug bool
}

func (s *stepCreateArchive) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	stepStartMsg(ui, s.Debug, "CreateArchive")
	ui.Say("\tCreating archive...")

	archive, err := s.createArchive(ctx, state)
	if err != nil {
		err := fmt.Errorf("Error creating archive: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("archive_id", archive.ID)
	state.Put("archive_name", archive.Name)
	stepEndMsg(ui, s.Debug, "CreateArchive")
	return multistep.ActionContinue
}

func (s *stepCreateArchive) createArchive(ctx context.Context, state multistep.StateBag) (*sacloud.Archive, error) {
	archiveClient := state.Get("archiveClient").(iaas.Archive)
	c := state.Get("config").(Config)
	diskID := state.Get("disk_id").(types.ID)

	req := &iaas.CreateArchiveRequest{
		DiskID:      diskID,
		Name:        c.ArchiveName,
		Tags:        c.ArchiveTags,
		Description: c.ArchiveDescription,
	}
	return archiveClient.Create(ctx, req)
}

func (s *stepCreateArchive) Cleanup(state multistep.StateBag) {
	// no cleanup
}
