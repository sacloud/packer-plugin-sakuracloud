package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/sacloud/libsacloud/api"
)

type stepCreateArchive struct {
	Debug bool
}

func (s *stepCreateArchive) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*api.Client)
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)
	diskID := state.Get("disk_id").(int64)

	stepStartMsg(ui, s.Debug, "CreateArchive")

	ui.Say(fmt.Sprintf("\tCreating archive: %v", c.ArchiveName))

	archiveReq := client.Archive.New()
	archiveReq.Name = c.ArchiveName
	archiveReq.SetSourceDisk(diskID)

	if len(c.ArchiveTags) > 0 {
		for _, tag := range c.ArchiveTags {
			if !archiveReq.HasTag(tag) {
				archiveReq.AppendTag(tag)
			}
		}
	} else {
		publicArchiveID, found := client.Disk.GetPublicArchiveIDFromAncestors(diskID)
		if found {
			sourceArchive, err := client.Archive.Read(publicArchiveID)
			if sourceArchive != nil && err == nil {
				for _, tag := range sourceArchive.Tags {
					if !archiveReq.HasTag(tag) {
						archiveReq.AppendTag(tag)
					}
				}
			}
		}
	}

	archiveReq.Description = c.ArchiveDescription

	archive, err := client.Archive.Create(archiveReq)
	if err != nil {
		err := fmt.Errorf("Error creating archive: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if err := client.Archive.SleepWhileCopying(archive.ID, c.APIClientTimeout); err != nil {
		err := fmt.Errorf("Error copying archive: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("archive_id", archive.ID)
	state.Put("archive_name", archive.Name)
	state.Put("zone", client.Zone)
	stepEndMsg(ui, s.Debug, "BootWait")
	return multistep.ActionContinue
}

func (s *stepCreateArchive) Cleanup(state multistep.StateBag) {
	// no cleanup
}
