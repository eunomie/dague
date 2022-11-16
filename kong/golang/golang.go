package golang

import (
	"context"
	"path/filepath"
	"strings"

	"dagger.io/dagger"
	"github.com/eunomie/dague"
	"github.com/eunomie/dague/daggers"
	"github.com/eunomie/dague/types"
)

type (
	Golang struct {
		Deps  DepsCmd       `cmd:"" group:"golang" help:"download go modules"`
		Mod   ModCmd        `cmd:"" group:"golang" help:"run go mod tidy and export go.mod and go.sum files"`
		Test  TestCmd       `cmd:"" group:"golang" help:"run go tests"`
		Doc   DocCmd        `cmd:"" group:"golang" help:"generate to documentation in README.md files"`
		Build BuildCmd      `cmd:"" group:"golang" help:"compiles go code and export it for the local architecture"`
		Cross CrossBuildCmd `cmd:"" group:"golang" help:"compiles go code and export it for multiple architectures"`
	}

	DepsCmd  struct{}
	ModCmd   struct{}
	TestCmd  struct{}
	DocCmd   struct{}
	BuildCmd struct {
		In  string `arg:"" help:"directory to build"`
		Out string `default:"dist" help:"directory where to export the binary"`
	}
	CrossBuildCmd struct {
		In        string   `arg:"" help:"directory to build"`
		Out       string   `default:"dist" help:"directory where to export the binary"`
		Platforms []string `default:"linux/amd64,linux/arm64,darwin/amd64,darwin/arm64,windows/amd64,windows/arm64" help:"platforms to build the binary"`
	}
)

func (DepsCmd) Run() error {
	ctx := context.Background()
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		daggers.GoDeps(c)
		return nil
	})
}

func (ModCmd) Run() error {
	ctx := context.Background()
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return daggers.ExportGoMod(ctx, c)
	})
}

func (TestCmd) Run() error {
	ctx := context.Background()
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return daggers.RunGoTests(ctx, c)
	})
}

func (DocCmd) Run() error {
	ctx := context.Background()
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return daggers.GoDoc(ctx, c)
	})
}

func (cmd BuildCmd) Run() error {
	ctx := context.Background()
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return daggers.LocalBuild(ctx, c, types.LocalBuildOpts{
			BuildOpts: types.BuildOpts{
				Dir: cmd.Out,
				In:  cmd.In,
				EnvVars: map[string]string{
					"CGO_ENABLED": "0",
					"GO11MODULE":  "auto",
				},
			},
			Out: filepath.Base(cmd.In),
		})
	})
}

func (cmd CrossBuildCmd) Run() error {
	ctx := context.Background()
	var platforms []types.Platform
	for _, p := range cmd.Platforms {
		t := strings.SplitN(p, "/", 2)
		platforms = append(platforms, types.Platform{t[0], t[1]})
	}
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return daggers.CrossBuild(ctx, c, types.CrossBuildOpts{
			BuildOpts: types.BuildOpts{
				Dir: cmd.Out,
				In:  cmd.In,
				EnvVars: map[string]string{
					"CGO_ENABLED": "0",
					"GO11MODULE":  "auto",
				},
			},
			OutFileFormat: filepath.Base(cmd.In) + "_%s_%s",
			Platforms:     platforms,
		})
	})
}
