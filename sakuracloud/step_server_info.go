package sakuracloud

import (
	"fmt"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

type stepServerInfo struct {
	Debug bool
}

func (s *stepServerInfo) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*api.Client)
	ui := state.Get("ui").(packer.Ui)
	serverID := state.Get("server_id").(int64)

	ui.Say("Waiting for server to become active...")

	// Set the IP on the state for later
	server, err := client.Server.Read(serverID)
	if err != nil {
		err := fmt.Errorf("Error retrieving server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ip := server.Interfaces[0].IPAddress
	if server.Interfaces[0].Switch.Scope != sacloud.ESCopeShared {
		ip = server.Interfaces[0].UserIPAddress
	}
	state.Put("server_ip", ip)

	return multistep.ActionContinue
}

func (s *stepServerInfo) Cleanup(state multistep.StateBag) {
	// no cleanup
}
