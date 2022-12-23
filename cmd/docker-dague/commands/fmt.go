package commands

import (
	"github.com/spf13/cobra"

	"github.com/eunomie/dague/config"
	"github.com/eunomie/dague/daggers"
)

// FmtCommands contains all commands related to code formatting
func FmtCommands(conf *config.Dague) []*cobra.Command {
	return []*cobra.Command{
		FmtPrint(conf),
		FmtWrite(conf),
		GoImports(conf),
		Fmt(conf),
	}
}

func Fmt(conf *config.Dague) *cobra.Command {
	return &cobra.Command{
		Use:   "fmt",
		Short: "Format files and imports",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
				return daggers.ApplyFormatAndImports(ctx, c, conf.Go.Fmt.Formatter, conf.Go.Fmt.Goimports.Locals)
			})
		},
	}
}

// FmtPrint is a command to print the result of Go formatter. Files will not be modified.
func FmtPrint(conf *config.Dague) *cobra.Command {
	return &cobra.Command{
		Use:   "fmt:print",
		Short: "Print result of go formatter",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
				return daggers.PrintGoformatter(ctx, c, conf.Go.Fmt.Formatter)
			})
		},
	}
}

// FmtWrite is a command to write to the existing files the result of the Go formatter.
func FmtWrite(conf *config.Dague) *cobra.Command {
	return &cobra.Command{
		Use:   "fmt:write",
		Short: "Write result of go formatter to existing files",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
				return daggers.ApplyGoformatter(ctx, c, conf.Go.Fmt.Formatter)
			})
		},
	}
}

func GoImports(conf *config.Dague) *cobra.Command {
	return &cobra.Command{
		Use:   "fmt:imports",
		Short: "Reorder imports",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
				return daggers.GoImports(ctx, c, conf.Go.Fmt.Goimports.Locals)
			})
		},
	}
}
