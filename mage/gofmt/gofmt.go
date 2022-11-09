package gofmt

import (
	"context"

	"dagger.io/dagger"
	"github.com/eunomie/dague"
	"github.com/eunomie/dague/stages"
	"github.com/magefile/mage/mg"
)

type Gofmt mg.Namespace

// Print runs gofmt and display the recommended changes
func (Gofmt) Print(ctx context.Context) error {
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return stages.PrintGofmt(ctx, c)
	})
}

// Write runs gofmt and write the recommended changes
func (Gofmt) Write(ctx context.Context) error {
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return stages.ApplyGofmt(ctx, c)
	})
}
