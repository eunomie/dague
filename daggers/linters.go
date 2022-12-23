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
		dagger.ContainerExecOpts{
			Args: []string{"govulncheck", "./..."},
		},
	)
}

func GolangCILint(ctx context.Context, c *Client) error {
	return dague.Exec(
		ctx,
		sources(c, GolangCILintBase(c)),
		dagger.ContainerExecOpts{
			Args: []string{"golangci-lint", "run", "-v"},
		},
	)
}

func GolangCILintBase(c *Client) *dagger.Container {
	return c.Dagger.Container().
		From(c.Config.Go.Lint.Golangci.Image).
		WithWorkdir(c.Config.Go.AppDir)
}
