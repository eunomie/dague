package dague

import (
	"context"
	"os"

	"dagger.io/dagger"
)

func RunInDagger(ctx context.Context, do func(*dagger.Client) error) error {
	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer c.Close()

	return do(c)
}
