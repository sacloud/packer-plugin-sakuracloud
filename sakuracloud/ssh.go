package sakuracloud

import (
	"fmt"

	"github.com/mitchellh/multistep"
	"golang.org/x/crypto/ssh"
)

func commHost(state multistep.StateBag) (string, error) {
	ipAddress := state.Get("server_ip").(string)
	return ipAddress, nil
}

func sshConfig(state multistep.StateBag) (*ssh.ClientConfig, error) {
	config := state.Get("config").(Config)
	privateKey := state.Get("ssh_private_key").(string)

	signer, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return nil, fmt.Errorf("Error setting up SSH config: %s", err)
	}

	auth := []ssh.AuthMethod{ssh.PublicKeys(signer)}
	if config.Password != "" {
		config.Comm.SSHPassword = config.Password
		auth = append(auth, ssh.Password(config.Password))
	}

	return &ssh.ClientConfig{
		User: config.Comm.SSHUsername,
		Auth: auth,
	}, nil
}
