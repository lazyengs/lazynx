package commands

import (
	"context"

	nxtypes "github.com/lazyengs/lazynx/pkg/nxlsclient/nx-types"
)

const (
	GeneratorOptionsRequestMethod = "nx/generatorOptions"
)

// GeneratorOptionsRequestOptions represents the options for a generator request.
type GeneratorOptionsRequestOptions struct {
	Collection string `json:"collection"` // Collection is the name of the collection.
	Name       string `json:"name"`       // Name is the name of the generator.
	Path       string `json:"path"`       // Path is the path where the generator should be applied.
}

// GeneratorOptionsRequestParams represents the parameters for a generator options request.
type GeneratorOptionsRequestParams struct {
	Options GeneratorOptionsRequestOptions `json:"options"` // Options are the generator options.
}

// SendGeneratorOptionsRequest sends a request to get generator options.
func (c *Commander) SendGeneratorOptionsRequest(ctx context.Context, params GeneratorOptionsRequestParams) ([]nxtypes.Option, error) {
	var result []nxtypes.Option
	err := c.sendRequest(ctx, GeneratorOptionsRequestMethod, params, &result)
	return result, err
}
