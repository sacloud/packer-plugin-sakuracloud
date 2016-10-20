package sakuracloud

import (
	"fmt"
	"github.com/sacloud/libsacloud/api"
	"log"
)

type Artifact struct {
	// The name of the archive
	archiveName string

	// The ID of the image
	archiveID int64

	// The name of the region
	zone string

	// The client for making API calls
	client *api.Client
}

func (*Artifact) BuilderId() string {
	return BuilderId
}

func (*Artifact) Files() []string {
	// No files with SakuraCloud
	return nil
}

func (a *Artifact) Id() string {
	return fmt.Sprintf("%d", a.archiveID)
}

func (a *Artifact) String() string {
	return fmt.Sprintf("A archive was created: '%v' (ID: %v) in zone '%v'", a.archiveName, a.archiveID, a.zone)
}

func (a *Artifact) State(name string) interface{} {
	return nil
}

func (a *Artifact) Destroy() error {
	log.Printf("Destroying image: %d (%s)", a.archiveID, a.archiveName)
	_, err := a.client.Archive.Delete(a.archiveID)
	return err
}
