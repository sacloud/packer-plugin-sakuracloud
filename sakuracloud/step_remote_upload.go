package sakuracloud

import (
	"bufio"
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"github.com/sacloud/libsacloud/api"
	"github.com/webguerilla/ftps"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type stepRemoteUpload struct{}

func (s *stepRemoteUpload) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*api.Client)
	ui := state.Get("ui").(packer.Ui)

	filepath, ok := state.Get("iso_path").(string)
	if !ok {
		return multistep.ActionContinue
	}

	config := state.Get("config").(Config)
	checksum := config.ISOChecksum

	if checksum != "" {
		//search ISO from SakuraCloud
		res, err := client.CDROM.Reset().WithNameLike(checksum).Find()
		if err != nil {
			err := fmt.Errorf("Error finding ISO image: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
		if len(res.CDROMs) > 0 {
			state.Put("iso_id", res.CDROMs[0].ID)
			return multistep.ActionContinue
		}
	}

	ui.Say("Uploading ISO to SakuraCloud...")
	log.Printf("Remote uploading: %s", filepath)

	req := client.CDROM.New()
	req.Name = config.ISOImageName
	req.SizeMB = config.ISOImageSizeGB * 1024
	req.Description = strings.Join(config.ISOConfig.ISOUrls, "\n")
	req.AppendTag("packer-for-sakuracloud")

	isoImage, ftp, err := client.CDROM.Create(req)
	if err != nil {
		err := fmt.Errorf("Error creating ISO image: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// upload iso by FTPS
	ftpsClient := &ftps.FTPS{}
	ftpsClient.TLSConfig.InsecureSkipVerify = true

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

	reader := bufio.NewReader(fs)
	fileBytes, _ := ioutil.ReadAll(reader)

	err = ftpsClient.StoreFile("packer-for-sakuracloud.iso", fileBytes)
	if err != nil {
		err := fmt.Errorf("Error store file on FTPS server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	ftpsClient.Quit()

	// close image FTP after upload
	_, err = client.CDROM.CloseFTP(isoImage.ID)
	if err != nil {
		err := fmt.Errorf("Error Closing FTPS connection: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("iso_id", isoImage.ID)
	return multistep.ActionContinue
}

func (s *stepRemoteUpload) Cleanup(state multistep.StateBag) {
}
