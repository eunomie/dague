package golang

import (
	"context"

	"github.com/eunomie/dague"
	"github.com/eunomie/dague/stages"
	"github.com/eunomie/dague/types"
	"github.com/magefile/mage/mg"

	"dagger.io/dagger"
)

type Go mg.Namespace

// Deps downloads go modules
func (Go) Deps(ctx context.Context) error {
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		stages.GoDeps(c)
		return nil
	})
}

// Mod runs go mod tidy and export go.mod and go.sum files
func (Go) Mod(ctx context.Context) error {
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return stages.ExportGoMod(ctx, c)
	})
}

// Test runs go tests
func (Go) Test(ctx context.Context) error {
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return stages.RunGoTests(ctx, c)
	})
}

// Local compile go code to a binary ane export it
func Local(ctx context.Context, buildOpts types.LocalBuildOpts) error {
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return stages.LocalBuild(ctx, c, buildOpts)
	})
}

// Cross cross compiles a go binary and export all of them
func Cross(ctx context.Context, buildOpts types.CrossBuildOpts) error {
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return stages.CrossBuild(ctx, c, buildOpts)
	})
}