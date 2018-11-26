package sakuracloud

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/stretchr/testify/assert"
)

func TestStepRemoteUpload(t *testing.T) {

	ioutil.WriteFile(dummyISOPath, []byte{}, 0644) // nolint
	defer func() {
		if _, err := os.Stat(dummyISOPath); err == nil {
			os.Remove(dummyISOPath)
		}
	}()

	ctx := context.Background()

	t.Run("with empty iso_path", func(t *testing.T) {
		step := &stepRemoteUpload{}
		state := initStepRemoteUpload(t)

		state.Put("iso_path", "")

		action := step.Run(ctx, state)
		assert.Equal(t, multistep.ActionContinue, action)
	})

	t.Run("with creating iso-image error", func(t *testing.T) {
		step := &stepRemoteUpload{}
		state := initStepRemoteUpload(t)

		client := state.Get("isoImageClient").(*dummyISOImageClient)
		client.newFunc = func() *sacloud.CDROM {
			return &sacloud.CDROM{}
		}
		client.createFunc = func(*sacloud.CDROM) (*sacloud.CDROM, *sacloud.FTPServer, error) {
			return nil, nil, errors.New("error")
		}

		action := step.Run(ctx, state)
		err := state.Get("error").(error)

		assert.Equal(t, multistep.ActionHalt, action)
		assert.Error(t, err)
	})

	t.Run("with connecting ftps server error", func(t *testing.T) {
		step := &stepRemoteUpload{}
		state := initStepRemoteUpload(t)

		client := state.Get("isoImageClient").(*dummyISOImageClient)
		client.newFunc = func() *sacloud.CDROM {
			return &sacloud.CDROM{}
		}
		client.createFunc = func(*sacloud.CDROM) (*sacloud.CDROM, *sacloud.FTPServer, error) {
			return createDummyISOImageWithFTPSInfo(dummyISOImageID, sacloud.EAAvailable)
		}

		ftps := state.Get("ftpsClient").(*dummyFTPSClient)
		ftps.connectFunc = func(string, int) error {
			return errors.New("error")
		}

		action := step.Run(ctx, state)
		err := state.Get("error").(error)

		assert.Equal(t, multistep.ActionHalt, action)
		assert.Error(t, err)
	})

	t.Run("with login error", func(t *testing.T) {
		step := &stepRemoteUpload{}
		state := initStepRemoteUpload(t)

		client := state.Get("isoImageClient").(*dummyISOImageClient)
		client.newFunc = func() *sacloud.CDROM {
			return &sacloud.CDROM{}
		}
		client.createFunc = func(*sacloud.CDROM) (*sacloud.CDROM, *sacloud.FTPServer, error) {
			return createDummyISOImageWithFTPSInfo(dummyISOImageID, sacloud.EAAvailable)
		}

		ftps := state.Get("ftpsClient").(*dummyFTPSClient)
		ftps.connectFunc = func(string, int) error {
			return nil
		}
		ftps.loginFunc = func(string, string) error {
			return errors.New("error")
		}

		action := step.Run(ctx, state)
		err := state.Get("error").(error)

		assert.Equal(t, multistep.ActionHalt, action)
		assert.Error(t, err)
	})

	t.Run("with saving remote file error", func(t *testing.T) {
		step := &stepRemoteUpload{}
		state := initStepRemoteUpload(t)

		client := state.Get("isoImageClient").(*dummyISOImageClient)
		client.newFunc = func() *sacloud.CDROM {
			return &sacloud.CDROM{}
		}
		client.createFunc = func(*sacloud.CDROM) (*sacloud.CDROM, *sacloud.FTPServer, error) {
			return createDummyISOImageWithFTPSInfo(dummyISOImageID, sacloud.EAAvailable)
		}

		ftps := state.Get("ftpsClient").(*dummyFTPSClient)
		ftps.connectFunc = func(string, int) error {
			return nil
		}
		ftps.loginFunc = func(string, string) error {
			return nil
		}
		ftps.storeFileFunc = func(string, *os.File) error {
			return errors.New("error")
		}

		action := step.Run(ctx, state)
		err := state.Get("error").(error)

		assert.Equal(t, multistep.ActionHalt, action)
		assert.Error(t, err)
	})

	t.Run("with closing ftps connection error", func(t *testing.T) {
		step := &stepRemoteUpload{}
		state := initStepRemoteUpload(t)

		client := state.Get("isoImageClient").(*dummyISOImageClient)
		client.newFunc = func() *sacloud.CDROM {
			return &sacloud.CDROM{}
		}
		client.createFunc = func(*sacloud.CDROM) (*sacloud.CDROM, *sacloud.FTPServer, error) {
			return createDummyISOImageWithFTPSInfo(dummyISOImageID, sacloud.EAAvailable)
		}
		client.closeFTPFunc = func(int64) (bool, error) {
			return false, errors.New("error")
		}

		ftps := state.Get("ftpsClient").(*dummyFTPSClient)
		ftps.connectFunc = func(string, int) error {
			return nil
		}
		ftps.loginFunc = func(string, string) error {
			return nil
		}
		ftps.storeFileFunc = func(string, *os.File) error {
			return nil
		}
		ftps.quitFunc = func() error {
			return nil
		}

		action := step.Run(ctx, state)
		err := state.Get("error").(error)

		assert.Equal(t, multistep.ActionHalt, action)
		assert.Error(t, err)
	})

	t.Run("normal case", func(t *testing.T) {
		step := &stepRemoteUpload{}
		state := initStepRemoteUpload(t)

		client := state.Get("isoImageClient").(*dummyISOImageClient)
		client.newFunc = func() *sacloud.CDROM {
			return &sacloud.CDROM{}
		}
		client.createFunc = func(*sacloud.CDROM) (*sacloud.CDROM, *sacloud.FTPServer, error) {
			return createDummyISOImageWithFTPSInfo(dummyISOImageID, sacloud.EAAvailable)
		}
		client.closeFTPFunc = func(int64) (bool, error) {
			return true, nil
		}

		ftps := state.Get("ftpsClient").(*dummyFTPSClient)
		ftps.connectFunc = func(string, int) error {
			return nil
		}
		ftps.loginFunc = func(string, string) error {
			return nil
		}
		ftps.storeFileFunc = func(string, *os.File) error {
			return nil
		}
		ftps.quitFunc = func() error {
			return nil
		}

		action := step.Run(ctx, state)
		isoID := state.Get("iso_id").(int64)
		_, hasError := state.GetOk("error")

		assert.Equal(t, multistep.ActionContinue, action)
		assert.Equal(t, dummyISOImageID, isoID)
		assert.False(t, hasError)
	})
}

