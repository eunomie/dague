package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"

	"github.com/eunomie/dague/config"
	"github.com/eunomie/dague/internal"
	"github.com/eunomie/dague/internal/commands"
)

const (
	PluginName = "dague"
)

func pluginMain() {
	plugin.Run(func(dockerCli command.Cli) *cobra.Command {
		var conf config.Dague
		c := &cobra.Command{
			Short:            "Docker Dague",
			Use:              PluginName,
			TraverseChildren: true,
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				c, err := config.Load(cmd.Context())
				if err != nil {
					return err
				}
				conf = c
				return nil
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) == 0 {
					return cmd.Help()
				}
				_ = cmd.Help()
				return cli.StatusError{
					StatusCode: 1,
					Status:     fmt.Sprintf("unknown docker command: %q", PluginName+" "+args[0]),
				}
			},
		}
		originalPreRun := c.PersistentPreRunE
		c.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			if err := plugin.PersistentPreRunE(cmd, args); err != nil {
				return err
			}
			if originalPreRun != nil {
				if err := originalPreRun(cmd, args); err != nil {
					return err
				}
			}
			return nil
		}
		l := commands.NewList()

		c.AddCommand(
			&cobra.Command{
				Use:   "version",
				Short: "Print version",
				Args:  cobra.NoArgs,
				Run: func(_ *cobra.Command, _ []string) {
					fmt.Println("version:", internal.Version)
					fmt.Println("git commit:", internal.Commit)
				},
			},
			&cobra.Command{
				Use:   "go:lint",
				Short: "Lint Go code (--help for subcommands)",
				Long: `Subcommands:
  go:lint:govuln   Lint Go code using govulncheck
  go:lint:golangci Lint Go code using golangci`,
				Args: cobra.NoArgs,
				RunE: func(cmd *cobra.Command, args []string) error {
					return l.Run(cmd.Context(), "go:lint", args, &conf, nil)
				},
			},
			&cobra.Command{
				Use:    "go:lint:govuln",
				Hidden: true,
				Short:  "Lint Go code using govulncheck",
				Args:   cobra.NoArgs,
				RunE: func(cmd *cobra.Command, args []string) error {
					return l.Run(cmd.Context(), "go:lint:govuln", args, &conf, nil)
				},
			},
			&cobra.Command{
				Use:    "go:lint:golangci",
				Hidden: true,
				Short:  "Lint Go code using golangci",
				Args:   cobra.NoArgs,
				RunE: func(cmd *cobra.Command, args []string) error {
					return l.Run(cmd.Context(), "go:lint:golangci", args, &conf, nil)
				},
			},
			func() *cobra.Command {
				type goFmtOptions struct {
					check bool
				}

				opts := goFmtOptions{
					check: false,
				}
				cmd := &cobra.Command{
					Use:   "go:fmt",
					Short: "Format files and imports (--help for subcommands)",
					Long: `Subcommands:
  go:fmt:print       Print result of configured Go formatter
  go:fmt:write       Write result of configured Go formatter
  go:imports:print   Print result of go imports
  go:imports:write   Reorder imports using configured locals`,
					Args: cobra.NoArgs,
					RunE: func(cmd *cobra.Command, args []string) error {
						return l.Run(cmd.Context(), "go:fmt", args, &conf, map[string]interface{}{
							"check": opts.check,
						})
					},
				}

				flags := cmd.Flags()
				flags.BoolVar(&opts.check, "check", false, "check the format is up-to-date")

				return cmd
			}(),
			&cobra.Command{
				Use:    "go:fmt:print",
				Hidden: true,
				Short:  "Print result of go formatter",
				Args:   cobra.NoArgs,
				RunE: func(cmd *cobra.Command, args []string) error {
					return l.Run(cmd.Context(), "go:fmt:print", args, &conf, nil)
				},
			},
			&cobra.Command{
				Use:    "go:fmt:write",
				Hidden: true,
				Short:  "Write result of go formatter to existing files",
				Args:   cobra.NoArgs,
				RunE: func(cmd *cobra.Command, args []string) error {
					return l.Run(cmd.Context(), "go:fmt:write", args, &conf, nil)
				},
			},
			&cobra.Command{
				Use:    "go:imports:write",
				Hidden: true,
				Short:  "Reorder imports",
				Args:   cobra.NoArgs,
				RunE: func(cmd *cobra.Command, args []string) error {
					return l.Run(cmd.Context(), "go:imports:write", args, &conf, nil)
				},
			},
			&cobra.Command{
				Use:    "go:imports:print",
				Hidden: true,
				Short:  "Print result of goimports",
				Args:   cobra.NoArgs,
				RunE: func(cmd *cobra.Command, args []string) error {
					return l.Run(cmd.Context(), "go:imports:print", args, &conf, nil)
				},
			},
			&cobra.Command{
				Use:   "go:mod",
				Short: "Run go mod download and go mod tidy (--help for subcommands)",
				Long: `Subcommands:
  go:mod:download  Download go modules`,
				Args: cobra.NoArgs,
				RunE: func(cmd *cobra.Command, args []string) error {
					return l.Run(cmd.Context(), "go:mod", args, &conf, nil)
				},
			},
			&cobra.Command{
				Use:    "go:mode:download",
				Hidden: true,
				Short:  "Download go modules",
				Args:   cobra.NoArgs,
				RunE: func(cmd *cobra.Command, args []string) error {
					return l.Run(cmd.Context(), "go:mod:download", args, &conf, nil)
				},
			},
			&cobra.Command{
				Use:   "go:test",
				Short: "Run go tests",
				Args:  cobra.NoArgs,
				RunE: func(cmd *cobra.Command, args []string) error {
					return l.Run(cmd.Context(), "go:test", args, &conf, nil)
				},
			},

			func() *cobra.Command {
				type goDocOptions struct {
					check bool
				}

				opts := goDocOptions{
					check: false,
				}
				cmd := &cobra.Command{
					Use:   "go:doc",
					Short: "Generate Go documentation into readme files",
					Args:  cobra.NoArgs,
					RunE: func(cmd *cobra.Command, args []string) error {
						return l.Run(cmd.Context(), "go:doc", args, &conf, map[string]interface{}{
							"check": opts.check,
						})
					},
				}

				flags := cmd.Flags()
				flags.BoolVar(&opts.check, "check", false, "check the documentation is up-to-date")

				return cmd
			}(),

			&cobra.Command{
				Use:   "go:build [TARGET]",
				Short: "Compile go code",
				Args:  cobra.MaximumNArgs(1),
				RunE: func(cmd *cobra.Command, args []string) error {
					return l.Run(cmd.Context(), "go:build", args, &conf, nil)
				},
			},

			&cobra.Command{
				Use:   "go:exec [TASK]",
				Short: "Execute scripts inside the build container",
				Args:  cobra.MaximumNArgs(1),
				RunE: func(cmd *cobra.Command, args []string) error {
					return l.Run(cmd.Context(), "go:exec", args, &conf, nil)
				},
			},

			&cobra.Command{
				Use:   "task [TASK]",
				Short: "Run tasks",
				Args:  cobra.MaximumNArgs(1),
				RunE: func(cmd *cobra.Command, args []string) error {
					return l.Run(cmd.Context(), "task", args, &conf, nil)
				},
			},
		)
		return c
	}, manager.Metadata{
		SchemaVersion: "0.1.0",
		Vendor:        "Docker Inc.",
		Version:       internal.Version,
	})
}

func main() {
	if plugin.RunningStandalone() {
		os.Args = append([]string{"docker", "dague"}, os.Args[1:]...)
	}
	pluginMain()
}
