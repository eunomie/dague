package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dagger.io/dagger"
	"github.com/spf13/cobra"

	"github.com/eunomie/dague"
	"github.com/eunomie/dague/config"
	"github.com/eunomie/dague/daggers"
	"github.com/eunomie/dague/types"
)

// GoCommands contains all commands related to Go like modules management or build.
func GoCommands(opts *config.Dague) []*cobra.Command {
	return []*cobra.Command{
		GoDeps(opts),
		GoMod(opts),
		GoTest(),
		GoDoc(),
		GoBuild(opts),
	}
}

// GoDeps is a command to download go modules.
func GoDeps(_ *config.Dague) *cobra.Command {
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
func GoMod(opts *config.Dague) *cobra.Command {
	return &cobra.Command{
		Use:   "go:mod MODULES...",
		Short: "Run go mod tidy and export go.mod and go.sum files",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return dague.RunInDagger(ctx, func(c *dagger.Client) error {
				if err := daggers.ExportGoMod(ctx, c); err != nil {
					return err
				}
				return nil
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

// GoBuild is a command to build a Go binary based on the local architecture.
func GoBuild(opts *config.Dague) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "go:build [OPTIONS] TARGET",
		Short: "Compile go code and export it for the local architecture",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			var targetName = args[0]
			var target config.Target
			var ok bool
			for _, t := range opts.Go.Build.Targets {
				if t.Name == targetName {
					target = t
					ok = true
					break
				}
			}
			if !ok {
				return fmt.Errorf("could not find the target %q to build", targetName)
			}

			var buildFlags []string
			if target.Ldflags != "" {
				buildFlags = append(buildFlags, "-ldflags="+os.ExpandEnv(target.Ldflags))
			}
			return dague.RunInDagger(ctx, func(c *dagger.Client) error {
				if target.Type == "local" {
					return daggers.LocalBuild(ctx, c, types.LocalBuildOpts{
						BuildOpts: types.BuildOpts{
							Dir: target.Out,
							In:  target.Path,
							EnvVars: map[string]string{
								"CGO_ENABLED": "0",
								"GO11MODULE":  "auto",
							},
							BuildFlags: buildFlags,
						},
						Out: filepath.Base(target.Path),
					})
				} else {
					var platforms []types.Platform
					for _, p := range target.Platforms {
						t := strings.SplitN(p, "/", 2)
						platforms = append(platforms, types.Platform{OS: t[0], Arch: t[1]})
					}
					return daggers.CrossBuild(ctx, c, types.CrossBuildOpts{
						BuildOpts: types.BuildOpts{
							Dir: target.Out,
							In:  target.Path,
							EnvVars: map[string]string{
								"CGO_ENABLED": "0",
								"GO11MODULE":  "auto",
							},
							BuildFlags: buildFlags,
						},
						OutFileFormat: filepath.Base(target.Path) + "_%s_%s",
						Platforms:     platforms,
					})
				}
			})
		},
	}

	return cmd
}
