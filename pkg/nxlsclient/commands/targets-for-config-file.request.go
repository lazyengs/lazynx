package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	TargetsForConfigFileRequestMethod = "nx/targetsForConfigFile"
)

type TargetsForConfigFileParams struct {
	ProjectName    string `json:"projectName"`
	ConfigFilePath string `json:"configFilePath"`
}

func (c *Commander) SendTargetsForConfigFileRequest(ctx context.Context, params TargetsForConfigFileParams) (map[string]nxtypes.TargetConfiguration, error) {
	var result map[string]nxtypes.TargetConfiguration
	err := c.sendRequest(ctx, TargetsForConfigFileRequestMethod, params, &result)
	return result, err
}
