package commands

import (
	"context"

	nxtypes "github.com/lazyengs/lazynx/pkg/nxlsclient/nx-types"
)

const (
	PDVDataRequestMethod = "nx/pdvData"
)

// PDVDataParams represents the parameters for a PDV data request.
type PDVDataParams struct {
	FilePath string `json:"filePath"` // FilePath is the path to the file for which PDV data is requested.
}

// SendPDVDataRequest sends a request to get PDV data.
func (c *Commander) SendPDVDataRequest(ctx context.Context, params PDVDataParams) (*nxtypes.PDVData, error) {
	var result *nxtypes.PDVData
	err := c.sendRequest(ctx, PDVDataRequestMethod, params, &result)
	return result, err
}
