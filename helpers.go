package dague

import (
	"context"
	"errors"
	"fmt"
	"os"

	"dagger.io/dagger"
)

var (
	goModFiles = []string{"go.mod", "go.sum"}
)

// RunInDagger initialize the dagger client and close it. In between run the specified function.
func RunInDagger(ctx context.Context, do func(*dagger.Client) error) error {
	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer c.Close()

	return do(c)
}

// Exec runs the specified command and check the error and exit code.
func Exec(ctx context.Context, cont *dagger.Container, opts dagger.ContainerExecOpts) error {
	_, err := ExecCont(ctx, cont, opts)
	return err
}

// ExecCont runs the specified command and check the error and exist code. Returns the container and the error if exists.
func ExecCont(ctx context.Context, src *dagger.Container, opts dagger.ContainerExecOpts) (*dagger.Container, error) {
	cont := src.Exec(opts)
	exitCode, err := cont.ExitCode(ctx)
	if err != nil {
		return nil, err
	}
	if exitCode != 0 {
		return nil, errors.New("exec failed")
	}
	return cont, nil
}

// ExecOut runs the specified command and return the content of stdout and stderr.
func ExecOut(ctx context.Context, src *dagger.Container, opts dagger.ContainerExecOpts) (string, string, error) {
	cont, err := ExecCont(ctx, src, opts)
	if err != nil {
		return "", "", err
	}

	stdout, err := cont.Stdout().Contents(ctx)
	if err != nil {
		return "", "", err
	}

	stderr, err := cont.Stderr().Contents(ctx)
	if err != nil {
		return "", "", err
	}

	return stdout, stderr, nil
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

// GoModFiles creates a directory containing the default go mod files.
func GoModFiles(c *dagger.Client) *dagger.Directory {
	src := c.Host().Workdir()
	goMods := c.Directory()
	for _, f := range goModFiles {
		goMods = goMods.WithFile(f, src.File(f))
	}
	return goMods
}

// GoModDownload runs the go mod download command.
func GoModDownload() dagger.ContainerExecOpts {
	return dagger.ContainerExecOpts{
		Args: []string{"go", "mod", "download"},
	}
}

// GoModTidy runs the go mod tidy command.
func GoModTidy() dagger.ContainerExecOpts {
	return dagger.ContainerExecOpts{
		Args: []string{"go", "mod", "tidy", "-v"},
	}
}

// ExportGoMod reads the default go mod tiles from the specified internal dir and export them to the host.
func ExportGoMod(ctx context.Context, cont *dagger.Container, contDir, exportDir string) error {
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
