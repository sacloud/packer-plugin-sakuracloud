package sakuracloud

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
	sacloudClient "github.com/sacloud/api-client-go"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/api"
	"github.com/sacloud/iaas-api-go/search"
)

func cleanupArchives() error {
	zone := os.Getenv("SAKURACLOUD_ZONE")
	client := api.NewCallerWithOptions(&api.CallerOptions{
		Options: &sacloudClient.Options{
			AccessToken:       os.Getenv("SAKURACLOUD_ACCESS_TOKEN"),
			AccessTokenSecret: os.Getenv("SAKURACLOUD_ACCESS_TOKEN_SECRET"),
		},
	})
	archiveOp := iaas.NewArchiveOp(client)
	found, err := archiveOp.Find(context.Background(), zone, &iaas.FindCondition{
		Filter: search.Filter{
			search.Key("Name"): search.PartialMatch("packer-acctest-"),
		},
	})
	if err != nil {
		return err
	}
	for _, archive := range found.Archives {
		if err := archiveOp.Delete(context.Background(), zone, archive.ID); err != nil {
			return err
		}
	}
	return nil
}

func TestBuilderAcc_withSSHKeyFile(t *testing.T) {
	zone := os.Getenv("SAKURACLOUD_ZONE")
	deferFunc := prepareTestPrivateKeyFile()
	acctest.TestPlugin(t, &acctest.PluginTestCase{
		Name: "sakuracloud-with-ssh-key",
		CheckInit: func(cmd *exec.Cmd, s string) error {
			testAccPreCheck(t)
			return nil
		},
		Template: testBuilderAccWithSSHPrivateKeyFile(zone, dummyPrivateKeyFile),
		Check:    testAccCheckFunc,
		Teardown: func() error {
			deferFunc()
			return cleanupArchives()
		},
	})
}

func testAccCheckFunc(buildCommand *exec.Cmd, logfile string) error {
	logs, err := os.Open(logfile) //nolint:gosec
	if err != nil {
		return fmt.Errorf("Unable find %s", logfile)
	}
	defer logs.Close() //nolint:errcheck

	logsBytes, err := io.ReadAll(logs)
	if err != nil {
		return fmt.Errorf("Unable to read %s", logfile)
	}

	if buildCommand.ProcessState != nil {
		if buildCommand.ProcessState.ExitCode() != 0 {
			log.Println(string(logsBytes))
			return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
		}
	}
	return nil
}

func testAccPreCheck(t *testing.T) {
	requiredEnvs := []string{"SAKURACLOUD_ACCESS_TOKEN", "SAKURACLOUD_ACCESS_TOKEN_SECRET", "SAKURACLOUD_ZONE"}

	for _, k := range requiredEnvs {
		if v := os.Getenv(k); v == "" {
			t.Fatalf("%q must be set for acceptance tests", k)
		}
	}
}

//go:embed test-fixtures/with-ssh-key.pkr.hcl
var testBuilderHCL2WithSshKey string

func testBuilderAccWithSSHPrivateKeyFile(zone, keyPath string) string {
	return fmt.Sprintf(testBuilderHCL2WithSshKey, zone, keyPath)
}
