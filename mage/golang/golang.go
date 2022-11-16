package golang

import (
	"context"
	"path/filepath"

	"github.com/eunomie/dague"
	"github.com/eunomie/dague/daggers"
	"github.com/eunomie/dague/types"
	"github.com/magefile/mage/mg"

	"dagger.io/dagger"
)

type Go mg.Namespace

// Deps downloads go modules
func (Go) Deps(ctx context.Context) error {
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		daggers.GoDeps(c)
		return nil
	})
}

// Mod runs go mod tidy and export go.mod and go.sum files
func (Go) Mod(ctx context.Context) error {
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return daggers.ExportGoMod(ctx, c)
	})
}

// Test runs go tests
func (Go) Test(ctx context.Context) error {
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return daggers.RunGoTests(ctx, c)
	})
}

// Doc generates go documentation in README.md files
func (Go) Doc(ctx context.Context) error {
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return daggers.GoDoc(ctx, c)
	})
}

// Local compiles go code from target and export it into dist/ folder for the local architecture
func (Go) Local(ctx context.Context, target string) error {
	return Local(ctx, types.LocalBuildOpts{
		BuildOpts: types.BuildOpts{
			Dir: "dist",
			In:  target,
			EnvVars: map[string]string{
				"CGO_ENABLED": "0",
				"GO11MODULE":  "auto",
			},
		},
		Out: filepath.Base(target),
	})
}

// Cross compiles go code from target and export it into dist/ for multiple architectures (linux|darwin|windows)/(amd64|arm64)
func (Go) Cross(ctx context.Context, target string) error {
	return Cross(ctx, types.CrossBuildOpts{
		BuildOpts: types.BuildOpts{
			Dir: "dist",
			In:  target,
			EnvVars: map[string]string{
				"CGO_ENABLED": "0",
				"GO11MODULE":  "auto",
			},
		},
		OutFileFormat: filepath.Base(target) + "_%s_%s",
		Platforms: []types.Platform{
			{"linux", "amd64"},
			{"linux", "arm64"},
			{"darwin", "amd64"},
			{"darwin", "arm64"},
			{"windows", "amd64"},
			{"windows", "arm64"},
		},
	})
}

// Local compiles go code to a binary ane export it
func Local(ctx context.Context, buildOpts types.LocalBuildOpts) error {
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return daggers.LocalBuild(ctx, c, buildOpts)
	})
}

// Cross compiles a go binary for multiple platforms and export all of them
func Cross(ctx context.Context, buildOpts types.CrossBuildOpts) error {
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return daggers.CrossBuild(ctx, c, buildOpts)
	})
}
