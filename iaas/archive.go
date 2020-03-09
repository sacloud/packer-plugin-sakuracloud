package iaas

import (
	"context"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/accessor"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/query"
	"github.com/sacloud/libsacloud/v2/utils/setup"
)

// Archive is responsible for API calls of archive handling
type Archive interface {
	//Read(ctx context.Context, zone string, id types.ID) (*sacloud.Archive, error)
	//Create(ctx context.Context, zone string, param *sacloud.ArchiveCreateRequest) (*sacloud.Archive, error)
	Delete(ctx context.Context, id types.ID) error
	Create(ctx context.Context, req *CreateArchiveRequest) (*sacloud.Archive, error)
}

func NewArchiveClient(caller sacloud.APICaller, zone string) Archive {
	return newArchiveClient(caller, zone)
}

type archiveClient struct {
	caller    sacloud.APICaller
	archiveOp sacloud.ArchiveAPI
	diskOp    sacloud.DiskAPI
	zone      string
}

func newArchiveClient(caller sacloud.APICaller, zone string) *archiveClient {
	return &archiveClient{
		caller:    caller,
		archiveOp: sacloud.NewArchiveOp(caller),
		diskOp:    sacloud.NewDiskOp(caller),
		zone:      zone,
	}
}

func (c *archiveClient) Delete(ctx context.Context, id types.ID) error {
	return c.archiveOp.Delete(ctx, c.zone, id)
}

func (c *archiveClient) Create(ctx context.Context, req *CreateArchiveRequest) (*sacloud.Archive, error) {
	tags := req.Tags
	if len(tags) == 0 {
		parentArchiveID, err := query.GetPublicArchiveIDFromAncestors(ctx, c.zone, query.NewArchiveSourceReader(c.caller), req.DiskID)
		if err != nil {
			return nil, err
		}
		parentArchive, err := c.archiveOp.Read(ctx, c.zone, parentArchiveID)
		if err != nil {
			return nil, err
		}
		tags = parentArchive.Tags
	}

	archiveBuilder := &setup.RetryableSetup{
		IsWaitForCopy: true,
		Create: func(ctx context.Context, zone string) (accessor.ID, error) {
			return c.archiveOp.Create(ctx, zone, &sacloud.ArchiveCreateRequest{
				SourceDiskID: req.DiskID,
				Name:         req.Name,
				Tags:         tags,
				Description:  req.Description,
			})
		},
		Read: func(ctx context.Context, zone string, id types.ID) (interface{}, error) {
			return c.archiveOp.Read(ctx, zone, id)
		},
		Delete: func(ctx context.Context, zone string, id types.ID) error {
			return c.archiveOp.Delete(ctx, zone, id)
		},
		RetryCount: 3,
	}
	res, err := archiveBuilder.Setup(ctx, c.zone)
	if err != nil {
		return nil, err
	}
	return res.(*sacloud.Archive), nil
}

type CreateArchiveRequest struct {
	DiskID      types.ID
	Name        string
	Tags        types.Tags
	Description string
}
