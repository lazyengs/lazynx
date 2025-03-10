package commands

import "context"

const (
	SourceMapFilesToProjectsMapRequestMethod = "nx/sourceMapFilesToProjectsMap"
)

// SendSourceMapFilesToProjectsMapRequest sends a request to map source files to projects.
func (c *Commander) SendSourceMapFilesToProjectsMapRequest(ctx context.Context) (map[string][]string, error) {
	var result map[string][]string
	err := c.sendRequest(ctx, SourceMapFilesToProjectsMapRequestMethod, nil, &result)
	return result, err
}
