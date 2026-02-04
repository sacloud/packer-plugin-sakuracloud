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
		"SAKURA_ACCESS_TOKEN",
		"SAKURA_ACCESS_TOKEN_SECRET",
		"SAKURA_ZONE",
	}
	for _, key := range clearEnvKeys {
		os.Setenv(key, "") //nolint:errcheck,gosec
	}

	t.Run("with minimum config", func(t *testing.T) {
		builder := &Builder{}
		_, warns, errs := builder.Prepare(testMinimumConfigValues)

		assert.Nil(t, warns)
		assert.Nil(t, errs)
	})

	// TODO add more unit tests after refactoring Builder/Config
}
