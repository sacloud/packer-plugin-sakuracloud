package sakuracloud

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"golang.org/x/crypto/ssh"
)

var (
	dummyDiskID             = types.ID(111111111111)
	dummyCreatedArchiveID   = types.ID(222222222222)
	dummyReadArchiveID      = types.ID(333333333333)
	dummyServerID           = types.ID(444444444444)
	dummyISOImageID         = types.ID(555555555555)
	dummyArchiveID          = types.ID(666666666666)
	dummyISOPath            = "test.iso"
	dummyArchiveName        = "testArchive"
	dummyArchiveTags        = []string{"archive1", "archive2"}
	dummyParentArchiveTags  = []string{"parent1", "parent2"}
	dummyDescription        = "testArchiveDescription"
	dummyServerCore         = 2
	dummyServerMemory       = 4
	dummyServerIP           = "192.2.0.11"
	dummyServerDefaultRoute = "192.2.0.1"
	dummyServerNwMaskLen    = 24
	dummyServerPassword     = "p@ssw0rd"
	dummyDiskSize           = 20
	dummyDNSServers         = []string{"ns1.example.com", "ns2.example.com"}
	dummySSHKeyBody         = "ssh-rsa AAAA..."
	dummyPrivateKeyFile     = "packer-test-private-key"

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
	if err := ioutil.WriteFile(dummyPrivateKeyFile, pem.EncodeToMemory(&privBlk), 0600); err != nil {
		defer deferFunc()
	}
	return deferFunc
}

func readTestKeyPair() (string, string, error) {
	bytes, err := ioutil.ReadFile(dummyPrivateKeyFile)
	if err != nil {
		return "", "", err
	}
	signer, err := ssh.ParsePrivateKey(bytes)
	if err != nil {
		return "", "", err
	}
	return string(bytes), string(ssh.MarshalAuthorizedKey(signer.PublicKey())), nil
}

func dummyUI() packer.Ui {
	return new(packer.NoopUi)
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

func dummyMinimumStateBag(config *Config) multistep.StateBag {
	state := dummyStateBag()
	state.Put("ui", dummyUI())
	if config == nil {
		state.Put("config", dummyConfig())
	} else {
		state.Put("config", *config)
	}
	return state
}

type dummyBasicClient struct {
	zoneFunc func() string
}

func (t *dummyBasicClient) Zone() string {
	if t.zoneFunc == nil {
		return "is1a"
	}
	return t.zoneFunc()
}

type dummyFTPSClient struct {
	connectFunc   func(string, int) error
	loginFunc     func(string, string) error
	storeFileFunc func(string, *os.File) error
	quitFunc      func() error
}

func (t *dummyFTPSClient) Connect(host string, port int) error {
	if t.connectFunc == nil {
		return nil
	}
	return t.connectFunc(host, port)
}

func (t *dummyFTPSClient) Login(user, password string) error {
	if t.loginFunc == nil {
		return nil
	}
	return t.loginFunc(user, password)
}

func (t *dummyFTPSClient) StoreFile(remoteFilepath string, file *os.File) error {
	if t.storeFileFunc == nil {
		return nil
	}
	return t.storeFileFunc(remoteFilepath, file)
}

func (t *dummyFTPSClient) Quit() error {
	if t.quitFunc == nil {
		return nil
	}
	return t.quitFunc()
}
