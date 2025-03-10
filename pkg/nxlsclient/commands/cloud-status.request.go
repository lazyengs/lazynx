package commands

import "context"

const (
	CloudStatusRequestMethod = "nx/cloudStatus"
)

// CloudStatusResult represents the result of a cloud status request.
type CloudStatusResult struct {
	IsConnected bool   `json:"isConnected"`
	NxCloudUrl  string `json:"nxCloudUrl,omitempty"`
	NxCloudId   string `json:"nxCloudId,omitempty"`
}

// SendCloudStatusRequest sends a request to retrieve the cloud status.
// It takes a context as a parameter and returns a pointer to CloudStatusResult and an error.
func (c *Commander) SendCloudStatusRequest(ctx context.Context) (*CloudStatusResult, error) {
	var result *CloudStatusResult
	err := c.sendRequest(ctx, CloudStatusRequestMethod, nil, &result)
	return result, err
}
