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

	stepStartMsg(ui, s.Debug, "Read Server Info")

	ui.Say("\tWaiting for server to become active...")

	// Set the Network informations on the state for later
	server, err := client.Server.Read(serverID)
	if err != nil {
		err := fmt.Errorf("Error retrieving server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("server_ip", "")
	state.Put("default_route", "")
	state.Put("network_mask_len", "")
	state.Put("dns1", "")
	state.Put("dns2", "")

	if len(server.Interfaces) > 0 && server.Interfaces[0].Switch != nil {
		ip := server.Interfaces[0].IPAddress
		if server.Interfaces[0].Switch.Scope != sacloud.ESCopeShared {
			ip = server.Interfaces[0].UserIPAddress
		}
		state.Put("server_ip", ip)

		state.Put("default_route", server.Interfaces[0].Switch.UserSubnet.DefaultRoute)
		state.Put("network_mask_len", server.Interfaces[0].Switch.UserSubnet.NetworkMaskLen)
	}
	if len(server.Zone.Region.NameServers) > 0 {
		state.Put("dns1", server.Zone.Region.NameServers[0])
	}
	if len(server.Zone.Region.NameServers) > 1 {
		state.Put("dns2", server.Zone.Region.NameServers[1])
	}

	// Set the VNC proxy on the state for later
	vnc, err := client.Server.GetVNCProxy(serverID)
	if err != nil {
		err := fmt.Errorf("Error vnc proxy info: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	state.Put("vnc", vnc)

	stepEndMsg(ui, s.Debug, "Read Server Info")
	return multistep.ActionContinue
}

func (s *stepServerInfo) Cleanup(state multistep.StateBag) {
	// no cleanup
}
