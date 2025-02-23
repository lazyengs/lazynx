package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	GeneratorsRequestMethod = "nx/generators"
)

// GeneratorsRequestOptions represents the options for a generators request.
type GeneratorsRequestOptions struct {
	IncludeHidden bool `json:"includeHidden"` // IncludeHidden specifies whether to include hidden generators.
	IncludeNgAdd  bool `json:"includeNgAdd"`  // IncludeNgAdd specifies whether to include ng-add generators.
}

// GeneratorsRequestParams represents the parameters for a generators request.
type GeneratorsRequestParams struct {
	Options *GeneratorsRequestOptions `json:"options,omitempty"` // Options are the generators request options.
}

// SendGeneratorsRequest sends a request to get generators.
func (c *Commander) SendGeneratorsRequest(ctx context.Context, params GeneratorsRequestParams) ([]nxtypes.GeneratorCollectionInfo, error) {
	var result []nxtypes.GeneratorCollectionInfo
	err := c.sendRequest(ctx, GeneratorsRequestMethod, params, &result)
	return result, err
}
