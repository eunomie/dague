package main

import (
	"fmt"
	"os"

	"github.com/eunomie/dague/cmd/docker-dague/commands"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/eunomie/dague/internal"
	"github.com/spf13/cobra"
)

const (
	PluginName = "dague"
)

func pluginMain() {
	plugin.Run(func(dockerCli command.Cli) *cobra.Command {
		c := &cobra.Command{
			Short:            "Docker Dague",
			Use:              PluginName,
			TraverseChildren: true,
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
		c.AddCommand(commands.VersionCommands...)
		c.AddCommand(commands.LintCommands...)
		c.AddCommand(commands.FmtCommands...)
		c.AddCommand(commands.GoCommands...)
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
