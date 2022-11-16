package gofumpt

import (
	"context"

	"dagger.io/dagger"
	"github.com/eunomie/dague"
	"github.com/eunomie/dague/daggers"
)

type (
	// Gofumpt regroups commands based on gofumpt
	Gofumpt struct {
		Print PrintCmd `cmd:"" group:"gofumpt" help:"print result of gofumpt"`
		Write WriteCmd `cmd:"" group:"gofumpt" help:"write result of gofumpt to existing files"`
	}

	// PrintCmd runs gofumpt and print out the result
	PrintCmd struct{}
	// WriteCmd runs gofumpt and write the result to the files
	WriteCmd struct{}
)

func (PrintCmd) Run() error {
	ctx := context.Background()
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return daggers.PrintGofumpt(ctx, c)
	})
}

func (WriteCmd) Run() error {
	ctx := context.Background()
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return daggers.ApplyGofumpt(ctx, c)
	})
}
