package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/eunomie/dague/cmd/docker-dague/commands"
	"github.com/eunomie/dague/config"
	"github.com/eunomie/dague/internal"
)

const (
	PluginName = "dague"
)

func pluginMain() {
	plugin.Run(func(dockerCli command.Cli) *cobra.Command {
		var opts config.Dague
		c := &cobra.Command{
			Short:            "Docker Dague",
			Use:              PluginName,
			TraverseChildren: true,
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				c, err := config.Load()
				if err != nil {
					return err
				}
				opts = c
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
		c.AddCommand(commands.VersionCommands(&opts)...)
		c.AddCommand(commands.LintCommands(&opts)...)
		c.AddCommand(commands.FmtCommands(&opts)...)
		c.AddCommand(commands.GoCommands(&opts)...)
		return c
	}, manager.Metadata{
		SchemaVersion: "0.1.0",
		Vendor:        "Docker Inc.",
		Version:       internal.Version,
	})
}

func main() {
	if plugin.RunningStandalone() {
		os.Args = append([]string{"docker"}, os.Args[1:]...)
	}
	pluginMain()
}
