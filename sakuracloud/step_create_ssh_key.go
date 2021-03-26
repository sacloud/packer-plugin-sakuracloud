package sakuracloud

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"golang.org/x/crypto/ssh"
)

type stepCreateSSHKey struct {
	Debug        bool
	DebugKeyPath string
}

func (s *stepCreateSSHKey) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)

	stepStartMsg(ui, s.Debug, "CreateSSHKey")

	var privateKeys []string

	if c.Comm.SSHPrivateKeyFile != "" {
		bytes, err := os.ReadFile(c.Comm.SSHPrivateKeyFile)
		if err != nil {
			err := fmt.Errorf("Error reading ssh key file: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
		privateKeys = append(privateKeys, string(bytes))
	}

	if len(c.Comm.SSHPrivateKey) > 0 {
		privateKeys = append(privateKeys, string(c.Comm.SSHPrivateKey))
	}

	if c.Comm.SSHPrivateKeyFile == "" && len(c.Comm.SSHPrivateKey) == 0 {
		ui.Say("\tCreating temporary SSH key for instance...")
		pkey, err := s.generatePrivateKey()
		if err != nil {
			err := fmt.Errorf("Error creating temporary ssh key: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
		privateKeys = append(privateKeys, pkey)
		state.Put("privateKey", pkey)
	}

	if !c.DisableGeneratePublicKey {
		// generate public keys for SakuraCloud Disk Edit API
		var publicKeys []string
		if len(c.Comm.SSHPublicKey) > 0 {
			publicKeys = append(publicKeys, string(c.Comm.SSHPublicKey))
		}
		for _, privateKey := range privateKeys {
			signer, err := ssh.ParsePrivateKey([]byte(privateKey))
			if err != nil {
				err := fmt.Errorf("Error creating ssh public key: %s", err)
				state.Put("error", err)
				ui.Error(err.Error())
				return multistep.ActionHalt
			}
			pubKey := string(ssh.MarshalAuthorizedKey(signer.PublicKey()))
			publicKeys = append(publicKeys, pubKey)
		}
		state.Put("publicKeys", publicKeys)
	}

	if s.Debug {
		pkey, ok := state.GetOk("privateKey")
		if ok {
			ui.Message(fmt.Sprintf("Saving key for debug purposes: %s", s.DebugKeyPath))
			if err := os.WriteFile(s.DebugKeyPath, []byte(pkey.(string)), 0600); err != nil {
				state.Put("error", fmt.Errorf("Error saving debug key: %s", err))
				return multistep.ActionHalt
			}
		}
	}

	stepEndMsg(ui, s.Debug, "CreateSSHKey")
	return multistep.ActionContinue
}

func (s *stepCreateSSHKey) Cleanup(state multistep.StateBag) {
	// no cleanup
}

func (s *stepCreateSSHKey) generatePrivateKey() (string, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", err
	}

	privBlk := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(priv),
	}
	return string(pem.EncodeToMemory(&privBlk)), nil
}
