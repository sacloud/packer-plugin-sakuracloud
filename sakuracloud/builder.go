package sakuracloud

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/communicator"
	"github.com/mitchellh/packer/packer"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/packer-builder-sakuracloud/sakuracloud/constants"
	"log"
)

// The unique id for the builder
const BuilderId = "packer.sakuracloud"

type Builder struct {
	config Config
	runner multistep.Runner
}

func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {
	c, warnings, errs := NewConfig(raws...)
	if errs != nil {
		return warnings, errs
	}
	b.config = *c

	return nil, nil
}

func (b *Builder) Cancel() {
	if b.runner != nil {
		log.Println("Cancelling the step runner...")
		b.runner.Cancel()
	}
}

func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	client := api.NewClient(b.config.AccessToken, b.config.AccessTokenSecret, b.config.Zone)

	// Set up the state
	state := new(multistep.BasicStateBag)
	state.Put("config", b.config)
	state.Put("client", client)
	state.Put("hook", hook)
	state.Put("ui", ui)
	state.Put("cache", cache)

	// Build the steps
	var steps []multistep.Step

	var communicateStep multistep.Step
	if b.config.OSType == constants.TargetOSWindows {
		communicateStep = &communicator.StepConnectWinRM{
			Config: &b.config.Comm,
			Host:   commHost,
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
			Host:      commHost,
			SSHConfig: sshConfig,
		}
	}

	var isoSteps []multistep.Step
	if b.config.ISOImageID > 0 {
		isoSteps = []multistep.Step{
			&stepPrepareISO{},
		}
	} else {
		isoSteps = []multistep.Step{
			&common.StepDownload{
				Checksum:     b.config.ISOChecksum,
				ChecksumType: b.config.ISOChecksumType,
				Description:  "ISO",
				Extension:    "iso",
				ResultKey:    "iso_path",
				TargetPath:   b.config.TargetPath,
				Url:          b.config.ISOUrls,
			},
			&stepRemoteUpload{},
			&stepPrepareISO{},
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
		new(stepBootWait),
		&stepServerInfo{
			Debug: b.config.PackerDebug,
		},
		&stepTypeBootCommand{
			Ctx: b.config.ctx,
		},
		communicateStep, // ssh or winrm
		new(common.StepProvision),
		&stepShutdown{
			Debug: b.config.PackerDebug,
		},
		&stepPowerOff{
			Debug: b.config.PackerDebug,
		},
		&stepCreateArchive{
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
			PauseFn: common.MultistepDebugFn(ui),
		}
	} else {
		b.runner = &multistep.BasicRunner{Steps: steps}
	}

	b.runner.Run(state)

	// If there was an error, return that
	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	if _, ok := state.GetOk("archive_id"); !ok {
		log.Println("Failed to find snapshot_name in state. Bug?")
		return nil, nil
	}

	artifact := &Artifact{
		archiveID:   state.Get("archive_id").(int64),
		archiveName: state.Get("archive_name").(string),
		zone:        state.Get("zone").(string),
		client:      client,
	}

	return artifact, nil
}
