package commands

import (
	"fmt"

	"github.com/eunomie/dague/internal"

	"github.com/spf13/cobra"
)

var (
	VersionCommands = []*cobra.Command{
		Version(),
	}
)

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
