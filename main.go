package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
	"github.com/hashicorp/packer-plugin-sdk/version"
	"github.com/sacloud/packer-plugin-sakuracloud/sakuracloud"
	pversion "github.com/sacloud/packer-plugin-sakuracloud/version"
)

var pluginVersion = version.InitializePluginVersion(pversion.Version, "")

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder(plugin.DEFAULT_NAME, new(sakuracloud.Builder))
	pps.SetVersion(pluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
