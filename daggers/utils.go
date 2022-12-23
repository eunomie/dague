package daggers

import (
	"context"

	"dagger.io/dagger"

	"github.com/eunomie/dague"
	"github.com/eunomie/dague/config"
)

func RunInDagger(ctx context.Context, conf *config.Dague, do func(*Client) error) error {
	return dague.RunInDagger(ctx, func(client *dagger.Client) error {
		c := NewClient(client, conf)
		return do(c)
	})
}
