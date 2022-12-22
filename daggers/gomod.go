package daggers

import (
	"context"

	"dagger.io/dagger"

	"github.com/eunomie/dague"
)

var goModDefaulFiles = []string{"go.mod", "go.sum"}

func GoMod(c *Client) *dagger.Container {
	return Sources(c).
		WithWorkdir(c.Config.Go.AppDir).
		Exec(goModTidy())
}

func ExportGoMod(ctx context.Context, c *Client) error {
	return dague.ExportFilePattern(
		ctx,
		GoMod(c),
		"go.*",
		"./",
	)
}

// GoModFiles creates a directory containing the default go mod files.
func goModFiles(c *Client) *dagger.Directory {
	src := c.Dagger.Host().Workdir()
	goMods := c.Dagger.Directory()
	for _, f := range goModDefaulFiles {
		goMods = goMods.WithFile(f, src.File(f))
	}
	return goMods
}

// GoModDownload runs the go mod download command.
func goModDownload() dagger.ContainerExecOpts {
	return dagger.ContainerExecOpts{
		Args: []string{"go", "mod", "download"},
	}
}

// GoModTidy runs the go mod tidy command.
func goModTidy() dagger.ContainerExecOpts {
	return dagger.ContainerExecOpts{
		Args: []string{"go", "mod", "tidy", "-v"},
	}
}
