package iaas

import (
	"fmt"
	"time"

	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/builder"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/libsacloud/sacloud/ostype"
	"github.com/sacloud/packer-builder-sakuracloud/version"
)

// Client represents SakuraCloud API Client
type Client struct {
	Builder  ServerBuilder
	Server   ServerClient
	Archive  ArchiveClient
	Disk     DiskClient
	ISOImage ISOImageClient
	Basic    BasicClient
	client   *api.Client
}

// NewClient returns new SakuraCloud API Client
func NewClient(token, secret, zone string, apiWaitTimeout time.Duration) *Client {

	client := api.NewClient(token, secret, zone)
	client.UserAgent = fmt.Sprintf("packer_for_sakuracloud:v%s", version.Version)
	client.DefaultTimeoutDuration = apiWaitTimeout

	return &Client{
		Builder:  &builderImpl{client: client},
		Server:   client.Server,
		Archive:  client.Archive,
		Disk:     client.Disk,
		ISOImage: client.CDROM,
		Basic:    &basicClient{client: client},
	}
}

type builderImpl struct {
	client *api.Client
}

func (b *builderImpl) FromBlankDisk(name string) builder.BlankDiskServerBuilder {
	return builder.ServerBlankDisk(b.client, name)
}

func (b *builderImpl) FromArchive(name string, sourceArchiveID int64) builder.CommonServerBuilder {
	return builder.ServerFromArchive(b.client, name, sourceArchiveID)
}

func (b *builderImpl) FromDisk(name string, sourceDiskID int64) builder.CommonServerBuilder {
	return builder.ServerFromDisk(b.client, name, sourceDiskID)
}

func (b *builderImpl) FromPublicArchiveWindows(os ostype.ArchiveOSTypes, name string) builder.PublicArchiveWindowsServerBuilder {
	return builder.ServerPublicArchiveWindows(b.client, os, name)
}

func (b *builderImpl) FromPublicArchiveFixedUnix(os ostype.ArchiveOSTypes, name string) builder.FixedUnixArchiveServerBuilder {
	return builder.ServerPublicArchiveFixedUnix(b.client, os, name)
}

func (b *builderImpl) FromPublicArchiveUnix(os ostype.ArchiveOSTypes, name string, password string) builder.PublicArchiveUnixServerBuilder {
	return builder.ServerPublicArchiveUnix(b.client, os, name, password)
}

// BasicClient is responsible for basic functions of API client
type BasicClient interface {
	Zone() string
}

type basicClient struct {
	client *api.Client
}

func (b *basicClient) Zone() string {
	return b.client.Zone
}

// ServerBuilder is responsible for API calls of server build
type ServerBuilder interface {
	FromBlankDisk(name string) builder.BlankDiskServerBuilder
	FromArchive(name string, sourceArchiveID int64) builder.CommonServerBuilder
	FromDisk(name string, sourceDiskID int64) builder.CommonServerBuilder
	FromPublicArchiveWindows(os ostype.ArchiveOSTypes, name string) builder.PublicArchiveWindowsServerBuilder
	FromPublicArchiveFixedUnix(os ostype.ArchiveOSTypes, name string) builder.FixedUnixArchiveServerBuilder
	FromPublicArchiveUnix(os ostype.ArchiveOSTypes, name string, password string) builder.PublicArchiveUnixServerBuilder
}

// ServerClient is responsible for API calls of server handling
type ServerClient interface {
	Read(id int64) (*sacloud.Server, error)
	Stop(id int64) (bool, error)
	Shutdown(id int64) (bool, error)
	SleepUntilDown(id int64, timeout time.Duration) error
	Delete(id int64) (*sacloud.Server, error)
	DeleteWithDisk(id int64, disks []int64) (*sacloud.Server, error)
	GetVNCProxy(serverID int64) (*sacloud.VNCProxyResponse, error)
}

// ArchiveClient is responsible for API calls of archive handling
type ArchiveClient interface {
	New() *sacloud.Archive
	Read(id int64) (*sacloud.Archive, error)
	Create(param *sacloud.Archive) (*sacloud.Archive, error)
	SleepWhileCopying(id int64, timeout time.Duration) error
	Delete(id int64) (*sacloud.Archive, error)
}

// DiskClient is responsible for API calls of disk handling
type DiskClient interface {
	GetPublicArchiveIDFromAncestors(id int64) (int64, bool)
}

// ISOImageClient is responsible for API calls of ISO-image handling
type ISOImageClient interface {
	New() *sacloud.CDROM
	Create(value *sacloud.CDROM) (*sacloud.CDROM, *sacloud.FTPServer, error)
	Read(id int64) (*sacloud.CDROM, error)
	SetEmpty()
	SetNameLike(name string)
	Find() (*sacloud.SearchResponse, error)
	CloseFTP(id int64) (bool, error)
}
