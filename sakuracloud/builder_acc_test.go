package sakuracloud

import (
	"os"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
)

func TestBuilderAcc_basic(t *testing.T) {
	acctest.Test(t, acctest.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Builder:  &Builder{},
		Template: testBuilderAccBasic,
	})
}

func TestBuilderAcc_withSSHKeyFile(t *testing.T) {
	var deferFunc func()
	acctest.Test(t, acctest.TestCase{
		PreCheck: func() {
			deferFunc = prepareTestPrivateKeyFile()
			testAccPreCheck(t)
		},
		Builder:  &Builder{},
		Template: testBuilderAccWithSSHPrivateKeyFile(dummyPrivateKeyFile),
		Teardown: func() error {
			deferFunc()
			return nil
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
        "type": "test",
        "zone": "is1a",
        "os_type": "centos",
        "password": "TestUserPassword01",
        "disk_size": 20,
        "disk_plan": "ssd",
        "core" : 2,
        "memory_size": 4,
        "archive_name": "packer-example-centos",
        "archive_description": "description of archive"
    }]
}
`

func testBuilderAccWithSSHPrivateKeyFile(keyPath string) string {
	return `
{
    "builders": [{
        "type": "test",
        "zone": "is1a",
        "os_type": "centos",
        "password": "TestUserPassword01",
        "disk_size": 20,
        "disk_plan": "ssd",
        "core" : 2,
        "memory_size": 4,
		"ssh_private_key_file": "` + keyPath + `",
        "archive_name": "packer-example-centos",
        "archive_description": "description of archive"
    }]
}
`
}
