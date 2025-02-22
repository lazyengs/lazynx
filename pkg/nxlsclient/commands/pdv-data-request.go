package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	PDVDataRequestMethod = "nx/pdvData"
)

type PDVDataParams struct {
	FilePath string `json:"filePath"`
}

func (c *Commander) SendPDVDataRequest(ctx context.Context, params PDVDataParams) (*nxtypes.PDVData, error) {
	var result *nxtypes.PDVData
	err := c.sendRequest(ctx, PDVDataRequestMethod, params, &result)
	return result, err
}
