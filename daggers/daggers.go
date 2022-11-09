package daggers

import (
	"context"
	"fmt"
	"path"
	"runtime"

	"github.com/eunomie/dague/config"

	"github.com/eunomie/dague/types"
	"golang.org/x/sync/errgroup"

	"dagger.io/dagger"
	"github.com/eunomie/dague"
)

func Base(c *dagger.Client) *dagger.Container {
	return c.Container().
		From(config.BuildImage).
		Exec(dague.ApkInstall("build-base")).
		Exec(dague.GoInstall("golang.org/x/vuln/cmd/govulncheck@latest")).
		Exec(dague.GoInstall("mvdan.cc/gofumpt@latest"))
}

func GoDeps(c *dagger.Client) *dagger.Container {
	return Base(c).
		WithWorkdir(config.AppDir).
		WithMountedDirectory(config.AppDir, dague.GoModFiles(c)).
		Exec(dague.GoModDownload())
}

func Sources(c *dagger.Client) *dagger.Container {
	return GoDeps(c).
		WithMountedDirectory(config.AppDir, c.Host().Workdir())
}

func GoMod(c *dagger.Client) *dagger.Container {
	return Sources(c).
		Exec(dague.GoModTidy())
}

func ExportGoMod(ctx context.Context, c *dagger.Client) error {
	return dague.ExportGoMod(ctx, GoMod(c), config.AppDir, "./")
}

func LocalBuild(ctx context.Context, c *dagger.Client, buildOpts types.LocalBuildOpts) error {
	file := path.Join(buildOpts.Dir, buildOpts.Out)
	return goBuild(ctx, Sources(c), runtime.GOOS, runtime.GOARCH, buildOpts.BuildOpts, file)
}

func CrossBuild(ctx context.Context, c *dagger.Client, buildOpts types.CrossBuildOpts) error {
	g, ctx := errgroup.WithContext(ctx)

	src := Sources(c)

	for _, platform := range buildOpts.Platforms {
		goos := platform.OS
		goarch := platform.Arch
		g.Go(func() error {
			file := fmt.Sprintf(path.Join(buildOpts.Dir, buildOpts.OutFileFormat), goos, goarch)
			return goBuild(ctx, src, goos, goarch, buildOpts.BuildOpts, file)
		})
	}
	return g.Wait()
}

func GoVulnCheck(ctx context.Context, c *dagger.Client) error {
	return dague.Exec(ctx,
		Sources(c).WithEnvVariable("CGO_ENABLED", "0"),
		dagger.ContainerExecOpts{
			Args: []string{"govulncheck", "./..."},
		})
}

func PrintGofmt(ctx context.Context, c *dagger.Client) error {
	return dague.Exec(ctx, Sources(c), dagger.ContainerExecOpts{
		Args: []string{"gofmt", "-d", "-e", "."},
	})
}

func ApplyGofmt(ctx context.Context, c *dagger.Client) error {
	return applyGoformatter(ctx, c, "gofmt")
}

func PrintGofumpt(ctx context.Context, c *dagger.Client) error {
	return dague.Exec(ctx, Sources(c), dagger.ContainerExecOpts{
		Args: []string{"gofumpt", "-d", "-e", "."},
	})
}

func ApplyGofumpt(ctx context.Context, c *dagger.Client) error {
	return applyGoformatter(ctx, c, "gofumpt")
}

func RunGoTests(ctx context.Context, c *dagger.Client) error {
	return dague.Exec(ctx, Sources(c), dagger.ContainerExecOpts{
		Args: []string{"go", "test", "-race", "-cover", "-shuffle=on", "./..."},
	})
}
