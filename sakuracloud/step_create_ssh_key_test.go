package sakuracloud

import (
	"context"
	"testing"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

func TestStepCreateSSHKey(t *testing.T) {
	t.Run("without private key config", func(t *testing.T) {
		ctx := context.Background()
		step := &stepCreateSSHKey{}

		state := dummyMinimumStateBag(nil)

		// run
		action := step.Run(ctx, state)

		privateKey := state.Get("privateKey").(string)
		publicKeys := state.Get("publicKeys").([]string)
		publicKey := publicKeys[0]

		assert.Equal(t, multistep.ActionContinue, action)
		assert.NotEmpty(t, privateKey)
		assert.NotEmpty(t, publicKey)

		// valid key?
		_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(publicKey))
		assert.NoError(t, err)

		_, err = ssh.ParsePrivateKey([]byte(privateKey))
		assert.NoError(t, err)
	})

	t.Run("with private key", func(t *testing.T) {
		defer prepareTestPrivateKeyFile()()

		ctx := context.Background()
		step := &stepCreateSSHKey{}
		config := dummyConfig()
		config.Comm.SSHPrivateKeyFile = dummyPrivateKeyFile
		state := dummyMinimumStateBag(&config)

		// run
		action := step.Run(ctx, state)

		publicKeys := state.Get("publicKeys").([]string)
		publicKey := publicKeys[0]
		privateKey, privateKeyExists := state.GetOk("privateKey")

		_, expectPublicKey, err := readTestKeyPair()
		if err != nil {
			panic(err)
		}

		assert.Equal(t, multistep.ActionContinue, action)
		assert.NotEmpty(t, publicKey)
		assert.Equal(t, expectPublicKey, publicKey)
		assert.Empty(t, privateKey)
		assert.False(t, privateKeyExists)
	})

	t.Run("with multiple private key", func(t *testing.T) {
		defer prepareTestPrivateKeyFile()()

		expectPrivateKey, expectPublicKey, err := readTestKeyPair()
		if err != nil {
			panic(err)
		}

		ctx := context.Background()
		step := &stepCreateSSHKey{}
		config := dummyConfig()
		config.Comm.SSHPrivateKeyFile = dummyPrivateKeyFile
		config.Comm.SSHPrivateKey = []byte(expectPrivateKey)
		state := dummyMinimumStateBag(&config)

		// run
		action := step.Run(ctx, state)

		publicKeys := state.Get("publicKeys").([]string)
		assert.Equal(t, multistep.ActionContinue, action)
		assert.Len(t, publicKeys, 2)
		assert.Equal(t, expectPublicKey, publicKeys[0])
		assert.Equal(t, expectPublicKey, publicKeys[1])
	})

	t.Run("skip generating public key ", func(t *testing.T) {
		defer prepareTestPrivateKeyFile()()

		ctx := context.Background()
		step := &stepCreateSSHKey{}
		config := dummyConfig()
		config.Comm.SSHPrivateKeyFile = dummyPrivateKeyFile
		config.DisableGeneratePublicKey = true
		state := dummyMinimumStateBag(&config)

		// run
		action := step.Run(ctx, state)
		_, ok := state.GetOk("publicKeys")

		assert.False(t, ok)
		assert.Equal(t, multistep.ActionContinue, action)
	})
}
