package sakuracloud

import (
	"context"
	"testing"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

func TestStepCreateSSHKey(t *testing.T) {

	ctx := context.Background()
	step := &stepCreateSSHKey{}

	state := dummyMinimumStateBag(nil)

	// run
	action := step.Run(ctx, state)

	privateKey := state.Get("privateKey").(string)
	publicKey := state.Get("ssh_public_key").(string)

	assert.Equal(t, multistep.ActionContinue, action)
	assert.NotEmpty(t, privateKey)
	assert.NotEmpty(t, publicKey)

	// valid key?
	_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(publicKey))
	assert.NoError(t, err)

	_, err = ssh.ParsePrivateKey([]byte(privateKey))
	assert.NoError(t, err)
}
