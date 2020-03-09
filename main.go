package main

import (
	"github.com/hashicorp/packer/packer/plugin"
	"github.com/sacloud/packer-builder-sakuracloud/sakuracloud"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterBuilder(new(sakuracloud.Builder))
	server.Serve()
}
