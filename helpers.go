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
//	err := dague.Exec(ctx, c.Container().From("golang"), []string{"go", "build"})
func Exec(ctx context.Context, src *dagger.Container, args []string) error {
	cont := src.WithExec(args)
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
//	c.Container().From("golang").WithExec(GoInstall("golang.org/x/vuln/cmd/govulncheck@latest"))
func GoInstall(packages ...string) []string {
	return append([]string{"go", "install"}, packages...)
}

// ApkInstall runs the apk add command with the specified packaged, to install packages on alpine based systems.
// Example:
//
//	c.Container().From("alpine").WithExec(ApkInstall("build-base"))
func ApkInstall(packages ...string) []string {
	return append([]string{"apk", "add"}, packages...)
}

// AptInstall runs apt-get to install the specified packages. It updates first, install, then clean and remove apt-get lists.
// Example:
//
//	dague.AptInstall(c.Container().From("debian"), "gcc", "git")
func AptInstall(cont *dagger.Container, packages ...string) *dagger.Container {
	return cont.
		WithExec([]string{"apt-get", "update"}).
		WithExec(append([]string{"apt-get", "install", "--no-install-recommends", "-y"}, packages...)).
		WithExec([]string{"apt-get", "clean"}).
		WithExec([]string{"rm", "-rf", "/var/lib/apt/lists/*"})
}

func ExportFilePattern(ctx context.Context, cont *dagger.Container, pattern, path string) error {
	ok, err := cont.
		WithExec([]string{"sh", "-c", fmt.Sprintf("find . -name '%s' | cpio -pdm __export_dague__", pattern)}).
		Directory("./__export_dague__").
		Export(ctx, path)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("could not export %s to %s", pattern, path)
	}
	return nil
}
