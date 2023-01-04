package daggers

import (
	"context"

	"github.com/eunomie/dague"
)

func RunGoTests(ctx context.Context, c *Client) error {
	return dague.Exec(
		ctx,
		Sources(c),
		[]string{"go", "test", "-race", "-cover", "-shuffle=on", "-v", "./..."},
	)
}
