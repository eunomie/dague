package commands

import (
	"fmt"

	"dagger.io/dagger"
	"github.com/spf13/cobra"

	"github.com/eunomie/dague"
	"github.com/eunomie/dague/config"
	"github.com/eunomie/dague/daggers"
)

// LintCommands contains all commands related to linters.
func LintCommands(opts *config.Dague) []*cobra.Command {
	return []*cobra.Command{
		Lint(opts),
		LintGovuln(opts),
	}
}

// Lint is a command running all configured linters.
func Lint(opts *config.Dague) *cobra.Command {
	return &cobra.Command{
		Use:   "lint",
		Short: "Lint Go code",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return dague.RunInDagger(ctx, func(c *dagger.Client) error {
				if opts.Go.Lint.Govulncheck.Enable {
					err := daggers.GoVulnCheck(ctx, c)
					if err != nil {
						return err
					}
				}
				if opts.Go.Lint.Golangci.Enable {
					//return daggers.GoLangCILint(ctx, c)
				}
				return nil
			})
		},
	}
}

// LintGovuln is a command running govulncheck against the Go source code.
func LintGovuln(opts *config.Dague) *cobra.Command {
	return &cobra.Command{
		Use:   "lint:govuln",
		Short: "Lint Go code using govulncheck",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return dague.RunInDagger(ctx, func(c *dagger.Client) error {
				if !opts.Go.Lint.Govulncheck.Enable {
					return fmt.Errorf("govulncheck must be enabled")
				}
				return daggers.GoVulnCheck(ctx, c)
			})
		},
	}
}
