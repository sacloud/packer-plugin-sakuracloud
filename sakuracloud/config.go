package sakuracloud

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/communicator"
	"github.com/mitchellh/packer/helper/config"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template/interpolate"
	"github.com/sacloud/libsacloud/sacloud/ostype"
	"github.com/sacloud/packer-builder-sakuracloud/sakuracloud/constants"
	"os"
	"time"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	Comm                communicator.Config `mapstructure:",squash"`

	// for API Auth
	AccessToken       string `mapstructure:"access_token"`
	AccessTokenSecret string `mapstructure:"access_token_secret"`
	Zone              string `mapstructure:"zone"`

	// for Communication
	UserName      string `mapstructure:"user_name"`
	Password      string `mapstructure:"password"`
	UseUSKeyboard bool   `mapstructure:"us_keyboard"`

	// for Disk
	DiskSize       int    `mapstructure:"disk_size"`
	DiskConnection string `mapstructure:"disk_connection"`
	DiskPlan       string `mapstructure:"disk_plan"`

	// for Server
	Core                int  `mapstructure:"core"`
	MemorySize          int  `mapstructure:"memory_size"`
	DisableVirtIONetPCI bool `mapstructure:"disable_virtio_net"`

	// for Source
	OSType        string `mapstructure:"os_type"`
	SourceArchive int64  `mapstructure:"source_archive"`
	SourceDisk    int64  `mapstructure:"source_disk"`

	// for ISO
	ISOImageID       int64 `mapstructure:"iso_id"`
	common.ISOConfig `mapstructure:",squash"`
	ISOImageSizeGB   int    `mapstructure:"iso_size"`
	ISOImageName     string `mapstructure:"iso_name"`

	// for artifact
	ArchiveName        string        `mapstructure:"archive_name"`
	ArchiveTags        []string      `mapstructure:"archive_tags"`
	ArchiveDescription string        `mapstructure:"archive_description"`
	APIClientTimeout   time.Duration `mapstructure:"api_client_timeout"`

	BootWait    time.Duration `mapstructure:"boot_wait"`
	BootCommand []string      `mapstructure:"boot_command"`

	ctx interpolate.Context
}

func NewConfig(raws ...interface{}) (*Config, []string, error) {
	c := &Config{}

	var md mapstructure.Metadata
	err := config.Decode(c, &config.DecodeOpts{
		Metadata:           &md,
		Interpolate:        true,
		InterpolateContext: &c.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{
				"boot_command",
			},
		},
	}, raws...)
	if err != nil {
		return nil, nil, err
	}

	setDefaultConfig(c)

	var errs *packer.MultiError
	warnings := make([]string, 0)
	if es := c.Comm.Prepare(&c.ctx); len(es) > 0 {
		errs = packer.MultiErrorAppend(errs, es...)
	}

	if c.OSType == constants.TargetOSISO && c.ISOImageID == 0 {
		isoWarnings, isoErrs := c.ISOConfig.Prepare(&c.ctx)
		warnings = append(warnings, isoWarnings...)
		errs = packer.MultiErrorAppend(errs, isoErrs...)
	}

	// validate
	errs = validateConfig(c, errs)

	if errs != nil && len(errs.Errors) > 0 {
		return nil, warnings, errs
	}

	return c, warnings, nil
}

