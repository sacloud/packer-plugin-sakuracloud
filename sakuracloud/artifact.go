package sakuracloud

import (
	"context"
	"fmt"
	"log"

	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/packer-builder-sakuracloud/iaas"
)

// Artifact is the result of a build and is the metadata that documents
// what a builder actually created.
type Artifact struct {
	// archiveName is the name of the archive
	archiveName string

	// archiveID is the ID of the image
	archiveID types.ID

	// client is the SakuraCloud API client
	client iaas.Archive
}

// BuilderId returns the ID of the builder that was used to create this artifact.
func (*Artifact) BuilderId() string {
	return BuilderId
}

// Files returns the set of files that comprise this artifact. If an
// artifact is not made up of files, then this will be empty.
func (*Artifact) Files() []string {
	// No files with SakuraCloud
	return nil
}

// Id for the artifact, if it has one.
func (a *Artifact) Id() string {
	return a.archiveID.String()
}

// String returns human-readable output that describes the artifact created.
func (a *Artifact) String() string {
	return fmt.Sprintf("A archive was created: %s (ID: %q)", a.archiveName, a.archiveID)
}

// State allows the caller to ask for builder specific state information
// relating to the artifact instance.
func (a *Artifact) State(name string) interface{} {
	return nil
}

// Destroy deletes the artifact.
func (a *Artifact) Destroy() error {
	log.Printf("Destroying image: %d (%s)", a.archiveID, a.archiveName)
	return a.client.Delete(context.TODO(), a.archiveID)
}
