package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/eunomie/dague/config"
	"github.com/eunomie/dague/internal"
)

func VersionCommands(_ *config.Dague) []*cobra.Command {
	return []*cobra.Command{
		Version(),
	}
}

func Version() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println("version:", internal.Version)
			fmt.Println("git commit:", internal.Commit)
		},
	}
}
