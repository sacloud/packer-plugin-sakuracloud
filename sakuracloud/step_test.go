package sakuracloud

import (
	"os"
	"time"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/sacloud/libsacloud/builder"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/libsacloud/sacloud/ostype"
)

var (
	dummyDiskID             = int64(111111111111)
	dummyCreatedArchiveID   = int64(222222222222)
	dummyReadArchiveID      = int64(333333333333)
	dummyServerID           = int64(444444444444)
	dummyISOImageID         = int64(555555555555)
	dummyArchiveID          = int64(666666666666)
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

	testMinimumConfigValues = map[string]interface{}{
		"access_token":        "aaaa",
		"access_token_secret": "bbbb",
		"zone":                "is1a",
		"os_type":             "centos",
	}
)

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

func createDummyServer(id int64, status string) *sacloud.Server {
	server := &sacloud.Server{Resource: sacloud.NewResource(id)}
	server.Instance = &sacloud.Instance{
		EServerInstanceStatus: &sacloud.EServerInstanceStatus{
			Status: status,
		},
	}
	return server
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

type dummyServerClient struct {
	readFunc           func(int64) (*sacloud.Server, error)
	stopFunc           func(int64) (bool, error)
	shutdownFunc       func(int64) (bool, error)
	sleepUntilDownFunc func(int64, time.Duration) error
	deleteFunc         func(int64) (*sacloud.Server, error)
	deleteWithDiskFunc func(int64, []int64) (*sacloud.Server, error)
	getVNCProxyFunc    func(int64) (*sacloud.VNCProxyResponse, error)
}

func (t *dummyServerClient) Read(id int64) (*sacloud.Server, error) {
	if t.readFunc == nil {
		return nil, nil
	}
	return t.readFunc(id)
}

func (t *dummyServerClient) Stop(id int64) (bool, error) {
	if t.stopFunc == nil {
		return false, nil
	}
	return t.stopFunc(id)
}

func (t *dummyServerClient) Shutdown(id int64) (bool, error) {
	if t.shutdownFunc == nil {
		return false, nil
	}
	return t.shutdownFunc(id)
}

func (t *dummyServerClient) SleepUntilDown(id int64, timeout time.Duration) error {
	if t.sleepUntilDownFunc == nil {
		return nil
	}
	return t.sleepUntilDownFunc(id, timeout)
}

func (t *dummyServerClient) Delete(id int64) (*sacloud.Server, error) {
	if t.deleteFunc == nil {
		return nil, nil
	}
	return t.deleteFunc(id)
}

func (t *dummyServerClient) DeleteWithDisk(id int64, disks []int64) (*sacloud.Server, error) {
	if t.deleteWithDiskFunc == nil {
		return nil, nil
	}
	return t.deleteWithDiskFunc(id, disks)
}

func (t *dummyServerClient) GetVNCProxy(serverID int64) (*sacloud.VNCProxyResponse, error) {
	if t.getVNCProxyFunc == nil {
		return nil, nil
	}
	return t.getVNCProxyFunc(serverID)
}

type dummyArchiveClient struct {
	newFunc               func() *sacloud.Archive
	readFunc              func(int64) (*sacloud.Archive, error)
	createFunc            func(*sacloud.Archive) (*sacloud.Archive, error)
	sleepWhileCopyingFunc func(int64, time.Duration) error
	deleteFunc            func(int64) (*sacloud.Archive, error)
}

func (t *dummyArchiveClient) New() *sacloud.Archive {
	if t.newFunc == nil {
		return &sacloud.Archive{}
	}
	return t.newFunc()
}

func (t *dummyArchiveClient) Read(id int64) (*sacloud.Archive, error) {
	if t.readFunc == nil {
		return nil, nil
	}
	return t.readFunc(id)
}

func (t *dummyArchiveClient) Create(param *sacloud.Archive) (*sacloud.Archive, error) {
	if t.createFunc == nil {
		return nil, nil
	}
	return t.createFunc(param)
}

func (t *dummyArchiveClient) SleepWhileCopying(id int64, timeout time.Duration) error {
	if t.sleepWhileCopyingFunc == nil {
		return nil
	}
	return t.sleepWhileCopyingFunc(id, timeout)
}

func (t *dummyArchiveClient) Delete(id int64) (*sacloud.Archive, error) {
	if t.deleteFunc == nil {
		return nil, nil
	}
	return t.deleteFunc(id)
}

type dummyDiskClient struct {
	getPublicArchiveIDFromAncestorsFunc func(int64) (int64, bool)
}

func (t *dummyDiskClient) GetPublicArchiveIDFromAncestors(id int64) (int64, bool) {
	if t.getPublicArchiveIDFromAncestorsFunc == nil {
		return 0, false
	}
	return t.getPublicArchiveIDFromAncestorsFunc(id)
}

type dummyISOImageClient struct {
	newFunc         func() *sacloud.CDROM
	createFunc      func(*sacloud.CDROM) (*sacloud.CDROM, *sacloud.FTPServer, error)
	readFunc        func(int64) (*sacloud.CDROM, error)
	setEmptyFunc    func()
	setNameLikeFunc func(string)
	findFunc        func() (*sacloud.SearchResponse, error)
	closeFTPFunc    func(int64) (bool, error)
}

func (t *dummyISOImageClient) New() *sacloud.CDROM {
	if t.newFunc == nil {
		return &sacloud.CDROM{}
	}
	return t.newFunc()
}

func (t *dummyISOImageClient) Create(value *sacloud.CDROM) (*sacloud.CDROM, *sacloud.FTPServer, error) {
	if t.createFunc == nil {
		return nil, nil, nil
	}
	return t.createFunc(value)
}

func (t *dummyISOImageClient) Read(id int64) (*sacloud.CDROM, error) {
	if t.readFunc == nil {
		return nil, nil
	}
	return t.readFunc(id)
}

func (t *dummyISOImageClient) SetEmpty() {
	if t.setEmptyFunc != nil {
		t.setEmptyFunc()
	}
}

func (t *dummyISOImageClient) SetNameLike(name string) {
	if t.setNameLikeFunc != nil {
		t.setNameLikeFunc(name)
	}
}

func (t *dummyISOImageClient) Find() (*sacloud.SearchResponse, error) {
	if t.findFunc == nil {
		return nil, nil
	}
	return t.findFunc()
}

func (t *dummyISOImageClient) CloseFTP(id int64) (bool, error) {
	if t.closeFTPFunc == nil {
		return false, nil
	}
	return t.closeFTPFunc(id)
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

type dummyBuilderAPIClient struct {
	getBySpecFunc           func(core int, memGB int, gen sacloud.PlanGenerations) (*sacloud.ProductServer, error)
	archiveFindByOSTypeFunc func(os ostype.ArchiveOSTypes) (*sacloud.Archive, error)
}

func (t *dummyBuilderAPIClient) ServerNew() *sacloud.Server                         { return nil }
func (t *dummyBuilderAPIClient) ServerRead(serverID int64) (*sacloud.Server, error) { return nil, nil }
func (t *dummyBuilderAPIClient) ServerCreate(value *sacloud.Server) (*sacloud.Server, error) {
	return nil, nil
}
func (t *dummyBuilderAPIClient) ServerSleepUntilUp(serverID int64, timeout time.Duration) error {
	return nil
}
func (t *dummyBuilderAPIClient) ServerInsertCDROM(serverID int64, cdromID int64) (bool, error) {
	return false, nil
}
func (t *dummyBuilderAPIClient) ServerBoot(serverID int64) (bool, error) { return false, nil }
func (t *dummyBuilderAPIClient) SSHKeyNew() *sacloud.SSHKey              { return nil }
func (t *dummyBuilderAPIClient) SSHKeyCreate(value *sacloud.SSHKey) (*sacloud.SSHKey, error) {
	return nil, nil
}
func (t *dummyBuilderAPIClient) SSHKeyDelete(sshKeyID int64) (*sacloud.SSHKey, error) { return nil, nil }
func (t *dummyBuilderAPIClient) SSHKeyGenerate(name string, passPhrase string, desc string) (*sacloud.SSHKeyGenerated, error) {
	return nil, nil
}
func (t *dummyBuilderAPIClient) NoteNew() *sacloud.Note { return nil }
func (t *dummyBuilderAPIClient) NoteCreate(value *sacloud.Note) (*sacloud.Note, error) {
	return nil, nil
}
func (t *dummyBuilderAPIClient) NoteDelete(noteID int64) (*sacloud.Note, error) { return nil, nil }
func (t *dummyBuilderAPIClient) DiskNew() *sacloud.Disk                         { return nil }
func (t *dummyBuilderAPIClient) DiskNewCondig() *sacloud.DiskEditValue          { return nil }
func (t *dummyBuilderAPIClient) DiskCreate(value *sacloud.Disk) (*sacloud.Disk, error) {
	return nil, nil
}
func (t *dummyBuilderAPIClient) DiskCreateWithConfig(value *sacloud.Disk, config *sacloud.DiskEditValue, bootAtAvailable bool) (*sacloud.Disk, error) {
	return nil, nil
}
func (t *dummyBuilderAPIClient) DiskSleepWhileCopying(id int64, timeout time.Duration) error {
	return nil
}
func (t *dummyBuilderAPIClient) DiskConnectToServer(diskID int64, serverID int64) (bool, error) {
	return false, nil
}
func (t *dummyBuilderAPIClient) InterfaceConnectToPacketFilter(interfaceID int64, packetFilterID int64) (bool, error) {
	return false, nil
}
func (t *dummyBuilderAPIClient) InterfaceSetDisplayIPAddress(interfaceID int64, ip string) (bool, error) {
	return false, nil
}
func (t *dummyBuilderAPIClient) GetTimeoutDuration() time.Duration { return time.Hour }

func (t *dummyBuilderAPIClient) ServerPlanGetBySpec(core int, memGB int, gen sacloud.PlanGenerations) (*sacloud.ProductServer, error) {
	if t.getBySpecFunc == nil {
		return nil, nil
	}
	return t.getBySpecFunc(core, memGB, gen)
}

func (t *dummyBuilderAPIClient) ArchiveFindByOSType(os ostype.ArchiveOSTypes) (*sacloud.Archive, error) {
	if t.archiveFindByOSTypeFunc == nil {
		return nil, nil
	}
	return t.archiveFindByOSTypeFunc(os)
}

type dummyBuilderFactory struct {
	builder *dummyServerBuilder
}

func (t *dummyBuilderFactory) createServerBuilder(multistep.StateBag) serverBuilder {
	return t.builder
}

type dummyServerBuilder struct {
	result *builder.ServerBuildResult
	err    error
}

func (t *dummyServerBuilder) Build() (*builder.ServerBuildResult, error) {
	return t.result, t.err
}

func (t *dummyServerBuilder) init() {
	t.result = nil
	t.err = nil
}
