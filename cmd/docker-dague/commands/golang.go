package commands

import (
	"path/filepath"
	"strings"

	"github.com/eunomie/dague"
	"github.com/eunomie/dague/daggers"
	"github.com/eunomie/dague/types"

	"dagger.io/dagger"
	"github.com/spf13/cobra"
)

var (
	// GoCommands contains all commands related to Go like modules management or build.
	GoCommands = []*cobra.Command{
		GoDeps(),
		GoMod(),
		GoTest(),
		GoDoc(),
		GoBuild(),
		GoCrossBuild(),
	}
)

// GoDeps is a command to download go modules.
func GoDeps() *cobra.Command {
	return &cobra.Command{
		Use:   "go:deps",
		Short: "Download go modules",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return dague.RunInDagger(ctx, func(c *dagger.Client) error {
				daggers.GoDeps(c)
				return nil
			})
		},
	}
}

// GoMod is a command to run go mod tidy and export go.mod and go.sum files.
func GoMod() *cobra.Command {
	return &cobra.Command{
		Use:   "go:mod",
		Short: "Run go mod tidy and export go.mod and go.sum files",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return dague.RunInDagger(ctx, func(c *dagger.Client) error {
				return daggers.ExportGoMod(ctx, c)
			})
		},
	}
}

// GoTest is a command running Go tests.
func GoTest() *cobra.Command {
	return &cobra.Command{
		Use:   "go:test",
		Short: "Run go tests",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return dague.RunInDagger(ctx, func(c *dagger.Client) error {
				return daggers.RunGoTests(ctx, c)
			})
		},
	}
}

type goDocOptions struct {
	check bool
}

// GoDoc is a command generating Go documentation into readme.md files.
func GoDoc() *cobra.Command {
	opts := goDocOptions{
		check: false,
	}
	cmd := &cobra.Command{
		Use:   "go:doc",
		Short: "Generate Go documentation into readme files",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return dague.RunInDagger(ctx, func(c *dagger.Client) error {
				if opts.check {
					return daggers.CheckGoDoc(ctx, c)
				}
				return daggers.GoDoc(ctx, c)
			})
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&opts.check, "check", false, "check the documentation is up-to-date")

	return cmd
}

type goBuildOptions struct {
	out     string
	ldflags string
}

// GoBuild is a command to build a Go binary based on the local architecture.
func GoBuild() *cobra.Command {
	opts := goBuildOptions{}

	cmd := &cobra.Command{
		Use:   "go:build [OPTIONS] DIRECTORY",
		Short: "Compile go code and export it for the local architecture",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			var buildFlags []string
			if opts.ldflags != "" {
				buildFlags = append(buildFlags, "-ldflags="+opts.ldflags)
			}
			return dague.RunInDagger(ctx, func(c *dagger.Client) error {
				return daggers.LocalBuild(ctx, c, types.LocalBuildOpts{
					BuildOpts: types.BuildOpts{
						Dir: opts.out,
						In:  args[0],
						EnvVars: map[string]string{
							"CGO_ENABLED": "0",
							"GO11MODULE":  "auto",
						},
						BuildFlags: buildFlags,
					},
					Out: filepath.Base(args[0]),
				})
			})
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.out, "out", "o", "dist", "directory where to export the binary")
	flags.StringVar(&opts.ldflags, "ldflags", "", "arguments to pass on each go tool link invocation")

	return cmd
}

type goCrossBuildOptions struct {
	out       string
	platforms []string
	ldflags   string
}

// GoCrossBuild is a command to perform cross compilation and generate Go binaries for multiple platforms.
func GoCrossBuild() *cobra.Command {
	opts := goCrossBuildOptions{}

	cmd := &cobra.Command{
		Use:   "go:cross [OPTIONS] DIRECTORY",
		Short: "Compile go code and export it for multiple architectures",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			var platforms []types.Platform
			for _, p := range opts.platforms {
				t := strings.SplitN(p, "/", 2)
				platforms = append(platforms, types.Platform{OS: t[0], Arch: t[1]})
			}
			var buildFlags []string
			if opts.ldflags != "" {
				buildFlags = append(buildFlags, "-ldflags="+opts.ldflags)
			}
			return dague.RunInDagger(ctx, func(c *dagger.Client) error {
				return daggers.CrossBuild(ctx, c, types.CrossBuildOpts{
					BuildOpts: types.BuildOpts{
						Dir: opts.out,
						In:  args[0],
						EnvVars: map[string]string{
							"CGO_ENABLED": "0",
							"GO11MODULE":  "auto",
						},
						BuildFlags: buildFlags,
					},
					OutFileFormat: filepath.Base(args[0]) + "_%s_%s",
					Platforms:     platforms,
				})
			})
		},
	}

	defaultPlatforms := []string{"linux/amd64", "linux/arm64", "darwin/amd64", "darwin/arm64", "windows/amd64", "windows/arm64"}
	flags := cmd.Flags()
	flags.StringVarP(&opts.out, "out", "o", "dist", "directory where to export the binary")
	flags.StringArrayVarP(&opts.platforms, "platform", "p", defaultPlatforms, "platform to build the binary")
	flags.StringVar(&opts.ldflags, "ldflags", "", "arguments to pass on each go tool link invocation")

	return cmd
}
