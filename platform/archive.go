package platform

import (
	"context"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/query"
	"github.com/sacloud/iaas-api-go/types"
	archiveBuilder "github.com/sacloud/iaas-service-go/archive/builder"
)

// Archive is responsible for API calls of archive handling
type Archive interface {
	Delete(ctx context.Context, id types.ID) error
	Create(ctx context.Context, req *CreateArchiveRequest) (*iaas.Archive, error)
	Transfer(ctx context.Context, zone string, req *TransferArchiveRequest) (*iaas.Archive, error)
}

type archiveClient struct {
	caller    iaas.APICaller
	archiveOp iaas.ArchiveAPI
	diskOp    iaas.DiskAPI
	zone      string
}

func newArchiveClient(caller iaas.APICaller, zone string) *archiveClient {
	return &archiveClient{
		caller:    caller,
		archiveOp: iaas.NewArchiveOp(caller),
		diskOp:    iaas.NewDiskOp(caller),
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

func (c *archiveClient) Create(ctx context.Context, req *CreateArchiveRequest) (*iaas.Archive, error) {
	tags, err := c.tagsFromAncestors(ctx, req.Tags, req.DiskID)
	if err != nil {
		return nil, err
	}

	builder := (&archiveBuilder.Director{
		Name:         req.Name,
		Description:  req.Description,
		Tags:         tags,
		SourceDiskID: req.DiskID,
		Client:       archiveBuilder.NewAPIClient(c.caller),
	}).Builder()

	return builder.Build(ctx, c.zone)
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

func (c *archiveClient) Transfer(ctx context.Context, zone string, req *TransferArchiveRequest) (*iaas.Archive, error) {
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
