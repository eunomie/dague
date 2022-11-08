package stages

import (
	"context"
	"errors"
	"fmt"
	"runtime"

	"github.com/eunomie/dague"

	"github.com/eunomie/dague/types"

	"golang.org/x/sync/errgroup"

	"dagger.io/dagger"
)

var (
	BuildImage = "golang:1.19.3-alpine3.16"
	AppDir     = "/go/src"
)

func Base(c *dagger.Client) *dagger.Container {
	return c.Container().
		From(BuildImage).
		Exec(dague.ApkInstall("build-base")).
		Exec(dague.GoInstall("golang.org/x/vuln/cmd/govulncheck@latest")).
		Exec(dague.GoInstall("mvdan.cc/gofumpt@latest"))
}

func GoDeps(c *dagger.Client) *dagger.Container {
	return Base(c).
		WithWorkdir(AppDir).
		WithMountedDirectory(AppDir, dague.GoModFiles(c)).
		Exec(dague.GoModDownload())
}

func Sources(c *dagger.Client) *dagger.Container {
	return GoDeps(c).
		WithMountedDirectory(AppDir, c.Host().Workdir())
}

func GoMod(c *dagger.Client) *dagger.Container {
	return Sources(c).
		Exec(dague.GoModTidy())
}

func ExportGoMod(ctx context.Context, c *dagger.Client) error {
	return dague.ExportGoMod(ctx, GoMod(c), AppDir, "./")
}

func LocalBuild(ctx context.Context, c *dagger.Client, buildOpts types.LocalBuildOpts) error {
	file := "/" + buildOpts.Dir + "/" + buildOpts.Out
	absFile := AppDir + file
	localFile := "." + file
	cont := Sources(c).
		WithEnvVariable("GOOS", runtime.GOOS).
		WithEnvVariable("GOARCH", runtime.GOARCH).
		WithEnvVariable("GO111MODULE", "auto")
	for k, v := range buildOpts.EnvVars {
		cont = cont.WithEnvVariable(k, v)
	}
	ok, err := cont.Exec(dagger.ContainerExecOpts{
		Args: append([]string{"go", "build"},
			append(buildOpts.BuildFlags, "-o", localFile, buildOpts.In)...),
	}).File(absFile).Export(ctx, localFile)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("could not export " + file)
	}
	return nil
}

func CrossBuild(ctx context.Context, c *dagger.Client, buildOpts types.CrossBuildOpts) error {
	g, ctx := errgroup.WithContext(ctx)

	src := Sources(c)

	for _, platform := range buildOpts.Platforms {
		goos := platform.OS
		goarch := platform.Arch
		g.Go(func() error {
			file := fmt.Sprintf("/"+buildOpts.Dir+"/"+buildOpts.OutFileFormat, goos, goarch)
			absFile := AppDir + file
			localFile := "." + file
			cont := src.
				WithEnvVariable("GOOS", goos).
				WithEnvVariable("GOARCH", goarch).
				WithEnvVariable("GO111MODULE", "auto")
			for k, v := range buildOpts.EnvVars {
				cont = cont.WithEnvVariable(k, v)
			}
			ok, err := cont.Exec(dagger.ContainerExecOpts{
				Args: append([]string{"go", "build"},
					append(buildOpts.BuildFlags, "-o", localFile, buildOpts.In)...),
			}).File(absFile).Export(ctx, localFile)
			if err != nil {
				return err
			}
			if !ok {
				return errors.New("could not export " + file)
			}
			return nil
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

func Gofumpt(ctx context.Context, c *dagger.Client) error {
	return dague.Exec(ctx, Sources(c), dagger.ContainerExecOpts{
		Args: []string{"gofumpt", "-d", "-e", "."},
	})
}

func Gofmt(ctx context.Context, c *dagger.Client) error {
	return dague.Exec(ctx, Sources(c), dagger.ContainerExecOpts{
		Args: []string{"gofmt", "-d", "-e", "."},
	})
}

func RunGoTests(ctx context.Context, c *dagger.Client) error {
	return dague.Exec(ctx, Sources(c), dagger.ContainerExecOpts{
		Args: []string{"go", "test", "-race", "-cover", "-shuffle=on", "./..."},
	})
}
