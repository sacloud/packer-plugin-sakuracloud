package sakuracloud

import (
	"fmt"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"github.com/sacloud/libsacloud/api"
)

type stepCreateArchive struct {
	Debug bool
}

func (s *stepCreateArchive) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*api.Client)
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)
	diskID := state.Get("disk_id").(int64)

	ui.Say(fmt.Sprintf("Creating archive: %v", c.ArchiveName))

	archiveReq := client.Archive.New()
	archiveReq.Name = c.ArchiveName
	archiveReq.SetSourceDisk(diskID)
	for _, tag := range c.ArchiveTags {
		if !archiveReq.HasTag(tag) {
			archiveReq.AppendTag(tag)
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

	if err := client.Archive.SleepWhileCopying(archive.ID, c.StateTimeout); err != nil {
		err := fmt.Errorf("Error copying archive: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("archive_id", archive.ID)
	state.Put("archive_name", archive.Name)
	state.Put("zone", client.Zone)
	return multistep.ActionContinue
}

func (s *stepCreateArchive) Cleanup(state multistep.StateBag) {
	// no cleanup
}
