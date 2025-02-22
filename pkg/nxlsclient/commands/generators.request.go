package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	GeneratorsRequestMethod = "nx/generators"
)

type GeneratorsRequestOptions struct {
	IncludeHidden bool `json:"includeHidden"`
	IncludeNgAdd  bool `json:"includeNgAdd"`
}

type GeneratorsRequestParams struct {
	Options *GeneratorsRequestOptions `json:"options,omitempty"`
}

func (c *Commander) SendGeneratorsRequest(ctx context.Context, params GeneratorsRequestParams) ([]nxtypes.GeneratorCollectionInfo, error) {
	var result []nxtypes.GeneratorCollectionInfo
	err := c.sendRequest(ctx, GeneratorsRequestMethod, params, &result)
	return result, err
}
