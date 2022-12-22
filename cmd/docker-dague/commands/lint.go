package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/eunomie/dague/config"
	"github.com/eunomie/dague/daggers"
)

// LintCommands contains all commands related to linters.
func LintCommands(conf *config.Dague) []*cobra.Command {
	return []*cobra.Command{
		Lint(conf),
		LintGovuln(conf),
		LintGolangCILint(conf),
	}
}

// Lint is a command running all configured linters.
func Lint(conf *config.Dague) *cobra.Command {
	return &cobra.Command{
		Use:   "lint",
		Short: "Lint Go code",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
				if conf.Go.Lint.Govulncheck.Enable {
					err := daggers.GoVulnCheck(ctx, c)
					if err != nil {
						return err
					}
				}
				if conf.Go.Lint.Golangci.Enable {
					return daggers.GolangCILint(ctx, c)
				}
				return nil
			})
		},
	}
}

// LintGovuln is a command running govulncheck against the Go source code.
func LintGovuln(conf *config.Dague) *cobra.Command {
	return &cobra.Command{
		Use:   "lint:govuln",
		Short: "Lint Go code using govulncheck",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
				if !conf.Go.Lint.Govulncheck.Enable {
					return fmt.Errorf("govulncheck must be enabled")
				}
				return daggers.GoVulnCheck(ctx, c)
			})
		},
	}
}

func LintGolangCILint(conf *config.Dague) *cobra.Command {
	return &cobra.Command{
		Use:   "lint:golangci",
		Short: "Lint Go code",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
				if !conf.Go.Lint.Golangci.Enable {
					return fmt.Errorf("golangci-lint must be enabled")
				}
				return daggers.GolangCILint(ctx, c)
			})
		},
	}
}
