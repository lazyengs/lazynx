package commands

import (
	"context"

	nxtypes "github.com/lazyengs/lazynx/pkg/nxlsclient/nx-types"
)

const (
	TransformedGeneratorSchemaRequestMethod = "nx/transformedGeneratorSchema"
)

// SendTransformedGeneratorSchemaRequest sends a request to get the transformed generator schema.
func (c *Commander) SendTransformedGeneratorSchemaRequest(ctx context.Context, schema nxtypes.GeneratorSchema) (*nxtypes.GeneratorSchema, error) {
	var result *nxtypes.GeneratorSchema
	err := c.sendRequest(ctx, TransformedGeneratorSchemaRequestMethod, schema, &result)
	return result, err
}
