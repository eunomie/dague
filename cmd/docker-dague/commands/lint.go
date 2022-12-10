package commands

import (
	"github.com/eunomie/dague"
	"github.com/eunomie/dague/daggers"

	"dagger.io/dagger"
	"github.com/spf13/cobra"
)

var (
	// LintCommands contains all commands related to linters.
	LintCommands = []*cobra.Command{
		LintGovuln(),
	}
)

// LintGovuln is a command running govulncheck against the Go source code.
func LintGovuln() *cobra.Command {
	return &cobra.Command{
		Use:   "lint:govuln",
		Short: "Lint Go code using govulncheck",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return dague.RunInDagger(ctx, func(c *dagger.Client) error {
				return daggers.GoVulnCheck(ctx, c)
			})
		},
	}
}
