package daggers

import (
	"context"

	"dagger.io/dagger"

	"github.com/eunomie/dague"
)

func RunGoTests(ctx context.Context, c *Client) error {
	return dague.Exec(
		ctx,
		Sources(c),
		dagger.ContainerExecOpts{
			Args: []string{"go", "test", "-race", "-cover", "-shuffle=on", "-v", "./..."},
		},
	)
}
