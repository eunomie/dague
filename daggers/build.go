package daggers

import (
	"context"
	"errors"
	"fmt"
	"path"
	"runtime"

	"dagger.io/dagger"
	"golang.org/x/sync/errgroup"

	"github.com/eunomie/dague/types"
)

func LocalBuild(ctx context.Context, c *Client, buildOpts types.LocalBuildOpts) error {
	file := path.Join(buildOpts.Dir, buildOpts.Out)
	return goBuild(ctx, c, Sources(c), runtime.GOOS, runtime.GOARCH, buildOpts.BuildOpts, file)
}

func CrossBuild(ctx context.Context, c *Client, buildOpts types.CrossBuildOpts) error {
	g, ctx := errgroup.WithContext(ctx)

	src := Sources(c)

	for _, platform := range buildOpts.Platforms {
		goos := platform.OS
		goarch := platform.Arch
		g.Go(func() error {
			file := fmt.Sprintf(path.Join(buildOpts.Dir, buildOpts.OutFileFormat), goos, goarch)
			return goBuild(ctx, c, src, goos, goarch, buildOpts.BuildOpts, file)
		})
	}
	return g.Wait()
}

func goBuild(ctx context.Context, c *Client, src *dagger.Container, os, arch string, buildOpts types.BuildOpts, buildFile string) error {
	var (
		absoluteFileInContainer = path.Join(c.Config.Go.AppDir, buildFile)
		localFile               = path.Join("./", buildFile)
	)

	cont := src.
		WithEnvVariable("GOOS", os).
		WithEnvVariable("GOARCH", arch)
	for k, v := range buildOpts.EnvVars {
		cont = cont.WithEnvVariable(k, v)
	}
	ok, err := cont.WithExec(
		append([]string{"go", "build"},
			append(buildOpts.BuildFlags, "-o", localFile, buildOpts.In)...),
	).File(absoluteFileInContainer).Export(ctx, localFile)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("could not export " + buildFile)
	}
	return nil
}
