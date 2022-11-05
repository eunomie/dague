package lint

import (
	"context"

	"dagger.io/dagger"
	"github.com/eunomie/dague"
	"github.com/eunomie/dague/stages"

	"github.com/magefile/mage/mg"
)

type Lint mg.Namespace

// All runs all linters at once
func (t Lint) All(_ context.Context) error {
	mg.Deps(
		t.Gofumpt,
		t.Govuln,
	)
	return nil
}

// Gofumpt runs gofumpt formatter and print out diff
func (Lint) Gofumpt(ctx context.Context) error {
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return stages.Gofumpt(ctx, c)
	})
}

// Gofmt runs gofmt formatter and print out diff
func (Lint) Gofmt(ctx context.Context) error {
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return stages.Gofmt(ctx, c)
	})
}

// Govuln checks vulnerabilities in Go code
func (Lint) Govuln(ctx context.Context) error {
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return stages.GoVulnCheck(ctx, c)
	})
}
