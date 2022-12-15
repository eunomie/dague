package daggers

import (
	"context"
	"errors"
	"path"
	"strings"

	"dagger.io/dagger"

	"github.com/eunomie/dague/config"
	"github.com/eunomie/dague/types"
)

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

func exportFiles(ctx context.Context, cont *dagger.Container, files []string) error {
	for _, f := range files {
		file := strings.TrimSpace(f)
		if file == "" {
			continue
		}
		ok, err := cont.File(file).Export(ctx, file)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("could not export " + f)
		}
	}
	return nil
}
