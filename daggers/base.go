package daggers

import (
	"dagger.io/dagger"

	"github.com/eunomie/dague"
)

// GoBase is a default container based on a Golang build image (see config.BuildImage) on top of which is installed several
// packages and Go packages.
// The workdir is also set based on config.AppDir.
//
// This container is used as the root of many other commands, allowing to share cache as much as possible.
func GoBase(c *Client) *dagger.Container {
	base := c.Dagger.Container().
		From(c.Config.Go.Image.Src).
		WithExec(dague.ApkInstall("build-base", "git")).
		WithExec(dague.GoInstall("golang.org/x/vuln/cmd/govulncheck@latest")).
		WithExec(dague.GoInstall("golang.org/x/tools/cmd/goimports@latest")).
		WithExec(dague.GoInstall("github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest"))

	base = applyBase(base, c.Dagger, c.Config)

	if len(c.Config.Go.Image.ApkPackages) > 0 {
		base = base.WithExec(dague.ApkInstall(c.Config.Go.Image.ApkPackages...))
	}
	if len(c.Config.Go.Image.AptPackages) > 0 {
		base = dague.AptInstall(base, c.Config.Go.Image.AptPackages...)
	}
	if len(c.Config.Go.Image.GoPackages) > 0 {
		base = base.WithExec(dague.GoInstall(c.Config.Go.Image.GoPackages...))
	}

	return base.WithWorkdir(c.Config.Go.AppDir)
}

// GoDeps mount the Go module files and download the needed dependencies.
func GoDeps(c *Client) *dagger.Container {
	return GoBase(c).
		WithMountedDirectory(c.Config.Go.AppDir, goModFiles(c)).
		WithExec(goModDownload())
}

func sources(c *Client, cont *dagger.Container) *dagger.Container {
	return cont.WithMountedDirectory(c.Config.Go.AppDir, c.Dagger.Host().Directory("."))
}

// Sources is a container based on GoDeps. It contains the Go source code but also all the needed dependencies from
// Go modules.
func Sources(c *Client) *dagger.Container {
	return sources(c, GoDeps(c))
}

// SourcesNoDeps is a container including all the source code, but without the Go modules downloaded.
// It can be helpful with projects where dependencies are vendored but also just minimise the number of steps when
// it's not required.
func SourcesNoDeps(c *Client) *dagger.Container {
	return sources(c, GoBase(c))
}
