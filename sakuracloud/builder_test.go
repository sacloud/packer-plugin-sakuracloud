package sakuracloud

import (
	"os"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_ImplementsBuilder(t *testing.T) {
	var raw interface{} = &Builder{}
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
		builder := &Builder{}
		_, warns, errs := builder.Prepare(testMinimumConfigValues)

		assert.Nil(t, warns)
		assert.Nil(t, errs)
	})

	// TODO add more unit tests after refactoring Builder/Config
}