func setDefaultConfig(c *Config) {
	// Defaults
	if c.AccessToken == "" {
		// Default to environment variable for api_token, if it exists
		c.AccessToken = os.Getenv("SAKURACLOUD_ACCESS_TOKEN")
	}
	if c.AccessTokenSecret == "" {
		c.AccessTokenSecret = os.Getenv("SAKURACLOUD_ACCESS_TOKEN_SECRET")
	}

	if c.DiskConnection == "" {
		c.DiskConnection = "virtio"
	}
	if c.DiskPlan == "" {
		c.DiskPlan = "ssd"
	}
	if c.DiskSize == 0 {
		c.DiskSize = 20
	}
	if c.Core == 0 {
		c.Core = 1
	}
	if c.MemorySize == 0 {
		c.MemorySize = 1
	}

	if c.ArchiveName == "" {
		def, err := interpolate.Render("packer-{{timestamp}}", nil)
		if err != nil {
			panic(err)
		}

		// Default to packer-{{ unix timestamp (utc) }}
		c.ArchiveName = def
	}

	if c.Comm.SSHUsername == "" {
		// Default to "root". You can override this if your
		// SourceImage has a different user account then the SakuraCloud default
		c.Comm.SSHUsername = c.UserName
		if c.Comm.SSHUsername == "" {
			c.Comm.SSHUsername = "root"
		}
	}
	os := ostype.StrToOSType(c.OSType)
	if os.IsWindows() {
		if c.Comm.WinRMUser == "" {
			c.Comm.WinRMUser = c.UserName
		}
		if c.Comm.WinRMUser == "" {
			c.Comm.WinRMUser = "Administrator"
		}
		c.Comm.WinRMPassword = c.Password
		c.Comm.WinRMTimeout = 10 * time.Minute
		if c.Comm.WinRMPort == 0 {
			if c.Comm.WinRMUseSSL {
				c.Comm.WinRMPort = 5986
			} else {
				c.Comm.WinRMPort = 5985
			}
		}
	}

	if c.APIClientTimeout == 0 {
		// Default to 20 minute timeouts waiting for
		// desired state. i.e waiting for droplet to become active
		c.APIClientTimeout = 20 * time.Minute
	}

	if c.ISOImageSizeGB == 0 {
		c.ISOImageSizeGB = 5
	}
	if c.ISOImageName == "" {
		def, err := interpolate.Render("packer-{{timestamp}}", nil)
		if err != nil {
			panic(err)
		}

		// Default to packer-{{ unix timestamp (utc) }}
		if c.ISOChecksum == "" {
			c.ISOImageName = def
		} else {
			c.ISOImageName = c.ISOChecksum
		}
	}

	//if c.BootWait == 0 {
	//	//c.BootWait = 10 * time.Second
	//}
}

func validateConfig(c *Config, errs *packer.MultiError) *packer.MultiError {
	// required
	if c.AccessToken == "" {
		// Required configurations that will display errors if not set
		errs = packer.MultiErrorAppend(
			errs, errors.New("access_token is required"))
	}
	if c.AccessTokenSecret == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("access_token_secret is required"))
	}
	if c.Zone == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("zone is required"))
	}

	// os type
	if c.OSType == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("os_type is required"))
	}
	if !isInWord(c.OSType, listOSType()) {
		errs = packer.MultiErrorAppend(
			errs, errors.New("os_type is invalid"))
	}

	os := ostype.StrToOSType(c.OSType)
	if c.OSType != constants.TargetOSISO {
		if os == ostype.Custom {
			if c.SourceArchive == 0 && c.SourceDisk == 0 {
				errs = packer.MultiErrorAppend(
					errs, errors.New("source_archive or source_disk is required when os_type is custom"))
			}
		}
	}

	// Disk
	if !isInWord(c.DiskConnection, listDiskConnection()) {
		errs = packer.MultiErrorAppend(
			errs, errors.New("disk_connection is invalid"))
	}
	if !isInWord(c.DiskPlan, listDiskPlan()) {
		errs = packer.MultiErrorAppend(
			errs, errors.New("disk_plan is invalid"))
	}
	if c.DiskSize < 0 {
		errs = packer.MultiErrorAppend(
			errs, errors.New("size is invalid"))
	}

	// server
	if c.Core < 0 {
		errs = packer.MultiErrorAppend(
			errs, errors.New("core is invalid"))
	}
	if c.MemorySize < 0 {
		errs = packer.MultiErrorAppend(
			errs, errors.New("memory_size is invalid"))
	}

	if c.ISOImageSizeGB != 5 && c.ISOImageSizeGB != 10 {
		errs = packer.MultiErrorAppend(
			errs, errors.New("iso_size is require 5 or 10"))
	}

	return errs
}

func listOSType() []string {
	return append(ostype.OSTypeShortNames, constants.TargetOSISO)
}

func listDiskConnection() []string {
	return []string{
		"ide",
		"virtio",
	}
}

func listDiskPlan() []string {
	return []string{
		"ssd",
		"hdd",
	}
}

func isInWord(target string, list []string) bool {
	if len(list) == 0 || target == "" {
		return true
	}

	for _, t := range list {
		if t == target {
			return true
		}
	}
	return false
}
