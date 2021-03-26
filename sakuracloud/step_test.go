package sakuracloud

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/sacloud/libsacloud/v2/sacloud/fake"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/packer-plugin-sakuracloud/iaas"
	"golang.org/x/crypto/ssh"
)

var (
	dummyDiskID         = types.ID(111111111111)
	dummyArchiveID      = types.ID(666666666666)
	dummyPrivateKeyFile = "packer-test-private-key"

	testMinimumConfigValues = map[string]interface{}{
		"access_token":        "aaaa",
		"access_token_secret": "bbbb",
		"zone":                "is1a",
		"os_type":             "centos",
	}
)

func prepareTestPrivateKeyFile() func() {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	privBlk := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(priv),
	}
	deferFunc := func() {
		if err := os.Remove(dummyPrivateKeyFile); err != nil {
			panic(err)
		}
	}
	if err := os.WriteFile(dummyPrivateKeyFile, pem.EncodeToMemory(&privBlk), 0600); err != nil {
		defer deferFunc()
	}
	return deferFunc
}

func readTestKeyPair() (string, string, error) {
	bytes, err := os.ReadFile(dummyPrivateKeyFile)
	if err != nil {
		return "", "", err
	}
	signer, err := ssh.ParsePrivateKey(bytes)
	if err != nil {
		return "", "", err
	}
	return string(bytes), string(ssh.MarshalAuthorizedKey(signer.PublicKey())), nil
}

func dummyStateBag() multistep.StateBag {
	return new(multistep.BasicStateBag)
}

func dummyConfig() Config {
	return dummyConfigWithValues(testMinimumConfigValues)
}

func dummyConfigWithValues(values map[string]interface{}) Config {
	conf, _, err := NewConfig(values)
	if err != nil {
		panic(err)
	}
	return *conf
}

func dummyMinimumStateBag(t *testing.T, config *Config) multistep.StateBag {
	state := dummyStateBag()
	state.Put("ui", packer.TestUi(t))
	if config == nil {
		state.Put("config", dummyConfig())
	} else {
		state.Put("config", *config)
	}

	// set fake API client
	fake.SwitchFactoryFuncToFake()
	iaasClient := iaas.NewClient("fake", "fake", "is1a")
	iaasClient.FTPS = &fakeFTPSClient{}
	state.Put("iaasClient", iaasClient)
	return state
}

type fakeFTPSClient struct {
	connectErr   error
	loginErr     error
	storeFileErr error
	quitErr      error
}

// Connect .
func (c *fakeFTPSClient) Connect(string, int) error {
	return c.connectErr
}

// Login .
func (c *fakeFTPSClient) Login(string, string) error {
	return c.loginErr
}

// StoreFile .
func (c *fakeFTPSClient) StoreFile(string, *os.File) error {
	return c.storeFileErr
}

// Quit .
func (c *fakeFTPSClient) Quit() error {
	return c.quitErr
}
