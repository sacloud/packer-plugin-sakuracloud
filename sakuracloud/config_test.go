package sakuracloud

import (
	"os"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/stretchr/testify/assert"
)

func TestConfig_Defaults(t *testing.T) {
	os.Setenv("SAKURA_ACCESS_TOKEN", "token")         //nolint
	os.Setenv("SAKURA_ACCESS_TOKEN_SECRET", "secret") //nolint

	config := &Config{}
	config.OSType = "windows2016"
	setDefaultConfig(config)

	expects := []struct {
		expect interface{}
		actual interface{}
	}{
		{"token", config.AccessToken},
		{"secret", config.AccessTokenSecret},
		{defaultDiskConnection, config.DiskConnection},
		{defaultDiskPlan, config.DiskPlan},
		{defaultDiskSize, config.DiskSize},
		{defaultCore, config.Core},
		{defaultMemory, config.MemorySize},
		{defaultAPITimeout, config.APIClientTimeout},
		{defaultSSHUser, config.Comm.SSHUsername},
		{defaultWinRMUser, config.Comm.WinRMUser},
		{defaultWinRMTimeout, config.Comm.WinRMTimeout},
		{5985, config.Comm.WinRMPort},
		{defaultISOIMageSize, config.ISOImageSizeGB},
	}

	for _, v := range expects {
		assert.Equal(t, v.expect, v.actual)
	}
}

func TestConfig_Validate(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		config := dummyValidConfig()
		setDefaultConfig(&config)

		var errs = &packer.MultiError{}
		errs = validateConfig(&config, errs)

		assert.True(t, len(errs.Errors) == 0)
	})

	t.Run("valid config", func(t *testing.T) {
		clearEnvKeys := []string{
			"SAKURA_ACCESS_TOKEN",
			"SAKURA_ACCESS_TOKEN_SECRET",
			"SAKURA_ZONE",
		}
		for _, key := range clearEnvKeys {
			os.Setenv(key, "") //nolint:errcheck,gosec
		}

		expects := []struct {
			caseName string
			filter   func(*Config)
			expect   bool
		}{
			{
				caseName: "access_token",
				filter:   func(c *Config) { c.AccessToken = "" },
			},
			{
				caseName: "access_token_secret",
				filter:   func(c *Config) { c.AccessTokenSecret = "" },
			},
			{
				caseName: "zone",
				filter:   func(c *Config) { c.Zone = "" },
			},
			{
				caseName: "zone:custom",
				filter:   func(c *Config) { c.Zone = "custom" },
				expect:   true,
			},
			{
				caseName: "os_type",
				filter:   func(c *Config) { c.OSType = "" },
			},
			{
				caseName: "os_type:ubuntu",
				filter:   func(c *Config) { c.OSType = "ubuntu" },
				expect:   true,
			},
			{
				caseName: "os_type:invalid",
				filter:   func(c *Config) { c.OSType = "invalid" },
			},
			{
				caseName: "os_type:custom",
				filter:   func(c *Config) { c.OSType = "custom" },
			},
			{
				caseName: "os_type:custom+disk",
				filter: func(c *Config) {
					c.OSType = "custom"
					c.SourceDisk = dummyDiskID
				},
				expect: true,
			},
			{
				caseName: "os_type:custom+archive",
				filter: func(c *Config) {
					c.OSType = "custom"
					c.SourceArchive = dummyArchiveID
				},
				expect: true,
			},
			{
				caseName: "os_type:iso",
				filter:   func(c *Config) { c.OSType = "iso" },
				expect:   true,
			},
			{
				caseName: "disk_connection:invalid",
				filter:   func(c *Config) { c.DiskConnection = "invalid" },
			},
			{
				caseName: "disk_plan:invalid",
				filter:   func(c *Config) { c.DiskPlan = "invalid" },
			},
			{
				caseName: "disk_size:invalid",
				filter:   func(c *Config) { c.DiskSize = -1 },
			},
			{
				caseName: "core:invalid",
				filter:   func(c *Config) { c.Core = -1 },
			},
			{
				caseName: "memory_size:invalid",
				filter:   func(c *Config) { c.MemorySize = -1 },
			},
			{
				caseName: "iso_image_size:invalid",
				filter:   func(c *Config) { c.ISOImageSizeGB = 1 },
			},
			{
				caseName: "iso_image_size:5g",
				filter:   func(c *Config) { c.ISOImageSizeGB = 5 },
				expect:   true,
			},
			{
				caseName: "iso_image_size:10g",
				filter:   func(c *Config) { c.ISOImageSizeGB = 10 },
				expect:   true,
			},
		}

		for _, expect := range expects {
			config := dummyValidConfig()
			expect.filter(&config)
			setDefaultConfig(&config)

			var errs = &packer.MultiError{}
			errs = validateConfig(&config, errs)
			assert.Equal(t, expect.expect, len(errs.Errors) == 0,
				"unexpected validate result: expect=%t target=%q errs=%s", expect.expect, expect.caseName, errs.Error())
		}
	})
}

func dummyValidConfig() Config {
	return Config{
		AccessToken:       "token",
		AccessTokenSecret: "secret",
		Zone:              "is1a",
		OSType:            "ubuntu",
	}
}
