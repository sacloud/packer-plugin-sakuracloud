package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/packer-builder-sakuracloud/iaas"
)

type stepServerInfo struct {
	Debug bool
}

func (s *stepServerInfo) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	c := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)

	caller := state.Get("iaasClient").(iaas.Client).Caller
	serverOp := sacloud.NewServerOp(caller)
	serverID := state.Get("server_id").(types.ID)

	stepStartMsg(ui, s.Debug, "Read Server Info")
	ui.Say("\tWaiting for server to become active...")

	// Set the Network informations on the state for later
	server, err := serverOp.Read(ctx, c.Zone, serverID)
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

	if len(server.Interfaces) > 0 && !server.Interfaces[0].SwitchID.IsEmpty() {
		iface := server.Interfaces[0]
		state.Put("server_ip", iface.IPAddress)
		state.Put("default_route", iface.SubnetDefaultRoute)
		state.Put("network_mask_len", iface.SubnetNetworkMaskLen)
	}
	if len(server.Zone.Region.NameServers) > 0 {
		state.Put("dns1", server.Zone.Region.NameServers[0])
	}
	if len(server.Zone.Region.NameServers) > 1 {
		state.Put("dns2", server.Zone.Region.NameServers[1])
	}

	// Set the VNC proxy on the state for later
	vnc, err := serverOp.GetVNCProxy(ctx, c.Zone, serverID)
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
