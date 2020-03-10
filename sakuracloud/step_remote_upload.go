package sakuracloud

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/sacloud/libsacloud/v2/pkg/size"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/search"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/packer-builder-sakuracloud/iaas"
	"github.com/sacloud/packer-builder-sakuracloud/sakuracloud/constants"
)

type stepRemoteUpload struct {
	Debug bool
}

func (s *stepRemoteUpload) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	c := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)

	client := state.Get("iaasClient").(*iaas.Client)
	caller := client.Caller
	isoImageOp := sacloud.NewCDROMOp(caller)

	stepStartMsg(ui, s.Debug, "ISO-Image Upload")

	filepath, ok := state.Get("iso_path").(string)
	if !ok || filepath == "" {
		return multistep.ActionContinue
	}

	if c.ISOChecksum != "" {
		searched, err := isoImageOp.Find(ctx, c.Zone, &sacloud.FindCondition{
			Filter: search.Filter{
				search.Key("Name"): search.ExactMatch(c.ISOChecksum),
			},
		})
		if err != nil {
			err := fmt.Errorf("Error finding ISO image: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
		if len(searched.CDROMs) > 0 {
			state.Put("iso_id", searched.CDROMs[0].ID)
			stepEndMsg(ui, s.Debug, "ISO-Image Upload")
			return multistep.ActionContinue
		}
	}

	ui.Say("\tUploading ISO to SakuraCloud...")
	log.Printf("Remote uploading: %s", filepath)

	isoImage, ftp, err := isoImageOp.Create(ctx, c.Zone, &sacloud.CDROMCreateRequest{
		SizeMB:      size.GiBToMiB(c.ISOImageSizeGB),
		Name:        c.ISOImageName,
		Description: strings.Join(c.ISOConfig.ISOUrls, "\n"),
		Tags:        types.Tags{constants.UploadedFromPackerMarkerTag},
	})
	if err != nil {
		err := fmt.Errorf("Error creating ISO image: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// upload iso by FTPS
	err = client.FTPS.Connect(ftp.HostName, 21)
	if err != nil {
		err := fmt.Errorf("Error connecting FTPS server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	err = client.FTPS.Login(ftp.User, ftp.Password)
	if err != nil {
		err := fmt.Errorf("Error login FTPS server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	fs, err := os.Open(filepath)
	if err != nil {
		err := fmt.Errorf("Error opening ISO file: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	defer fs.Close()

	err = client.FTPS.StoreFile("packer-for-sakuracloud.iso", fs)
	if err != nil {
		err := fmt.Errorf("Error store file on FTPS server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	if err := client.FTPS.Quit(); err != nil {
		err := fmt.Errorf("Error quit on FTPS server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// close image FTP after upload
	if err := isoImageOp.CloseFTP(ctx, c.Zone, isoImage.ID); err != nil {
		err := fmt.Errorf("Error Closing FTPS connection: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("iso_id", isoImage.ID)
	stepEndMsg(ui, s.Debug, "ISO-Image Upload")
	return multistep.ActionContinue
}

func (s *stepRemoteUpload) Cleanup(state multistep.StateBag) {
}
