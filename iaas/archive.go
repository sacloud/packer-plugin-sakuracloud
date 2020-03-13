package iaas

import (
	"context"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/accessor"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	archiveBuilder "github.com/sacloud/libsacloud/v2/utils/builder/archive"
	"github.com/sacloud/libsacloud/v2/utils/query"
	"github.com/sacloud/libsacloud/v2/utils/setup"
)

// Archive is responsible for API calls of archive handling
type Archive interface {
	Delete(ctx context.Context, id types.ID) error
	Create(ctx context.Context, req *CreateArchiveRequest) (*sacloud.Archive, error)
	Transfer(ctx context.Context, zone string, req *TransferArchiveRequest) (*sacloud.Archive, error)
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

func (c *archiveClient) tagsFromAncestors(ctx context.Context, currentTags types.Tags, id types.ID) (types.Tags, error) {
	if len(currentTags) != 0 {
		return currentTags, nil
	}
	parentArchiveID, err := query.GetPublicArchiveIDFromAncestors(ctx, c.zone, query.NewArchiveSourceReader(c.caller), id)
	if err != nil {
		return nil, err
	}
	parentArchive, err := c.archiveOp.Read(ctx, c.zone, parentArchiveID)
	if err != nil {
		return nil, err
	}
	return parentArchive.Tags, nil
}

func (c *archiveClient) Create(ctx context.Context, req *CreateArchiveRequest) (*sacloud.Archive, error) {
	tags, err := c.tagsFromAncestors(ctx, req.Tags, req.DiskID)
	if err != nil {
		return nil, err
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

// CreateArchiveRequest is a parameter of creating SakuraCloud Archive
type CreateArchiveRequest struct {
	DiskID      types.ID
	Name        string
	Tags        types.Tags
	Description string
}

// TransferArchiveRequest is a parameter of creating SakuraCloud Archive
type TransferArchiveRequest struct {
	Name              string
	Tags              types.Tags
	Description       string
	SourceArchiveID   types.ID
	SourceArchiveZone string
}

func (c *archiveClient) Transfer(ctx context.Context, zone string, req *TransferArchiveRequest) (*sacloud.Archive, error) {
	tags, err := c.tagsFromAncestors(ctx, req.Tags, req.SourceArchiveID)
	if err != nil {
		return nil, err
	}

	builder := &archiveBuilder.TransferArchiveBuilder{
		Name:              req.Name,
		Description:       req.Description,
		Tags:              tags,
		SourceArchiveID:   req.SourceArchiveID,
		SourceArchiveZone: req.SourceArchiveZone,
		Client:            archiveBuilder.NewAPIClient(c.caller),
	}

	return builder.Build(ctx, zone)
}
