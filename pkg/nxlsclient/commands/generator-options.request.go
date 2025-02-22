package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	GeneratorOptionsRequestMethod = "nx/generatorOptions"
)

type GeneratorOptionsRequestOptions struct {
	Collection string `json:"collection"`
	Name       string `json:"name"`
	Path       string `json:"path"`
}

type GeneratorOptionsRequestParams struct {
	Options GeneratorOptionsRequestOptions `json:"options"`
}

func (c *Commander) SendGeneratorOptionsRequest(ctx context.Context, params GeneratorOptionsRequestParams) ([]nxtypes.Option, error) {
	var result []nxtypes.Option
	err := c.sendRequest(ctx, GeneratorOptionsRequestMethod, params, &result)
	return result, err
}
