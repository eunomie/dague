package daggers

import (
	"context"
	"fmt"
	"path"
	"runtime"
	"strings"

	"dagger.io/dagger"
	"golang.org/x/sync/errgroup"

	"github.com/eunomie/dague"
	"github.com/eunomie/dague/config"
	"github.com/eunomie/dague/types"
)

// Base is a default container based on a Golang build image (see config.BuildImage) on top of which is installed several
// packages and Go packages.
// The workdir is also set based on config.AppDir.
//
// This container is used as the root of many other commands, allowing to share cache as much as possible.
func Base(c *dagger.Client) *dagger.Container {
	return c.Container().
		From(config.BuildImage).
		Exec(dague.ApkInstall("build-base", "git")).
		Exec(dague.GoInstall("golang.org/x/vuln/cmd/govulncheck@latest")).
		Exec(dague.GoInstall("golang.org/x/tools/cmd/goimports@latest")).
		Exec(dague.GoInstall("mvdan.cc/gofumpt@latest")).
		Exec(dague.GoInstall("github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest")).
		WithWorkdir(config.AppDir)
}

// GoDeps mount the Go module files and download the needed dependencies.
func GoDeps(c *dagger.Client) *dagger.Container {
	return Base(c).
		WithMountedDirectory(config.AppDir, dague.GoModFiles(c)).
		Exec(dague.GoModDownload())
}

func sources(c *dagger.Client, cont *dagger.Container) *dagger.Container {
	return cont.WithMountedDirectory(config.AppDir, c.Host().Workdir())
}

// Sources is a container based on GoDeps. It contains the Go source code but also all the needed dependencies from
// Go modules.
func Sources(c *dagger.Client) *dagger.Container {
	return sources(c, GoDeps(c))
}

// SourcesNoDeps is a container including all the source code, but without the Go modules downloaded.
// It can be helpful with projects where dependencies are vendored but also just minimise the number of steps when
// it's not required.
func SourcesNoDeps(c *dagger.Client) *dagger.Container {
	return sources(c, Base(c))
}

func GoMod(c *dagger.Client) *dagger.Container {
	return Sources(c).
		WithWorkdir(config.AppDir).
		Exec(dague.GoModTidy())
}

func ExportGoMod(ctx context.Context, c *dagger.Client) error {
	return dague.ExportFilePattern(ctx, GoMod(c), "go.*", "./")
}

func LocalBuild(ctx context.Context, c *dagger.Client, buildOpts types.LocalBuildOpts) error {
	file := path.Join(buildOpts.Dir, buildOpts.Out)
	return goBuild(ctx, Sources(c).WithWorkdir(buildOpts.Dir), runtime.GOOS, runtime.GOARCH, buildOpts.BuildOpts, file)
}

func CrossBuild(ctx context.Context, c *dagger.Client, buildOpts types.CrossBuildOpts) error {
	g, ctx := errgroup.WithContext(ctx)

	src := Sources(c).WithWorkdir(buildOpts.Dir)

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

func PrintGoformatter(ctx context.Context, c *dagger.Client, formatter string) error {
	return dague.Exec(ctx, SourcesNoDeps(c), formatPrint(formatter))
}

func ApplyGoformatter(ctx context.Context, c *dagger.Client, formatter string) error {
	return dague.ExportFilePattern(
		ctx,
		Sources(c).Exec(formatWrite(formatter)),
		"*.go",
		"./",
	)
}

func ApplyFormatAndImports(ctx context.Context, c *dagger.Client, formatter string, locals []string) error {
	return dague.ExportFilePattern(
		ctx,
		Sources(c).Exec(goImports(locals)).Exec(formatWrite(formatter)),
		"*.go",
		"./",
	)
}

func goImports(locals []string) dagger.ContainerExecOpts {
	args := []string{"goimports", "-w", "-format-only"}
	if len(locals) > 0 {
		args = append(args, "-local", strings.Join(locals, ","))
	}
	args = append(args, ".")
	return dagger.ContainerExecOpts{Args: args}
}

func formatWrite(formatter string) dagger.ContainerExecOpts {
	return dagger.ContainerExecOpts{
		Args: []string{formatter, "-w", "."},
	}
}

func formatPrint(formatter string) dagger.ContainerExecOpts {
	return dagger.ContainerExecOpts{
		Args: []string{formatter, "-d", "-e", "."},
	}
}

func GoImports(ctx context.Context, c *dagger.Client, locals []string) error {
	return dague.ExportFilePattern(
		ctx,
		SourcesNoDeps(c).Exec(goImports(locals)),
		"*.go",
		"./")
}

func RunGoTests(ctx context.Context, c *dagger.Client) error {
	return dague.Exec(ctx, Sources(c), dagger.ContainerExecOpts{
		Args: []string{"go", "test", "-race", "-cover", "-shuffle=on", "./..."},
	})
}

func GoDoc(ctx context.Context, c *dagger.Client) error {
	return dague.ExportFilePattern(
		ctx,
		SourcesNoDeps(c).Exec(dagger.ContainerExecOpts{
			Args: []string{"gomarkdoc", "-u", "-e", "-o", "{{.Dir}}/README.md", "./..."},
		}),
		"*.md",
		".",
	)
}

func CheckGoDoc(ctx context.Context, c *dagger.Client) error {
	return dague.Exec(ctx, SourcesNoDeps(c), dagger.ContainerExecOpts{
		Args: []string{"gomarkdoc", "-c", "-u", "-e", "-o", "{{.Dir}}/README.md", "./..."},
	})
}
