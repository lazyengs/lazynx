package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	GeneratorContextV2RequestMethod = "nx/generatorContextV2"
)

type GeneratorContextV2Params struct {
	Path string `json:"path"`
}

func (c *Commander) SendGeneratorContextV2Request(ctx context.Context, params GeneratorContextV2Params) (*nxtypes.GeneratorContext, error) {
	var result *nxtypes.GeneratorContext
	err := c.sendRequest(ctx, GeneratorContextV2RequestMethod, params, &result)
	return result, err
}
