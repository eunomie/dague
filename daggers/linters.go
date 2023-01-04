package daggers

import (
	"context"

	"dagger.io/dagger"

	"github.com/eunomie/dague"
)

func GoVulnCheck(ctx context.Context, c *Client) error {
	return dague.Exec(
		ctx,
		Sources(c).WithEnvVariable("CGO_ENABLED", "0"),
		[]string{"govulncheck", "./..."},
	)
}

func GolangCILint(ctx context.Context, c *Client) error {
	return dague.Exec(
		ctx,
		sources(c, GolangCILintBase(c)),
		[]string{"golangci-lint", "run", "-v", "--timeout", "5m"},
	)
}

func GolangCILintBase(c *Client) *dagger.Container {
	base := c.Dagger.Container().
		From(c.Config.Go.Lint.Golangci.Image).
		WithWorkdir(c.Config.Go.AppDir)

	base = applyBase(base, c.Dagger, c.Config)

	return base
}
