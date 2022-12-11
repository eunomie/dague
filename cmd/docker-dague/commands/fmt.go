package commands

import (
	"github.com/eunomie/dague"
	"github.com/eunomie/dague/daggers"

	"dagger.io/dagger"
	"github.com/spf13/cobra"
)

var (
	// FmtCommands contains all commands related to code formatting
	FmtCommands = []*cobra.Command{
		FmtPrint(),
		FmtWrite(),
	}
)

// FmtPrint is a command to print the result of Go formatter. Files will not be modified.
func FmtPrint() *cobra.Command {
	return &cobra.Command{
		Use:   "fmt:print",
		Short: "Print result of gofumpt",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return dague.RunInDagger(ctx, func(c *dagger.Client) error {
				return daggers.PrintGofumpt(ctx, c)
			})
		},
	}
}

// FmtWrite is a command to write to the existing files the result of the Go formatter.
func FmtWrite() *cobra.Command {
	return &cobra.Command{
		Use:   "fmt:write",
		Short: "Write result of gofumpt to existing files",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			return dague.RunInDagger(ctx, func(c *dagger.Client) error {
				return daggers.ApplyGofumpt(ctx, c)
			})
		},
	}
}
