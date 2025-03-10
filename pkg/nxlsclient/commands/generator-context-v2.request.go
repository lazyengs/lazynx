package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	GeneratorContextV2RequestMethod = "nx/generatorContextV2"
)

// GeneratorContextV2Params represents the parameters for the generator context V2 request.
type GeneratorContextV2Params struct {
	Path string `json:"path"` // The path for which to generate the context.
}

// SendGeneratorContextV2Request sends a request to generate context V2 from a path.
// It takes a context and GeneratorContextV2Params as parameters and returns a pointer to GeneratorContext and an error.
func (c *Commander) SendGeneratorContextV2Request(ctx context.Context, params GeneratorContextV2Params) (*nxtypes.GeneratorContext, error) {
	var result *nxtypes.GeneratorContext
	err := c.sendRequest(ctx, GeneratorContextV2RequestMethod, params, &result)
	return result, err
}
