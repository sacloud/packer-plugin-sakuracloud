package sakuracloud

import (
	"fmt"
	"log"

	"github.com/sacloud/libsacloud/api"
)

// Artifact is the result of a build and is the metadata that documents
// what a builder actually created.
type Artifact struct {
	// archiveName is the name of the archive
	archiveName string

	// archiveID is the ID of the image
	archiveID int64

	// zone is the name of the region
	zone string

	// client is the SakuraCloud API client
	client *api.Client
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
	return fmt.Sprintf("%d", a.archiveID)
}

// String returns human-readable output that describes the artifact created.
func (a *Artifact) String() string {
	return fmt.Sprintf("A archive was created: '%v' (ID: %v) in zone '%v'", a.archiveName, a.archiveID, a.zone)
}

// State allows the caller to ask for builder specific state information
// relating to the artifact instance.
func (a *Artifact) State(name string) interface{} {
	return nil
}

// Destroy deletes the artifact.
func (a *Artifact) Destroy() error {
	log.Printf("Destroying image: %d (%s)", a.archiveID, a.archiveName)
	_, err := a.client.Archive.Delete(a.archiveID)
	return err
}
