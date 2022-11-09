package gofumpt

import (
	"context"

	"dagger.io/dagger"
	"github.com/eunomie/dague"
	"github.com/eunomie/dague/stages"
	"github.com/magefile/mage/mg"
)

type Gofumpt mg.Namespace

// Print runs gofumpt and print the recommended changes
func (Gofumpt) Print(ctx context.Context) error {
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return stages.PrintGofumpt(ctx, c)
	})
}

// Write runs gofumpt and write the recommended changes
func (Gofumpt) Write(ctx context.Context) error {
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return stages.ApplyGofumpt(ctx, c)
	})
}