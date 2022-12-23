package dague

import (
	"context"
	"errors"
	"fmt"
	"os"

	"dagger.io/dagger"
)

// RunInDagger initialize the dagger client and close it. In between it runs the specified function.
// Example:
//
//	dague.RunInDagger(ctx, func(c *dagger.Client) error {
//	    c.Container().From("alpine")
//	})
func RunInDagger(ctx context.Context, do func(*dagger.Client) error) error {
	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer c.Close()

	return do(c)
}

// Exec runs the specified command and check the error and exit code.
// Example:
//
//	err := dague.Exec(ctx, c.Container().From("golang"), dagger.ContainerExecOpts{
//	    Args: []string{"go", "build"},
//	})
func Exec(ctx context.Context, src *dagger.Container, opts dagger.ContainerExecOpts) error {
	cont := src.Exec(opts)
	exitCode, err := cont.ExitCode(ctx)
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return errors.New("exec failed")
	}
	return nil
}

// GoInstall installs the specified go packages.
// Example:
//
//	c.Container().From("golang").Exec(GoInstall("golang.org/x/vuln/cmd/govulncheck@latest"))
func GoInstall(packages ...string) dagger.ContainerExecOpts {
	return dagger.ContainerExecOpts{
		Args: append([]string{"go", "install"}, packages...),
	}
}

// ApkInstall runs the apk add command with the specified packaged, to install packages on alpine based systems.
// Example:
//
//	c.Container().From("alpine").Exec(ApkInstall("build-base"))
func ApkInstall(packages ...string) dagger.ContainerExecOpts {
	return dagger.ContainerExecOpts{
		Args: append([]string{"apk", "add"}, packages...),
	}
}

// AptInstall runs apt-get to install the specified packages. It updates first, install, then clean and remove apt-get lists.
// Example:
//
//	dague.AptInstall(c.Container().From("debian"), "gcc", "git")
func AptInstall(cont *dagger.Container, packages ...string) *dagger.Container {
	return cont.Exec(dagger.ContainerExecOpts{
		Args: []string{"apt-get", "update"},
	}).Exec(dagger.ContainerExecOpts{
		Args: append([]string{"apt-get", "install", "--no-install-recommends", "-y"}, packages...),
	}).Exec(dagger.ContainerExecOpts{
		Args: []string{"apt-get", "clean"},
	}).Exec(dagger.ContainerExecOpts{
		Args: []string{"rm", "-rf", "/var/lib/apt/lists/*"},
	})
}

func ExportFilePattern(ctx context.Context, cont *dagger.Container, pattern, path string) error {
	ok, err := cont.Exec(dagger.ContainerExecOpts{
		Args: []string{"sh", "-c", fmt.Sprintf("find . -name '%s' | cpio -pdm __export_dague__", pattern)},
	}).Directory("./__export_dague__").Export(ctx, path)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("could not export %s to %s", pattern, path)
	}
	return nil
}
