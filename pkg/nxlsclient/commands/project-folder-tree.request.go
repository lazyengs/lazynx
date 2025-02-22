package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	ProjectFolderTreeRequestMethod = "nx/projectFolderTree"
)

type ProjectFolderTreeResult struct {
	SerializedTreeMap []struct {
		Dir  string           `json:"dir"`
		Node nxtypes.TreeNode `json:"node"`
	} `json:"serializedTreeMap"`
	Roots []nxtypes.TreeNode `json:"roots"`
}

func (c *Commander) SendProjectFolderTreeRequest(ctx context.Context) (*ProjectFolderTreeResult, error) {
	var result *ProjectFolderTreeResult
	err := c.sendRequest(ctx, ProjectFolderTreeRequestMethod, nil, &result)
	return result, err
}
