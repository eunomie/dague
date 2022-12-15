package commands

import (
	"dagger.io/dagger"
	"github.com/spf13/cobra"

	"github.com/eunomie/dague"
	"github.com/eunomie/dague/config"
	"github.com/eunomie/dague/daggers"
)

// FmtCommands contains all commands related to code formatting
func FmtCommands(opts *config.Dague) []*cobra.Command {
	return []*cobra.Command{
		FmtPrint(opts),
		FmtWrite(opts),
		GoImports(opts),
		Fmt(opts),
	}
}

func Fmt(opts *config.Dague) *cobra.Command {
	return &cobra.Command{
		Use:   "fmt",
		Short: "Format files and imports",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return dague.RunInDagger(ctx, func(c *dagger.Client) error {
				return daggers.ApplyFormatAndImports(ctx, c, opts.Go.Fmt.Formatter, opts.Go.Fmt.Goimports.Locals)
			})
		},
	}
}

// FmtPrint is a command to print the result of Go formatter. Files will not be modified.
func FmtPrint(opts *config.Dague) *cobra.Command {
	return &cobra.Command{
		Use:   "fmt:print",
		Short: "Print result of go formatter",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return dague.RunInDagger(ctx, func(c *dagger.Client) error {
				return daggers.PrintGoformatter(ctx, c, opts.Go.Fmt.Formatter)
			})
		},
	}
}

// FmtWrite is a command to write to the existing files the result of the Go formatter.
func FmtWrite(opts *config.Dague) *cobra.Command {
	return &cobra.Command{
		Use:   "fmt:write",
		Short: "Write result of go formatter to existing files",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return dague.RunInDagger(ctx, func(c *dagger.Client) error {
				return daggers.ApplyGoformatter(ctx, c, opts.Go.Fmt.Formatter)
			})
		},
	}
}

func GoImports(opts *config.Dague) *cobra.Command {
	return &cobra.Command{
		Use:   "fmt:imports",
		Short: "Reorder imports",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return dague.RunInDagger(ctx, func(c *dagger.Client) error {
				return daggers.GoImports(ctx, c, opts.Go.Fmt.Goimports.Locals)
			})
		},
	}
}
