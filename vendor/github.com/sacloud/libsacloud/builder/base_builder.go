package builder

import (
	"fmt"
	"github.com/sacloud/libsacloud/api"
)

type baseBuilder struct {
	client *api.Client
	errors []error
}

func (b *baseBuilder) toStringList(values []int64) []string {
	keys := []string{}
	for _, k := range values {
		keys = append(keys, fmt.Sprintf("%d", k))
	}
	return keys
}
