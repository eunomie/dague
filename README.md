# Dague

⚠️ In Progress!

Collection of tools to develop, build, lint, test, etc Go softwares.

Based on:
- [dagger.io](https://dagger.io)
- [mage](https://magefile.org)

## Helpers on top of Dagger

Here is some example of helpers on top of `dagger`, to help deal with Docker and Go projects.
(in the examples, `c` is a `*dagger.Client`)

* Install `apk` packages (for alpine based containers)

    ```go
   c.Container().
       From("alpine").
       Exec(dague.ApkInstall("build-base"))
   ```

* Install `go` packages

   ```go
   c.Container().
       From("golang").
       Exec(dague.GoInstall("golang.org/x/vuln/cmd/govulncheck@latest"))
   ```

* Install `apt` packages (for debian like containers). This is a bit different than with `apk` but it includes
cleaning at the end

   ```go
   dague.AptInstall(
       c.Container().From("debian"),
       "gcc", "git")
   ```

* Exec a command but return an error if any (just to avoid some boilerplate)

   ```go
   if err := dague.Exec(
       context.Background(),
       c.Container().From(baseImage),
       dagger.ContainerExecOpts{
           Args: []string{...},
       },
   }); err != nil {
       ...
   }
   ```

## Shared Tooling

The goal for this collection of tools is to make them easy to share across multiple repositories. So you can more easily bootstrap
and maintain projects as all the default actions are available from the start. Also as everything is running in containers, you
don't have to worry about not having the right version of a tool, or even the right version of the compiler.

### Examples

#### WIP Example

You can find a real, work in progress, example/experiment/PoC on `dagger` branch of [`docker/scan-cli-plugin`](https://github.com/eunomie/scan-cli-plugin/tree/dagger)

Checkout this branch, then run `./dague build:cross` to build all binaries of this Docker CLI plugin.

The output can be split in several pieces:

```shell
❯ ./dague build:cross
```

As you just checkout, you don't have the tool binary built. This script will start doing it.

```text
[+] Building 15.7s (13/13) FINISHED
 => [internal] load .dockerignore
 => => transferring context: 2B
 => [internal] load build definition from Dockerfile
 => => transferring dockerfile: 269B
 => [internal] load metadata for docker.io/eunomie/dague:0.0.1
 => [auth] eunomie/dague:pull token for registry-1.docker.io
 => [builder 1/6] FROM docker.io/eunomie/dague:0.0.1@sha256:cd534505dba93caf5e0f623147a18e562b7e4f38b35dc15772e424ada018c577
 => => resolve docker.io/eunomie/dague:0.0.1@sha256:cd534505dba93caf5e0f623147a18e562b7e4f38b35dc15772e424ada018c577
 => [internal] load build context
 => => transferring context: 2.14kB
 => CACHED [builder 2/6] COPY go.mod .
 => CACHED [builder 3/6] COPY go.sum .
 => CACHED [builder 4/6] RUN go mod download
 => [builder 5/6] COPY magefile.go .
 => [builder 6/6] RUN mage -compile /go/src/dague -goos darwin -goarch arm64
 => [stage-1 1/1] COPY --from=builder /go/src/dague /
 => exporting to client
 => => copying files 27.48MB
```

This first part will build the static binary for the tools, for the local architecture. The logs comes from a Mac M1, that's why you can see `-goos darwin -goarch arm64`.

This build is made using a multi-stage Docker image and using the `output` feature of `buildx` allowing to create files instead of images.

The result if the creation of the file `tools/dague`, static binary for the local architecture.

```text
#1 resolve image config for docker.io/library/golang:1.19.3-alpine3.16
#1 DONE 2.0s

#2 mkfile /Dockerfile
#2 DONE 0.0s

#3 mkfile /main.go
#3 CACHED

#4 [internal] load metadata for docker.io/library/golang:1.18.2-alpine
#4 DONE 1.1s

#4 [internal] load metadata for docker.io/library/golang:1.18.2-alpine
#4 DONE 1.5s

#4 [internal] load metadata for docker.io/library/golang:1.18.2-alpine
#4 DONE 1.9s

#5 mkdir /meta
#5 CACHED

#6 [build 1/6] FROM docker.io/library/golang:1.18.2-alpine@sha256:4795c5d21f01e0777707ada02408debe77fe31848be97cf9fa8a1462da78d949
#6 resolve docker.io/library/golang:1.18.2-alpine@sha256:4795c5d21f01e0777707ada02408debe77fe31848be97cf9fa8a1462da78d949 done
#6 DONE 0.0s

#7 [build 4/6] RUN go mod init go.dagger.io/dagger/shim/cmd
#7 CACHED

#8 [build 5/6] COPY . .
#8 CACHED

...

#30 copy /dist/docker-scan_linux_amd64 /docker-scan_linux_amd64
#30 DONE 0.0s

#31 copy /dist/docker-scan_darwin_amd64 /docker-scan_darwin_amd64
#31 DONE 0.1s

#32 copy /dist/docker-scan_darwin_arm64 /docker-scan_darwin_arm64
#32 DONE 0.0s

#29 exporting to client
#29 copying files 13.05MB 0.4s done
#29 copying files 13.57MB 0.3s done
#29 DONE 0.5s
```

This second step will run the `build:cross` defined in the `magefile.go`. This will launch in parallel several builds for all the different architectures.
I removed some noise from the output, but as you can see at the end the different generated files have been created and exported.

We can check that by listing the `dist` folder:

```text
❯ ls dist/
docker-scan_darwin_amd64  docker-scan_darwin_arm64  docker-scan_linux_amd64   docker-scan_linux_arm64   docker-scan_windows_amd64
```

And that's it 🎉

You just cross compiled the plugin to 5 different platforms, just by requiring Docker (and `git` in this specific case).
And a lot of what's available can be shared across multiple projects.

#### Basic Example

With the following `magefile.go`

```go
//go:build mage

package main

import (
	//mage:import
	_ "github.com/eunomie/dague/target/lint"
	//mage:import
	_ "github.com/eunomie/dague/target/golang"
)
```

This is what you get out of the box:

```text
❯ ./dague -l
Targets:
  go:deps         downloads go modules
  go:mod          runs go mod tidy and export go.mod and go.sum files
  go:test         runs go tests
  lint:all        runs all linters at once
  lint:gofmt      runs gofmt formatter and print out diff
  lint:gofumpt    runs gofumpt formatter and print out diff
  lint:govuln     checks vulnerabilities in Go code
```

#### Build Example

Now add a little bit more code in your `magefile.go`.

```go
//go:build mage

package main

import (
	"context"

	"github.com/eunomie/dague/types"
	//mage:import
	_ "github.com/eunomie/dague/target/lint"
	//mage:import
	"github.com/eunomie/dague/target/golang"

	"github.com/magefile/mage/mg"
)

type Build mg.Namespace

// Local builds local binary of myapp, for the running platform and export it to dist/
func (Build) Local(ctx context.Context) error {
	return golang.Local(ctx, types.LocalBuildOpts{
		BuildOpts: types.BuildOpts{
			Dir: "dist",
			In:  "./cmd/myapp",
			EnvVars: map[string]string{
				"CGO_ENABLED": "0",
			},
        },
		Out: "myapp",
	})
}

// Cross builds myapp binaries for all the supported platforms and export them to dist/
func (Build) Cross(ctx context.Context) error {
	return golang.Cross(ctx, types.CrossBuildOpts{
		BuildOpts: types.BuildOpts{
			Dir:           "dist",
			In:            "./cmd/myapp",
			EnvVars: map[string]string{
				"CGO_ENABLED": "0",
			},
        },
		OutFileFormat: "myapp_%s_%s",
		Platforms: []types.Platform{
			{"linux", "amd64"},
			{"linux", "arm64"},
			{"darwin", "amd64"},
			{"darwin", "arm64"},
			{"windows", "amd64"},
		},
	})
}
```

The main piece is the `type Build mg.Namespace` used to expose the two commands to the user.

On the two commands, it's using `golang` `Local` and `Cross` helpers. They are an abstraction on top of `dagger` with pre-defined behaviour
to build Go binaries and to cross compile.

With that, this is what you get:

```text
❯ ./dague -l
Targets:
  build:cross     builds myapp binaries for all the supported platforms and export them to dist/
  build:local     builds local binary of myapp, for the running platform and export it to dist/
  go:deps         downloads go modules
  go:mod          runs go mod tidy and export go.mod and go.sum files
  go:test         runs go tests
  lint:all        runs all linters at once
  lint:gofmt      runs gofmt formatter and print out diff
  lint:gofumpt    runs gofumpt formatter and print out diff
  lint:govuln     checks vulnerabilities in Go code
```

### Installation and Requirements

1. In the repository of you go project, creates a `tools` folder (name is up to you)
2. Add your `magefile.go` in the `tools` folder
3. Add `go.mod` and `go.sum` files, to make this folder independent from you main project
4. Add the following `Dockerfile` to compile this `magefile.go` to a static binary

    ```Dockerfile
    FROM eunomie/dague:0.0.1 as builder
    ARG OS
    ARG ARCH

    COPY go.mod .
    COPY go.sum .

    RUN go mod download

    COPY magefile.go .

    RUN mage -compile /go/src/dague -goos $OS -goarch $ARCH

    FROM scratch

    COPY --from=builder /go/src/dague /
    ```
5. Add this little `dague` wrapper at the root of your repository:

    ```shell
    #!/bin/sh
    
    build () {
      docker buildx build --output=type=local,dest=tools -f tools/Dockerfile --build-arg OS=`uname | tr '[:upper:]' '[:lower:]'` --build-arg ARCH=`uname -m` tools
    }
    
    if [ ! -f "tools/dague" ]; then
      build
    elif [ $# -eq 1 ] && [ "$1" = "refresh" ]; then
      build
    fi
    
    if [ $# -gt 1 ] || [ ! "$1" = "refresh" ]; then
      ./tools/dague "$@"
    fi
    ```

This wrapper will build the static tool if not exists, contains a `refresh` command if you need to rebuild it. And then it just pass everything to the tool.
So you just have to run `./dague -l` and let the magic happen.