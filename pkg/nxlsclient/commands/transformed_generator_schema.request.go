package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	TransformedGeneratorSchemaRequestMethod = "nx/transformedGeneratorSchema"
)

func (c *Commander) SendTransformedGeneratorSchemaRequest(ctx context.Context, schema nxtypes.GeneratorSchema) (*nxtypes.GeneratorSchema, error) {
	var result *nxtypes.GeneratorSchema
	err := c.sendRequest(ctx, TransformedGeneratorSchemaRequestMethod, schema, &result)
	return result, err
}
