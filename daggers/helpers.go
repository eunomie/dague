package daggers

import (
	"context"
	"errors"
	"path"
	"strings"

	"github.com/eunomie/dague/config"

	"github.com/eunomie/dague/types"

	"dagger.io/dagger"
	"github.com/eunomie/dague"
)

func applyGoformatter(ctx context.Context, c *dagger.Client, formatter string) error {
	list, _, err := dague.ExecOut(ctx, Sources(c), dagger.ContainerExecOpts{
		Args: []string{formatter, "-l", "."},
	})
	if err != nil {
		return err
	}

	cont, err := dague.ExecCont(ctx, Sources(c), dagger.ContainerExecOpts{
		Args: []string{formatter, "-w", "."},
	})
	if err != nil {
		return err
	}

	for _, f := range strings.Split(list, "\n") {
		file := strings.TrimSpace(f)
		if file == "" {
			continue
		}
		ok, err := cont.File(f).Export(ctx, f)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("could not export " + f)
		}
	}

	return nil
}

func goBuild(ctx context.Context, src *dagger.Container, os, arch string, buildOpts types.BuildOpts, buildFile string) error {
	var (
		absoluteFileInContainer = path.Join(config.AppDir, buildFile)
		localFile               = path.Join(".", buildFile)
	)

	cont := src.
		WithEnvVariable("GOOS", os).
		WithEnvVariable("GOARCH", arch)
	for k, v := range buildOpts.EnvVars {
		cont = cont.WithEnvVariable(k, v)
	}
	ok, err := cont.Exec(dagger.ContainerExecOpts{
		Args: append([]string{"go", "build"},
			append(buildOpts.BuildFlags, "-o", localFile, buildOpts.In)...),
	}).File(absoluteFileInContainer).Export(ctx, localFile)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("could not export " + buildFile)
	}
	return nil
}
