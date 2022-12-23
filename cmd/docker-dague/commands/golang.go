package commands

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mvdan.cc/sh/v3/syntax"

	"mvdan.cc/sh/v3/expand"

	"mvdan.cc/sh/v3/interp"

	"github.com/spf13/cobra"

	"mvdan.cc/sh/v3/shell"

	"github.com/eunomie/dague/config"
	"github.com/eunomie/dague/daggers"
	"github.com/eunomie/dague/types"
)

// GoCommands contains all commands related to Go like modules management or build.
func GoCommands(conf *config.Dague) []*cobra.Command {
	return []*cobra.Command{
		GoDeps(conf),
		GoMod(conf),
		GoTest(conf),
		GoDoc(conf),
		GoBuild(conf),
	}
}

// GoDeps is a command to download go modules.
func GoDeps(conf *config.Dague) *cobra.Command {
	return &cobra.Command{
		Use:   "go:deps",
		Short: "Download go modules",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
				daggers.GoDeps(c)
				return nil
			})
		},
	}
}

// GoMod is a command to run go mod tidy and export go.mod and go.sum files.
func GoMod(conf *config.Dague) *cobra.Command {
	return &cobra.Command{
		Use:   "go:mod MODULES...",
		Short: "Run go mod tidy and export go.mod and go.sum files",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
				if err := daggers.ExportGoMod(ctx, c); err != nil {
					return err
				}
				return nil
			})
		},
	}
}

// GoTest is a command running Go tests.
func GoTest(conf *config.Dague) *cobra.Command {
	return &cobra.Command{
		Use:   "go:test",
		Short: "Run go tests",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
				return daggers.RunGoTests(ctx, c)
			})
		},
	}
}

type goDocOptions struct {
	check bool
}

// GoDoc is a command generating Go documentation into readme.md files.
func GoDoc(conf *config.Dague) *cobra.Command {
	opts := goDocOptions{
		check: false,
	}
	cmd := &cobra.Command{
		Use:   "go:doc",
		Short: "Generate Go documentation into readme files",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
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
func GoBuild(conf *config.Dague) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "go:build [OPTIONS] TARGET",
		Short: "Compile go code and export it for the local architecture",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			targetName := args[0]
			var target config.Target
			var ok bool
			for _, t := range conf.Go.Build.Targets {
				if t.Name == targetName {
					target = t
					ok = true
					break
				}
			}
			if !ok {
				return fmt.Errorf("could not find the target %q to build", targetName)
			}

			env := map[string]string{}

			for k, v := range target.Env {
				if strings.HasPrefix(v, "$ ") {
					shellCmd := strings.TrimPrefix(v, "$ ")
					value, err := interpretShell(ctx, shellCmd, env)
					if err != nil {
						return err
					}
					env[k] = value
				} else {
					env[k] = v
				}
			}

			var buildFlags []string
			if target.Ldflags != "" {
				flags, err := shell.Expand(target.Ldflags, func(s string) string {
					return env[s]
				})
				if err != nil {
					return err
				}
				buildFlags = append(buildFlags, "-ldflags="+flags)
			}
			return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
				if target.Type == "local" {
					return daggers.LocalBuild(ctx, c, types.LocalBuildOpts{
						BuildOpts: types.BuildOpts{
							Dir:        target.Out,
							In:         target.Path,
							EnvVars:    env,
							BuildFlags: buildFlags,
						},
						Out: filepath.Base(target.Path),
					})
				}
				var platforms []types.Platform
				for _, p := range target.Platforms {
					t := strings.SplitN(p, "/", 2)
					platforms = append(platforms, types.Platform{OS: t[0], Arch: t[1]})
				}
				return daggers.CrossBuild(ctx, c, types.CrossBuildOpts{
					BuildOpts: types.BuildOpts{
						Dir:        target.Out,
						In:         target.Path,
						EnvVars:    env,
						BuildFlags: buildFlags,
					},
					OutFileFormat: filepath.Base(target.Path) + "_%s_%s",
					Platforms:     platforms,
				})
			})
		},
	}

	return cmd
}

func interpretShell(ctx context.Context, cmd string, env map[string]string) (string, error) {
	script, err := syntax.NewParser().Parse(strings.NewReader(cmd), "")
	if err != nil {
		return "", err
	}

	out := bytes.NewBufferString("")

	pairs := os.Environ()
	for k, v := range env {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
	}
	runner, err := interp.New(interp.Env(expand.ListEnviron(pairs...)), interp.StdIO(nil, out, out))
	if err != nil {
		return "", err
	}

	if err = runner.Run(ctx, script); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}
