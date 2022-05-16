package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/packer-plugin-sakuracloud/platform"
)

type stepTransferArchive struct {
	Debug bool
}

func (s *stepTransferArchive) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	state.Put("transferred_ids", []types.ID{})
	state.Put("transferred_zones", []string{})

	c := state.Get("config").(Config)
	if len(c.Zones) == 0 {
		return multistep.ActionContinue
	}

	ui := state.Get("ui").(packer.Ui)
	stepStartMsg(ui, s.Debug, "TransferArchive")
	ui.Say("\tTransferring archive to other zones...")

	if err := s.transferArchives(ctx, state); err != nil {
		err := fmt.Errorf("Error creating archive: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	stepEndMsg(ui, s.Debug, "TransferArchive")
	return multistep.ActionContinue
}

func (s *stepTransferArchive) transferArchives(ctx context.Context, state multistep.StateBag) error {
	ui := state.Get("ui").(packer.Ui)
	archiveClient := state.Get("iaasClient").(*platform.Client).Archive
	c := state.Get("config").(Config)
	archiveID := state.Get("archive_id").(types.ID)

	var transferredIDs []types.ID
	var transferredZones []string

	// SAKURA Cloud API probably doesn't support multiple Transfer API calling, so we call Transfer API synchronously.
	for _, zone := range c.Zones {
		if c.Zone == zone {
			continue
		}
		archive, err := archiveClient.Transfer(ctx, zone, &platform.TransferArchiveRequest{
			Name:              c.ArchiveName,
			Tags:              c.ArchiveTags,
			Description:       c.ArchiveDescription,
			SourceArchiveID:   archiveID,
			SourceArchiveZone: c.Zone,
		})
		if err != nil {
			return err
		}
		ui.Say(fmt.Sprintf("\tArchive[%s:%s] is transferred from [%s:%s]", zone, archive.ID.String(), c.Zone, archiveID.String()))
		transferredIDs = append(transferredIDs, archive.ID)
		transferredZones = append(transferredZones, zone)
	}

	state.Put("transferred_ids", transferredIDs)
	state.Put("transferred_zones", transferredZones)
	return nil
}

func (s *stepTransferArchive) Cleanup(state multistep.StateBag) {
	// no cleanup
}
