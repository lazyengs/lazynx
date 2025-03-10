package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	GeneratorContextFromPathRequestMethod = "nx/generatorContextFromPath"
)

// GeneratorContextFromPathParams represents the parameters for the generator context from path request.
type GeneratorContextFromPathParams struct {
	Generator *nxtypes.TaskExecutionSchema `json:"generator,omitempty"` // The generator task execution schema.
	Path      string                       `json:"path"`                // The path for which to generate the context.
}

// GeneratorContextFromPathResult represents the result of the generator context from path request.
type GeneratorContextFromPathResult struct {
	Path        string `json:"path,omitempty"`        // The path for which the context was generated.
	Directory   string `json:"directory,omitempty"`   // The directory of the path.
	Project     string `json:"project,omitempty"`     // The project associated with the path.
	ProjectName string `json:"projectName,omitempty"` // The name of the project.
}

// SendGeneratorContextFromPathRequest sends a request to generate context from a path.
// It takes a context and GeneratorContextFromPathParams as parameters and returns a pointer to GeneratorContextFromPathResult and an error.
func (c *Commander) SendGeneratorContextFromPathRequest(ctx context.Context, params GeneratorContextFromPathParams) (*GeneratorContextFromPathResult, error) {
	var result *GeneratorContextFromPathResult
	err := c.sendRequest(ctx, GeneratorContextFromPathRequestMethod, params, &result)
	return result, err
}
