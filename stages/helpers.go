package stages

import (
	"context"
	"errors"
	"fmt"

	"dagger.io/dagger"
)

func exec(ctx context.Context, cont *dagger.Container, opts dagger.ContainerExecOpts) error {
	exitCode, err := cont.Exec(opts).ExitCode(ctx)
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return errors.New("exec failed")
	}
	return nil
}

func goInstall(packages ...string) dagger.ContainerExecOpts {
	return dagger.ContainerExecOpts{
		Args: append([]string{"go", "install"}, packages...),
	}
}

func apkAdd(packages ...string) dagger.ContainerExecOpts {
	return dagger.ContainerExecOpts{
		Args: append([]string{"apk", "add"}, packages...),
	}
}

func goModDownload() dagger.ContainerExecOpts {
	return dagger.ContainerExecOpts{
		Args: []string{"go", "mod", "download"},
	}
}

func goModTidy() dagger.ContainerExecOpts {
	return dagger.ContainerExecOpts{
		Args: []string{"go", "mod", "tidy", "-v"},
	}
}

func exportGoMod(ctx context.Context, cont *dagger.Container, contDir, exportDir string) error {
	for _, f := range goModFiles {
		file := fmt.Sprintf("%s/%s", contDir, f)
		ok, err := cont.File(file).Export(ctx, exportDir+f)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("could not export " + file)
		}
	}
	return nil
}
