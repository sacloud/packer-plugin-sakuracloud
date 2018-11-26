package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/packer-builder-sakuracloud/iaas"
)

type stepCreateArchive struct {
	Debug bool
}

func (s *stepCreateArchive) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	archiveClient := state.Get("archiveClient").(iaas.ArchiveClient)
	basicClient := state.Get("basicClient").(iaas.BasicClient)

	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)

	stepStartMsg(ui, s.Debug, "CreateArchive")

	ui.Say(fmt.Sprintf("\tCreating archive: %v", c.ArchiveName))

	archiveReq := s.createArchiveReq(state)

	archive, err := archiveClient.Create(archiveReq)
	if err != nil {
		err := fmt.Errorf("Error creating archive: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if err := archiveClient.SleepWhileCopying(archive.ID, c.APIClientTimeout); err != nil {
		err := fmt.Errorf("Error copying archive: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("archive_id", archive.ID)
	state.Put("archive_name", archive.Name)
	state.Put("zone", basicClient.Zone())
	stepEndMsg(ui, s.Debug, "BootWait")
	return multistep.ActionContinue
}

func (s *stepCreateArchive) createArchiveReq(state multistep.StateBag) *sacloud.Archive {
	archiveClient := state.Get("archiveClient").(iaas.ArchiveClient)
	diskClient := state.Get("diskClient").(iaas.DiskClient)

	c := state.Get("config").(Config)
	diskID := state.Get("disk_id").(int64)

	archiveReq := archiveClient.New()
	archiveReq.Name = c.ArchiveName
	archiveReq.SetSourceDisk(diskID)

	if len(c.ArchiveTags) > 0 {
		for _, tag := range c.ArchiveTags {
			if !archiveReq.HasTag(tag) {
				archiveReq.AppendTag(tag)
			}
		}
	} else {
		publicArchiveID, found := diskClient.GetPublicArchiveIDFromAncestors(diskID)
		if found {
			sourceArchive, err := archiveClient.Read(publicArchiveID)
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

	return archiveReq
}

func (s *stepCreateArchive) Cleanup(state multistep.StateBag) {
	// no cleanup
}
