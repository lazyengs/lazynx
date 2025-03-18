package commands

import (
	"context"

	nxtypes "github.com/lazyengs/lazynx/pkg/nxlsclient/nx-types"
)

const (
	TargetsForConfigFileRequestMethod = "nx/targetsForConfigFile"
)

// TargetsForConfigFileParams represents the parameters for a targets for config file request.
type TargetsForConfigFileParams struct {
	ProjectName    string `json:"projectName"`    // ProjectName is the name of the project.
	ConfigFilePath string `json:"configFilePath"` // ConfigFilePath is the path to the configuration file.
}

// SendTargetsForConfigFileRequest sends a request to get targets for a configuration file.
func (c *Commander) SendTargetsForConfigFileRequest(ctx context.Context, params TargetsForConfigFileParams) (map[string]nxtypes.TargetConfiguration, error) {
	var result map[string]nxtypes.TargetConfiguration
	err := c.sendRequest(ctx, TargetsForConfigFileRequestMethod, params, &result)
	return result, err
}
