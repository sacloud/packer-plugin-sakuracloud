package sakuracloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/sacloud/libsacloud/v2/sacloud/ostype"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/packer-plugin-sakuracloud/iaas"
	"github.com/sacloud/packer-plugin-sakuracloud/sakuracloud/constants"
)

// BuilderId is the unique id for the builder
const BuilderId = "packer.sakuracloud"

// Builder implememts packer.Builder interface for
// handling actions for SakuraCloud
type Builder struct {
	config Config
	runner multistep.Runner
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec {
	return b.config.FlatMapstructure().HCL2Spec()
}

// Prepare is responsible for configuring the builder and validating
// that configuration.
func (b *Builder) Prepare(raws ...interface{}) ([]string, []string, error) {
	c, warnings, errs := NewConfig(raws...)
	if errs != nil {
		return nil, warnings, errs
	}
	b.config = *c

	return nil, nil, nil
}

// Run is where the actual build should take place.
func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {
	client := iaas.NewClient(b.config.AccessToken, b.config.AccessTokenSecret, b.config.Zone)

	// Set up the state
	state := new(multistep.BasicStateBag)
	state.Put("config", b.config)

	state.Put("iaasClient", client)

	state.Put("hook", hook)
	state.Put("ui", ui)

	// Build the steps
	var steps []multistep.Step

	var communicateStep multistep.Step
	os := ostype.StrToOSType(b.config.OSType)

	getHostIPFunc := func(state multistep.StateBag) (string, error) {
		ipAddress := state.Get("server_ip").(string)
		return ipAddress, nil
	}

	if os.IsWindows() {
		communicateStep = &communicator.StepConnectWinRM{
			Config: &b.config.Comm,
			Host:   getHostIPFunc,
			WinRMConfig: func(multistep.StateBag) (*communicator.WinRMConfig, error) {
				return &communicator.WinRMConfig{
					Username: b.config.UserName,
					Password: b.config.Password,
				}, nil
			},
		}
	} else {
		communicateStep = &communicator.StepConnectSSH{
			Config:    &b.config.Comm,
			Host:      getHostIPFunc,
			SSHConfig: b.config.Comm.SSHConfigFunc(),
		}
	}

	var isoSteps []multistep.Step
	if b.config.ISOImageID > 0 {
		isoSteps = []multistep.Step{
			&stepPrepareISO{
				Debug: b.config.PackerDebug,
			},
		}
	} else {
		isoSteps = []multistep.Step{
			&commonsteps.StepDownload{
				Checksum:    b.config.ISOChecksum,
				Description: "ISO",
				Extension:   "iso",
				ResultKey:   "iso_path",
				TargetPath:  b.config.TargetPath,
				Url:         b.config.ISOUrls,
			},
			&stepRemoteUpload{
				Debug: b.config.PackerDebug,
			},
			&stepPrepareISO{
				Debug: b.config.PackerDebug,
			},
		}
	}

	steps = []multistep.Step{
		&stepCreateSSHKey{
			Debug:        b.config.PackerDebug,
			DebugKeyPath: fmt.Sprintf("sakuracloud_%s.pem", b.config.PackerBuildName),
		},
		&stepCreateServer{
			Debug: b.config.PackerDebug,
		},
		&stepBootWait{
			Debug: b.config.PackerDebug,
		},
		&stepServerInfo{
			Debug: b.config.PackerDebug,
		},
		&stepTypeBootCommand{
			Debug: b.config.PackerDebug,
			Ctx:   b.config.ctx,
		},
		communicateStep, // ssh or winrm
		new(commonsteps.StepProvision),
		&stepPowerOff{
			Debug: b.config.PackerDebug,
		},
		&stepCreateArchive{
			Debug: b.config.PackerDebug,
		},
		&stepTransferArchive{
			Debug: b.config.PackerDebug,
		},
	}

	if b.config.OSType == constants.TargetOSISO {
		steps = append(isoSteps, steps...)
	}

	// Run the steps
	if b.config.PackerDebug {
		b.runner = &multistep.DebugRunner{
			Steps:   steps,
			PauseFn: commonsteps.MultistepDebugFn(ui),
		}
	} else {
		b.runner = &multistep.BasicRunner{Steps: steps}
	}

	b.runner.Run(ctx, state)

	// If there was an error, return that
	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	if _, ok := state.GetOk("archive_id"); !ok {
		log.Println("Failed to find snapshot_name in state. Bug?")
		return nil, nil
	}

	artifact := &Artifact{
		archiveID:        state.Get("archive_id").(types.ID),
		archiveName:      state.Get("archive_name").(string),
		transferredIDs:   state.Get("transferred_ids").([]types.ID),
		transferredZones: state.Get("transferred_zones").([]string),
		client:           client.Archive,
	}

	return artifact, nil
}
