package sakuracloud

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
	"github.com/sacloud/libsacloud/v2/helper/api"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/search"
)

func cleanupArchives() error {
	client := api.NewCaller(&api.CallerOptions{
		AccessToken:       os.Getenv("SAKURACLOUD_ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("SAKURACLOUD_ACCESS_TOKEN_SECRET"),
	})
	archiveOp := sacloud.NewArchiveOp(client)
	found, err := archiveOp.Find(context.Background(), "is1a", &sacloud.FindCondition{
		Filter: search.Filter{
			search.Key("Name"): search.PartialMatch("packer-acctest-"),
		},
	})
	if err != nil {
		return err
	}
	for _, archive := range found.Archives {
		if err := archiveOp.Delete(context.Background(), "is1a", archive.ID); err != nil {
			return err
		}
	}
	return nil
}

func TestBuilderAcc_basic(t *testing.T) {
	acctest.TestPlugin(t, &acctest.PluginTestCase{
		Name: "sakuracloud-minimum",
		CheckInit: func(cmd *exec.Cmd, s string) error {
			testAccPreCheck(t)
			return nil
		},
		Template: testBuilderAccBasic,
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}
			return nil
		},
		Teardown: cleanupArchives,
	})
}

func TestBuilderAcc_withSSHKeyFile(t *testing.T) {
	deferFunc := prepareTestPrivateKeyFile()
	acctest.TestPlugin(t, &acctest.PluginTestCase{
		Name: "sakuracloud-with-ssh-key",
		CheckInit: func(cmd *exec.Cmd, s string) error {
			testAccPreCheck(t)
			return nil
		},
		Template: testBuilderAccWithSSHPrivateKeyFile(dummyPrivateKeyFile),
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}
			return nil
		},
		Teardown: func() error {
			deferFunc()
			return cleanupArchives()
		},
	})
}

func testAccPreCheck(t *testing.T) {
	requiredEnvs := []string{"SAKURACLOUD_ACCESS_TOKEN", "SAKURACLOUD_ACCESS_TOKEN_SECRET"}

	for _, k := range requiredEnvs {
		if v := os.Getenv(k); v == "" {
			t.Fatalf("%q must be set for acceptance tests", k)
		}
	}
}

const testBuilderAccBasic = `
{
    "builders": [{
        "type": "sakuracloud",
        "zone": "is1a",
        "os_type": "ubuntu2004",
        "user_name": "ubuntu",
        "password": "TestUserPassword01",
        "disk_size": 20,
        "disk_plan": "ssd",
        "core" : 2,
        "memory_size": 4,
        "archive_name": "packer-acctest-minimum",
        "archive_description": "description of archive"
    }]
}
`

func testBuilderAccWithSSHPrivateKeyFile(keyPath string) string {
	return `
{
   "builders": [{
       "type": "sakuracloud",
       "zone": "is1a",
       "os_type": "ubuntu2004",
       "user_name": "ubuntu",
       "password": "TestUserPassword01",
       "disk_size": 20,
       "disk_plan": "ssd",
       "core" : 2,
       "memory_size": 4,
       "ssh_private_key_file": "` + keyPath + `",
       "archive_name": "packer-acctest-sshkey",
       "archive_description": "description of archive"
   }]
}
`
}
