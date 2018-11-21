package sakuracloud

import (
	"os"
	"testing"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_ImplementsBuilder(t *testing.T) {
	var raw interface{}
	raw = &Builder{}
	if _, ok := raw.(packer.Builder); !ok {
		t.Fatalf("Builder should be a builder")
	}
}

func TestBuilder_Prepare(t *testing.T) {

	clearEnvKeys := []string{
		"SAKURACLOUD_ACCESS_TOKEN",
		"SAKURACLOUD_ACCESS_TOKEN_SECRET",
		"SAKURACLOUD_ZONE",
	}
	for _, key := range clearEnvKeys {
		os.Setenv(key, "")
	}

	t.Run("with minimum config", func(t *testing.T) {
		conf := map[string]interface{}{
			"access_token":        "aaaa",
			"access_token_secret": "bbbb",
			"zone":                "is1a",
			"os_type":             "centos",
		}
		builder := &Builder{}
		warns, errs := builder.Prepare(conf)

		assert.Nil(t, warns)
		assert.Nil(t, errs)
	})

	// TODO add more unit tests after refactoring Builder/Config
}

type testBuildRunner struct {
	cancelCalled bool
}

func (t *testBuildRunner) Run(bag multistep.StateBag) {}
func (t *testBuildRunner) Cancel()                    { t.cancelCalled = true }

func TestBuilder_Cancel(t *testing.T) {
	runner := &testBuildRunner{}
	builder := &Builder{
		runner: runner,
	}

	builder.Cancel()
	assert.True(t, runner.cancelCalled)
}
