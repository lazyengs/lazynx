package commands

import "context"

const (
	SourceMapFilesToProjectsMapRequestMethod = "nx/sourceMapFilesToProjectsMap"
)

func (c *Commander) SendSourceMapFilesToProjectsMapRequest(ctx context.Context) (map[string][]string, error) {
	var result map[string][]string
	err := c.sendRequest(ctx, SourceMapFilesToProjectsMapRequestMethod, nil, &result)
	return result, err
}
