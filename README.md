# Dague

Build Go projects, better. Based on [`dagger`](https://dagger.io).

## Why?

`dague` is a `docker` cli plugin. It acts as a opinionated Go toolchain only relying on Docker as dependency.

You don't need to have the right version of Go or any other dependencies, if you have Docker you have everything.

## Installation

1. Download the binary corresponding to your platform from the [`latest` release](https://github.com/eunomie/dague/releases/latest).

    binaries are available for linux, mac and windows, for amd64 and arm64

2. Rename the binary to `docker-dague` and make it executable

    ```
    chmod +x docker-dague
    ```

3. On Mac, authorize the binary (as not signed):

    ```
    xattr -d com.apple.quarantine docker-dague
    ```

4. Copy it to the docker directory for CLI plugins:

    ```
    mkdir -p ~/.docker/cli-plugins
    install docker-dague ~/.docker/cli-plugins/ 
    ```

## Usage

Create a `.dague.yml` file to configure your build targets.

If you want to build a binary from `main/path` to `dist` this is the minimal file you need:

```yaml
go:
  build:
    targets:
      - name: local-build
        type: local # to allow local build, can be 'cross' to enable cross platform build
        path: ./main/path
        out: ./dist
```

With that, you can run `docker dague go:build local-build` and it will build your binary and put it under `./dist/`.

The build is performed inside containers, so you don't have to worry about the needed dependencies, tools, versions, etc.

By default `dague` comes with handy go tools already configured like:

- `go:fmt`: runs `goimports` and `gofumpt` to re-format the code
- `go:lint`: runs `golangci-lint` and `govulncheck`
- `go:doc`: generate Go documentation in markdown inside README.me files
- `go:test`: run go unit tests with handy defaults (`-race -cover -shuffle=on`)
- `go:mod`: run `go mod tidy` and update `go.mod` and `go.sum` files

It's also possible to define any script that will be run from the inside of the build container.
The exec task can also define files to export to the host.

```yaml
go:
  exec:
    info:
      cmds: |
        uname -a > info.txt
        go version >> info.txt
      export:
        pattern: info.txt
        path: .
```

Then you can run:

```text
❯ docker dague go:exec info
# ...

❯ cat info.txt
Linux buildkitsandbox 5.15.49-linuxkit #1 SMP PREEMPT Tue Sep 13 07:51:32 UTC 2022 aarch64 Linux
go version go1.19.4 linux/arm64
```

In comparison, this is the output of `go version` directly on my host:

```text
❯ go version
go version go1.19.4 darwin/arm64
```

You can also define any arbitrary task to be run using `go:task`:

```yaml
tasks:
  install:
    deps:
      - go:build local
    cmds: |
      mkdir -p ~/.docker/cli-plugins
      install dist/docker-dague ~/.docker/cli-plugins/docker-dague
```

The command `docker dague task install` will first run `go:build local` then run the shell script to install the binary.
The shell script is run using a Go shell implementation so is portable across platforms.

To know more about the possibilities and available configuration, please refer to [the configuration reference file](./.dague.reference.yml).

<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# dague

```go
import "github.com/eunomie/dague"
```

## Index

- [func ApkInstall(packages ...string) dagger.ContainerExecOpts](<#func-apkinstall>)
- [func AptInstall(cont *dagger.Container, packages ...string) *dagger.Container](<#func-aptinstall>)
- [func Exec(ctx context.Context, src *dagger.Container, opts dagger.ContainerExecOpts) error](<#func-exec>)
- [func ExportFilePattern(ctx context.Context, cont *dagger.Container, pattern, path string) error](<#func-exportfilepattern>)
- [func GoInstall(packages ...string) dagger.ContainerExecOpts](<#func-goinstall>)
- [func RunInDagger(ctx context.Context, do func(*dagger.Client) error) error](<#func-runindagger>)


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
func Exec(ctx context.Context, src *dagger.Container, opts dagger.ContainerExecOpts) error
```

Exec runs the specified command and check the error and exit code. Example:

```
err := dague.Exec(ctx, c.Container().From("golang"), dagger.ContainerExecOpts{
    Args: []string{"go", "build"},
})
```

## func ExportFilePattern

```go
func ExportFilePattern(ctx context.Context, cont *dagger.Container, pattern, path string) error
```

## func GoInstall

```go
func GoInstall(packages ...string) dagger.ContainerExecOpts
```

GoInstall installs the specified go packages. Example:

```
c.Container().From("golang").Exec(GoInstall("golang.org/x/vuln/cmd/govulncheck@latest"))
```

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