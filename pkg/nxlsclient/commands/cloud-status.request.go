package commands

import "context"

const (
	CloudStatusRequestMethod = "nx/cloudStatus"
)

type CloudStatusResult struct {
	IsConnected bool   `json:"isConnected"`
	NxCloudUrl  string `json:"nxCloudUrl,omitempty"`
	NxCloudId   string `json:"nxCloudId,omitempty"`
}

func (c *Commander) SendCloudStatusRequest(ctx context.Context) (*CloudStatusResult, error) {
	var result *CloudStatusResult
	err := c.sendRequest(ctx, CloudStatusRequestMethod, nil, &result)
	return result, err
}
