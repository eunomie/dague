# Dague

Build Go projects, better. Based on [`dagger`](https://dagger.io).

## Why?

`dague` is a `docker` cli plugin. It acts as a opinionated Go toolchain only relying on Docker as dependency.

You don't need to have the right version of Go or any other dependencies, if you have Docker you have everything.



Especially when you're working on multiple projects, there's always the question of the tooling, the different versions
of all the tools, do you have all the needed requirements, etc.

## Usage

Just drop `docker-dague` binary in your `~/.docker/cli-plugin` directory and you're setup.

```
❯ docker dague --help

Usage:  docker dague COMMAND

Docker Dague

Commands:
  fmt:print   Print result of gofumpt
  fmt:write   Write result of gofumpt to existing files
  go:build    Compile go code and export it for the local architecture
  go:cross    Compile go code and export it for multiple architectures
  go:deps     Download go modules
  go:doc      Generate Go documentation into readme files
  go:mod      Run go mod tidy and export go.mod and go.sum files
  go:test     Run go tests
  lint:govuln Lint Go code using govulncheck
  version     Print version

Run 'docker dague COMMAND --help' for more information on a command.
```

<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# dague

```go
import "github.com/eunomie/dague"
```

## Index

- [Variables](<#variables>)
- [func ApkInstall(packages ...string) dagger.ContainerExecOpts](<#func-apkinstall>)
- [func AptInstall(cont *dagger.Container, packages ...string) *dagger.Container](<#func-aptinstall>)
- [func Exec(ctx context.Context, cont *dagger.Container, opts dagger.ContainerExecOpts) error](<#func-exec>)
- [func ExecCont(ctx context.Context, src *dagger.Container, opts dagger.ContainerExecOpts) (*dagger.Container, error)](<#func-execcont>)
- [func ExecOut(ctx context.Context, src *dagger.Container, opts dagger.ContainerExecOpts) (string, string, error)](<#func-execout>)
- [func ExportFilePattern(ctx context.Context, cont *dagger.Container, pattern, path string) error](<#func-exportfilepattern>)
- [func ExportGoMod(ctx context.Context, cont *dagger.Container, contDir, exportDir string) error](<#func-exportgomod>)
- [func GoInstall(packages ...string) dagger.ContainerExecOpts](<#func-goinstall>)
- [func GoModDownload() dagger.ContainerExecOpts](<#func-gomoddownload>)
- [func GoModFiles(c *dagger.Client) *dagger.Directory](<#func-gomodfiles>)
- [func GoModTidy() dagger.ContainerExecOpts](<#func-gomodtidy>)
- [func RunInDagger(ctx context.Context, do func(*dagger.Client) error) error](<#func-runindagger>)


## Variables

```go
var goModFiles = []string{"go.mod", "go.sum"}
```

## func ApkInstall

```go
func ApkInstall(packages ...string) dagger.ContainerExecOpts
```

ApkInstall runs the apk add command with the specified packaged, to install packages on alpine based systems. Example:

```
c.Container().From("alpine").Exec(ApkInstall("build-base"))
```

## func AptInstall

```go
func AptInstall(cont *dagger.Container, packages ...string) *dagger.Container
```

AptInstall runs apt\-get to install the specified packages. It updates first, install, then clean and remove apt\-get lists. Example:

```
dague.AptInstall(c.Container().From("debian"), "gcc", "git")
```

## func Exec

```go
func Exec(ctx context.Context, cont *dagger.Container, opts dagger.ContainerExecOpts) error
```

Exec runs the specified command and check the error and exit code. Example:

```
err := dague.Exec(ctx, c.Container().From("golang"), dagger.ContainerExecOpts{
    Args: []string{"go", "build"},
})
```

## func ExecCont

```go
func ExecCont(ctx context.Context, src *dagger.Container, opts dagger.ContainerExecOpts) (*dagger.Container, error)
```

ExecCont runs the specified command and check the error and exist code. Returns the container and the error if exists. Example:

```
cont, err := dague.ExecCont(ctx, c.Container().From("golang"), dagger.ContainerExecOpts{
    Args: []string{"go", "build"},
})
```

## func ExecOut

```go
func ExecOut(ctx context.Context, src *dagger.Container, opts dagger.ContainerExecOpts) (string, string, error)
```

ExecOut runs the specified command and return the content of stdout and stderr. Example:

```
stdout, stderr, err := dague.ExecOut(ctx, c.Container().From("golang"), dagger.ContainerExecOpts{
    Args: []string{"go", "build"},
})
```

## func ExportFilePattern

```go
func ExportFilePattern(ctx context.Context, cont *dagger.Container, pattern, path string) error
```

## func ExportGoMod

```go
func ExportGoMod(ctx context.Context, cont *dagger.Container, contDir, exportDir string) error
```

ExportGoMod reads the default go mod tiles from the specified internal dir and export them to the host.

## func GoInstall

```go
func GoInstall(packages ...string) dagger.ContainerExecOpts
```

GoInstall installs the specified go packages. Example:

```
c.Container().From("golang").Exec(GoInstall("golang.org/x/vuln/cmd/govulncheck@latest"))
```

## func GoModDownload

```go
func GoModDownload() dagger.ContainerExecOpts
```

GoModDownload runs the go mod download command.

## func GoModFiles

```go
func GoModFiles(c *dagger.Client) *dagger.Directory
```

GoModFiles creates a directory containing the default go mod files.

## func GoModTidy

```go
func GoModTidy() dagger.ContainerExecOpts
```

GoModTidy runs the go mod tidy command.

## func RunInDagger

```go
func RunInDagger(ctx context.Context, do func(*dagger.Client) error) error
```

RunInDagger initialize the dagger client and close it. In between it runs the specified function. Example:

```
dague.RunInDagger(ctx, func(c *dagger.Client) error {
    c.Container().From("alpine")
})
```



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)


<!-- gomarkdoc:embed:end -->