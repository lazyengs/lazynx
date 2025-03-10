package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	ProjectFolderTreeRequestMethod = "nx/projectFolderTree"
)

// ProjectFolderTreeResult represents the result of a project folder tree request.
type ProjectFolderTreeResult struct {
	SerializedTreeMap []struct {
		Dir  string           `json:"dir"`  // Dir is the directory path.
		Node nxtypes.TreeNode `json:"node"` // Node is the tree node.
	} `json:"serializedTreeMap"`
	Roots []nxtypes.TreeNode `json:"roots"` // Roots are the root nodes of the tree.
}

// SendProjectFolderTreeRequest sends a request to get the project folder tree.
func (c *Commander) SendProjectFolderTreeRequest(ctx context.Context) (*ProjectFolderTreeResult, error) {
	var result *ProjectFolderTreeResult
	err := c.sendRequest(ctx, ProjectFolderTreeRequestMethod, nil, &result)
	return result, err
}
