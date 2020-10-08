// Code generated by "mapstructure-to-hcl2 -type Config"; DO NOT EDIT.
package sakuracloud

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/zclconf/go-cty/cty"
)

// FlatConfig is an auto-generated flat version of Config.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatConfig struct {
	PackerBuildName           *string           `mapstructure:"packer_build_name" cty:"packer_build_name"`
	PackerBuilderType         *string           `mapstructure:"packer_builder_type" cty:"packer_builder_type"`
	PackerDebug               *bool             `mapstructure:"packer_debug" cty:"packer_debug"`
	PackerForce               *bool             `mapstructure:"packer_force" cty:"packer_force"`
	PackerOnError             *string           `mapstructure:"packer_on_error" cty:"packer_on_error"`
	PackerUserVars            map[string]string `mapstructure:"packer_user_variables" cty:"packer_user_variables"`
	PackerSensitiveVars       []string          `mapstructure:"packer_sensitive_variables" cty:"packer_sensitive_variables"`
	Type                      *string           `mapstructure:"communicator" cty:"communicator"`
	PauseBeforeConnect        *string           `mapstructure:"pause_before_connecting" cty:"pause_before_connecting"`
	SSHHost                   *string           `mapstructure:"ssh_host" cty:"ssh_host"`
	SSHPort                   *int              `mapstructure:"ssh_port" cty:"ssh_port"`
	SSHUsername               *string           `mapstructure:"ssh_username" cty:"ssh_username"`
	SSHPassword               *string           `mapstructure:"ssh_password" cty:"ssh_password"`
	SSHKeyPairName            *string           `mapstructure:"ssh_keypair_name" cty:"ssh_keypair_name"`
	SSHTemporaryKeyPairName   *string           `mapstructure:"temporary_key_pair_name" cty:"temporary_key_pair_name"`
	SSHClearAuthorizedKeys    *bool             `mapstructure:"ssh_clear_authorized_keys" cty:"ssh_clear_authorized_keys"`
	SSHPrivateKeyFile         *string           `mapstructure:"ssh_private_key_file" cty:"ssh_private_key_file"`
	SSHPty                    *bool             `mapstructure:"ssh_pty" cty:"ssh_pty"`
	SSHTimeout                *string           `mapstructure:"ssh_timeout" cty:"ssh_timeout"`
	SSHAgentAuth              *bool             `mapstructure:"ssh_agent_auth" cty:"ssh_agent_auth"`
	SSHDisableAgentForwarding *bool             `mapstructure:"ssh_disable_agent_forwarding" cty:"ssh_disable_agent_forwarding"`
	SSHHandshakeAttempts      *int              `mapstructure:"ssh_handshake_attempts" cty:"ssh_handshake_attempts"`
	SSHBastionHost            *string           `mapstructure:"ssh_bastion_host" cty:"ssh_bastion_host"`
	SSHBastionPort            *int              `mapstructure:"ssh_bastion_port" cty:"ssh_bastion_port"`
	SSHBastionAgentAuth       *bool             `mapstructure:"ssh_bastion_agent_auth" cty:"ssh_bastion_agent_auth"`
	SSHBastionUsername        *string           `mapstructure:"ssh_bastion_username" cty:"ssh_bastion_username"`
	SSHBastionPassword        *string           `mapstructure:"ssh_bastion_password" cty:"ssh_bastion_password"`
	SSHBastionPrivateKeyFile  *string           `mapstructure:"ssh_bastion_private_key_file" cty:"ssh_bastion_private_key_file"`
	SSHFileTransferMethod     *string           `mapstructure:"ssh_file_transfer_method" cty:"ssh_file_transfer_method"`
	SSHProxyHost              *string           `mapstructure:"ssh_proxy_host" cty:"ssh_proxy_host"`
	SSHProxyPort              *int              `mapstructure:"ssh_proxy_port" cty:"ssh_proxy_port"`
	SSHProxyUsername          *string           `mapstructure:"ssh_proxy_username" cty:"ssh_proxy_username"`
	SSHProxyPassword          *string           `mapstructure:"ssh_proxy_password" cty:"ssh_proxy_password"`
	SSHKeepAliveInterval      *string           `mapstructure:"ssh_keep_alive_interval" cty:"ssh_keep_alive_interval"`
	SSHReadWriteTimeout       *string           `mapstructure:"ssh_read_write_timeout" cty:"ssh_read_write_timeout"`
	SSHRemoteTunnels          []string          `mapstructure:"ssh_remote_tunnels" cty:"ssh_remote_tunnels"`
	SSHLocalTunnels           []string          `mapstructure:"ssh_local_tunnels" cty:"ssh_local_tunnels"`
	SSHPublicKey              []byte            `mapstructure:"ssh_public_key" cty:"ssh_public_key"`
	SSHPrivateKey             []byte            `mapstructure:"ssh_private_key" cty:"ssh_private_key"`
	WinRMUser                 *string           `mapstructure:"winrm_username" cty:"winrm_username"`
	WinRMPassword             *string           `mapstructure:"winrm_password" cty:"winrm_password"`
	WinRMHost                 *string           `mapstructure:"winrm_host" cty:"winrm_host"`
	WinRMPort                 *int              `mapstructure:"winrm_port" cty:"winrm_port"`
	WinRMTimeout              *string           `mapstructure:"winrm_timeout" cty:"winrm_timeout"`
	WinRMUseSSL               *bool             `mapstructure:"winrm_use_ssl" cty:"winrm_use_ssl"`
	WinRMInsecure             *bool             `mapstructure:"winrm_insecure" cty:"winrm_insecure"`
	WinRMUseNTLM              *bool             `mapstructure:"winrm_use_ntlm" cty:"winrm_use_ntlm"`
	AccessToken               *string           `mapstructure:"access_token" cty:"access_token"`
	AccessTokenSecret         *string           `mapstructure:"access_token_secret" cty:"access_token_secret"`
	Zone                      *string           `mapstructure:"zone" cty:"zone"`
	Zones                     []string          `mapstructure:"zones" cty:"zones"`
	UserName                  *string           `mapstructure:"user_name" cty:"user_name"`
	Password                  *string           `mapstructure:"password" cty:"password"`
	UseUSKeyboard             *bool             `mapstructure:"us_keyboard" cty:"us_keyboard"`
	DiskSize                  *int              `mapstructure:"disk_size" cty:"disk_size"`
	DiskConnection            *string           `mapstructure:"disk_connection" cty:"disk_connection"`
	DiskPlan                  *string           `mapstructure:"disk_plan" cty:"disk_plan"`
	Core                      *int              `mapstructure:"core" cty:"core"`
	MemorySize                *int              `mapstructure:"memory_size" cty:"memory_size"`
	DisableVirtIONetPCI       *bool             `mapstructure:"disable_virtio_net" cty:"disable_virtio_net"`
	ForceShutdown             *bool             `mapstructure:"force_shutdown" cty:"force_shutdown"`
	OSType                    *string           `mapstructure:"os_type" cty:"os_type"`
	SourceArchive             *types.ID         `mapstructure:"source_archive" cty:"source_archive"`
	SourceDisk                *types.ID         `mapstructure:"source_disk" cty:"source_disk"`
	ISOImageID                *types.ID         `mapstructure:"iso_id" cty:"iso_id"`
	ISOChecksum               *string           `mapstructure:"iso_checksum" required:"true" cty:"iso_checksum"`
	ISOChecksumURL            *string           `mapstructure:"iso_checksum_url" cty:"iso_checksum_url"`
	ISOChecksumType           *string           `mapstructure:"iso_checksum_type" cty:"iso_checksum_type"`
	RawSingleISOUrl           *string           `mapstructure:"iso_url" required:"true" cty:"iso_url"`
	ISOUrls                   []string          `mapstructure:"iso_urls" cty:"iso_urls"`
	TargetPath                *string           `mapstructure:"iso_target_path" cty:"iso_target_path"`
	TargetExtension           *string           `mapstructure:"iso_target_extension" cty:"iso_target_extension"`
	ISOImageSizeGB            *int              `mapstructure:"iso_size" cty:"iso_size"`
	ISOImageName              *string           `mapstructure:"iso_name" cty:"iso_name"`
	ArchiveName               *string           `mapstructure:"archive_name" cty:"archive_name"`
	ArchiveTags               []string          `mapstructure:"archive_tags" cty:"archive_tags"`
	ArchiveDescription        *string           `mapstructure:"archive_description" cty:"archive_description"`
	APIClientTimeout          *string           `mapstructure:"api_client_timeout" cty:"api_client_timeout"`
	DisableGeneratePublicKey  *bool             `mapstructure:"disable_generate_public_key" cty:"disable_generate_public_key"`
	BootWait                  *string           `mapstructure:"boot_wait" cty:"boot_wait"`
	BootCommand               []string          `mapstructure:"boot_command" cty:"boot_command"`
}

