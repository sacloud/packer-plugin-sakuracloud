package sakuracloud

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/sacloud/packer-builder-sakuracloud/iaas"
	"github.com/sacloud/packer-builder-sakuracloud/sakuracloud/constants"
)

type stepRemoteUpload struct {
	Debug bool
}

func (s *stepRemoteUpload) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	isoImageClient := state.Get("isoImageClient").(iaas.ISOImageClient)
	ui := state.Get("ui").(packer.Ui)

	stepStartMsg(ui, s.Debug, "ISO-Image Upload")

	filepath, ok := state.Get("iso_path").(string)
	if !ok || filepath == "" {
		return multistep.ActionContinue
	}

	config := state.Get("config").(Config)
	checksum := config.ISOChecksum

	if checksum != "" {
		//search ISO from SakuraCloud
		isoImageClient.SetEmpty()
		isoImageClient.SetNameLike(checksum)
		res, err := isoImageClient.Find()
		if err != nil {
			err := fmt.Errorf("Error finding ISO image: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
		if len(res.CDROMs) > 0 {
			state.Put("iso_id", res.CDROMs[0].ID)
			stepEndMsg(ui, s.Debug, "ISO-Image Upload")
			return multistep.ActionContinue
		}
	}

	ui.Say("\tUploading ISO to SakuraCloud...")
	log.Printf("Remote uploading: %s", filepath)

	req := isoImageClient.New()
	req.Name = config.ISOImageName
	req.SizeMB = config.ISOImageSizeGB * 1024
	req.Description = strings.Join(config.ISOConfig.ISOUrls, "\n")
	req.AppendTag(constants.UploadedFromPackerMarkerTag)

	isoImage, ftp, err := isoImageClient.Create(req)
	if err != nil {
		err := fmt.Errorf("Error creating ISO image: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// upload iso by FTPS
	ftpsClient := state.Get("ftpsClient").(iaas.FTPSClient)

	err = ftpsClient.Connect(ftp.HostName, 21)
	if err != nil {
		err := fmt.Errorf("Error connecting FTPS server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt

	}

	err = ftpsClient.Login(ftp.User, ftp.Password)
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

	err = ftpsClient.StoreFile("packer-for-sakuracloud.iso", fs)
	if err != nil {
		err := fmt.Errorf("Error store file on FTPS server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	ftpsClient.Quit()

	// close image FTP after upload
	_, err = isoImageClient.CloseFTP(isoImage.ID)
	if err != nil {
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
