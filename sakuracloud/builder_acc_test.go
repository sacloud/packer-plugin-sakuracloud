package sakuracloud

import (
	"os"
	"testing"

	builderT "github.com/hashicorp/packer/helper/builder/testing"
)

func TestBuilderAcc_basic(t *testing.T) {
	builderT.Test(t, builderT.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Builder:  &Builder{},
		Template: testBuilderAccBasic,
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