func initStepRemoteUpload(t *testing.T) multistep.StateBag {
	config := dummyConfig()
	state := dummyMinimumStateBag(&config)
	state.Put("iso_path", dummyISOPath)

	state.Put("isoImageClient", &dummyISOImageClient{
		findFunc: func() (*sacloud.SearchResponse, error) {
			t.Fail()
			return nil, nil
		},
		newFunc: func() *sacloud.CDROM {
			t.Fail()
			return nil
		},
		createFunc: func(*sacloud.CDROM) (*sacloud.CDROM, *sacloud.FTPServer, error) {
			t.Fail()
			return nil, nil, nil
		},
		closeFTPFunc: func(int64) (bool, error) {
			t.Fail()
			return false, nil
		},
	})

	state.Put("ftpsClient", &dummyFTPSClient{
		connectFunc: func(string, int) error {
			t.Fail()
			return nil
		},
		loginFunc: func(string, string) error {
			t.Fail()
			return nil
		},
		storeFileFunc: func(string, *os.File) error {
			t.Fail()
			return nil
		},
		quitFunc: func() error {
			t.Fail()
			return nil
		},
	})

	return state
}

func createDummyISOImageWithFTPSInfo(id int64, availability sacloud.EAvailability) (*sacloud.CDROM, *sacloud.FTPServer, error) {
	isoImage := createDummyISOImage(id, availability)
	ftpsInfo := &sacloud.FTPServer{
		HostName:  "ftps.example.com",
		IPAddress: "192.2.0.1",
		User:      "example",
		Password:  "p@ssw0rd",
	}
	return isoImage, ftpsInfo, nil
}