// FlatMapstructure returns a new FlatConfig.
// FlatConfig is an auto-generated flat version of Config.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*Config) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatConfig)
}

// HCL2Spec returns the hcl spec of a Config.
// This spec is used by HCL to read the fields of Config.
// The decoded values from this spec will then be applied to a FlatConfig.
func (*FlatConfig) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"packer_build_name":            &hcldec.AttrSpec{Name: "packer_build_name", Type: cty.String, Required: false},
		"packer_builder_type":          &hcldec.AttrSpec{Name: "packer_builder_type", Type: cty.String, Required: false},
		"packer_debug":                 &hcldec.AttrSpec{Name: "packer_debug", Type: cty.Bool, Required: false},
		"packer_force":                 &hcldec.AttrSpec{Name: "packer_force", Type: cty.Bool, Required: false},
		"packer_on_error":              &hcldec.AttrSpec{Name: "packer_on_error", Type: cty.String, Required: false},
		"packer_user_variables":        &hcldec.BlockAttrsSpec{TypeName: "packer_user_variables", ElementType: cty.String, Required: false},
		"packer_sensitive_variables":   &hcldec.AttrSpec{Name: "packer_sensitive_variables", Type: cty.List(cty.String), Required: false},
		"communicator":                 &hcldec.AttrSpec{Name: "communicator", Type: cty.String, Required: false},
		"pause_before_connecting":      &hcldec.AttrSpec{Name: "pause_before_connecting", Type: cty.String, Required: false},
		"ssh_host":                     &hcldec.AttrSpec{Name: "ssh_host", Type: cty.String, Required: false},
		"ssh_port":                     &hcldec.AttrSpec{Name: "ssh_port", Type: cty.Number, Required: false},
		"ssh_username":                 &hcldec.AttrSpec{Name: "ssh_username", Type: cty.String, Required: false},
		"ssh_password":                 &hcldec.AttrSpec{Name: "ssh_password", Type: cty.String, Required: false},
		"ssh_keypair_name":             &hcldec.AttrSpec{Name: "ssh_keypair_name", Type: cty.String, Required: false},
		"temporary_key_pair_name":      &hcldec.AttrSpec{Name: "temporary_key_pair_name", Type: cty.String, Required: false},
		"ssh_clear_authorized_keys":    &hcldec.AttrSpec{Name: "ssh_clear_authorized_keys", Type: cty.Bool, Required: false},
		"ssh_private_key_file":         &hcldec.AttrSpec{Name: "ssh_private_key_file", Type: cty.String, Required: false},
		"ssh_pty":                      &hcldec.AttrSpec{Name: "ssh_pty", Type: cty.Bool, Required: false},
		"ssh_timeout":                  &hcldec.AttrSpec{Name: "ssh_timeout", Type: cty.String, Required: false},
		"ssh_agent_auth":               &hcldec.AttrSpec{Name: "ssh_agent_auth", Type: cty.Bool, Required: false},
		"ssh_disable_agent_forwarding": &hcldec.AttrSpec{Name: "ssh_disable_agent_forwarding", Type: cty.Bool, Required: false},
		"ssh_handshake_attempts":       &hcldec.AttrSpec{Name: "ssh_handshake_attempts", Type: cty.Number, Required: false},
		"ssh_bastion_host":             &hcldec.AttrSpec{Name: "ssh_bastion_host", Type: cty.String, Required: false},
		"ssh_bastion_port":             &hcldec.AttrSpec{Name: "ssh_bastion_port", Type: cty.Number, Required: false},
		"ssh_bastion_agent_auth":       &hcldec.AttrSpec{Name: "ssh_bastion_agent_auth", Type: cty.Bool, Required: false},
		"ssh_bastion_username":         &hcldec.AttrSpec{Name: "ssh_bastion_username", Type: cty.String, Required: false},
		"ssh_bastion_password":         &hcldec.AttrSpec{Name: "ssh_bastion_password", Type: cty.String, Required: false},
		"ssh_bastion_private_key_file": &hcldec.AttrSpec{Name: "ssh_bastion_private_key_file", Type: cty.String, Required: false},
		"ssh_file_transfer_method":     &hcldec.AttrSpec{Name: "ssh_file_transfer_method", Type: cty.String, Required: false},
		"ssh_proxy_host":               &hcldec.AttrSpec{Name: "ssh_proxy_host", Type: cty.String, Required: false},
		"ssh_proxy_port":               &hcldec.AttrSpec{Name: "ssh_proxy_port", Type: cty.Number, Required: false},
		"ssh_proxy_username":           &hcldec.AttrSpec{Name: "ssh_proxy_username", Type: cty.String, Required: false},
		"ssh_proxy_password":           &hcldec.AttrSpec{Name: "ssh_proxy_password", Type: cty.String, Required: false},
		"ssh_keep_alive_interval":      &hcldec.AttrSpec{Name: "ssh_keep_alive_interval", Type: cty.String, Required: false},
		"ssh_read_write_timeout":       &hcldec.AttrSpec{Name: "ssh_read_write_timeout", Type: cty.String, Required: false},
		"ssh_remote_tunnels":           &hcldec.AttrSpec{Name: "ssh_remote_tunnels", Type: cty.List(cty.String), Required: false},
		"ssh_local_tunnels":            &hcldec.AttrSpec{Name: "ssh_local_tunnels", Type: cty.List(cty.String), Required: false},
		"ssh_public_key":               &hcldec.AttrSpec{Name: "ssh_public_key", Type: cty.List(cty.Number), Required: false},
		"ssh_private_key":              &hcldec.AttrSpec{Name: "ssh_private_key", Type: cty.List(cty.Number), Required: false},
		"winrm_username":               &hcldec.AttrSpec{Name: "winrm_username", Type: cty.String, Required: false},
		"winrm_password":               &hcldec.AttrSpec{Name: "winrm_password", Type: cty.String, Required: false},
		"winrm_host":                   &hcldec.AttrSpec{Name: "winrm_host", Type: cty.String, Required: false},
		"winrm_port":                   &hcldec.AttrSpec{Name: "winrm_port", Type: cty.Number, Required: false},
		"winrm_timeout":                &hcldec.AttrSpec{Name: "winrm_timeout", Type: cty.String, Required: false},
		"winrm_use_ssl":                &hcldec.AttrSpec{Name: "winrm_use_ssl", Type: cty.Bool, Required: false},
		"winrm_insecure":               &hcldec.AttrSpec{Name: "winrm_insecure", Type: cty.Bool, Required: false},
		"winrm_use_ntlm":               &hcldec.AttrSpec{Name: "winrm_use_ntlm", Type: cty.Bool, Required: false},
		"access_token":                 &hcldec.AttrSpec{Name: "access_token", Type: cty.String, Required: false},
		"access_token_secret":          &hcldec.AttrSpec{Name: "access_token_secret", Type: cty.String, Required: false},
		"zone":                         &hcldec.AttrSpec{Name: "zone", Type: cty.String, Required: false},
		"zones":                        &hcldec.AttrSpec{Name: "zones", Type: cty.List(cty.String), Required: false},
		"user_name":                    &hcldec.AttrSpec{Name: "user_name", Type: cty.String, Required: false},
		"password":                     &hcldec.AttrSpec{Name: "password", Type: cty.String, Required: false},
		"us_keyboard":                  &hcldec.AttrSpec{Name: "us_keyboard", Type: cty.Bool, Required: false},
		"disk_size":                    &hcldec.AttrSpec{Name: "disk_size", Type: cty.Number, Required: false},
		"disk_connection":              &hcldec.AttrSpec{Name: "disk_connection", Type: cty.String, Required: false},
		"disk_plan":                    &hcldec.AttrSpec{Name: "disk_plan", Type: cty.String, Required: false},
		"core":                         &hcldec.AttrSpec{Name: "core", Type: cty.Number, Required: false},
		"memory_size":                  &hcldec.AttrSpec{Name: "memory_size", Type: cty.Number, Required: false},
		"disable_virtio_net":           &hcldec.AttrSpec{Name: "disable_virtio_net", Type: cty.Bool, Required: false},
		"force_shutdown":               &hcldec.AttrSpec{Name: "force_shutdown", Type: cty.Bool, Required: false},
		"os_type":                      &hcldec.AttrSpec{Name: "os_type", Type: cty.String, Required: false},
		"source_archive":               &hcldec.AttrSpec{Name: "source_archive", Type: cty.Number, Required: false},
		"source_disk":                  &hcldec.AttrSpec{Name: "source_disk", Type: cty.Number, Required: false},
		"iso_id":                       &hcldec.AttrSpec{Name: "iso_id", Type: cty.Number, Required: false},
		"iso_checksum":                 &hcldec.AttrSpec{Name: "iso_checksum", Type: cty.String, Required: false},
		"iso_checksum_url":             &hcldec.AttrSpec{Name: "iso_checksum_url", Type: cty.String, Required: false},
		"iso_checksum_type":            &hcldec.AttrSpec{Name: "iso_checksum_type", Type: cty.String, Required: false},
		"iso_url":                      &hcldec.AttrSpec{Name: "iso_url", Type: cty.String, Required: false},
		"iso_urls":                     &hcldec.AttrSpec{Name: "iso_urls", Type: cty.List(cty.String), Required: false},
		"iso_target_path":              &hcldec.AttrSpec{Name: "iso_target_path", Type: cty.String, Required: false},
		"iso_target_extension":         &hcldec.AttrSpec{Name: "iso_target_extension", Type: cty.String, Required: false},
		"iso_size":                     &hcldec.AttrSpec{Name: "iso_size", Type: cty.Number, Required: false},
		"iso_name":                     &hcldec.AttrSpec{Name: "iso_name", Type: cty.String, Required: false},
		"archive_name":                 &hcldec.AttrSpec{Name: "archive_name", Type: cty.String, Required: false},
		"archive_tags":                 &hcldec.AttrSpec{Name: "archive_tags", Type: cty.List(cty.String), Required: false},
		"archive_description":          &hcldec.AttrSpec{Name: "archive_description", Type: cty.String, Required: false},
		"api_client_timeout":           &hcldec.AttrSpec{Name: "api_client_timeout", Type: cty.String, Required: false},
		"disable_generate_public_key":  &hcldec.AttrSpec{Name: "disable_generate_public_key", Type: cty.Bool, Required: false},
		"boot_wait":                    &hcldec.AttrSpec{Name: "boot_wait", Type: cty.String, Required: false},
		"boot_command":                 &hcldec.AttrSpec{Name: "boot_command", Type: cty.List(cty.String), Required: false},
	}
	return s
}