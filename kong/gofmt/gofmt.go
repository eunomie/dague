package gofmt

import (
	"context"

	"dagger.io/dagger"
	"github.com/eunomie/dague"
	"github.com/eunomie/dague/daggers"
)

type (
	// Gofmt regroups commands based on gofmt
	Gofmt struct {
		Print PrintCmd `cmd:"" group:"gofmt" help:"print result of gofmt"`
		Write WriteCmd `cmd:"" group:"gofmt" help:"write result of gofmt to existing files"`
	}

	// PrintCmd runs gofmt and print out the result
	PrintCmd struct{}
	// WriteCmd runs gofmt and write the result to the files
	WriteCmd struct{}
)

func (PrintCmd) Run() error {
	ctx := context.Background()
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return daggers.PrintGofmt(ctx, c)
	})
}

func (WriteCmd) Run() error {
	ctx := context.Background()
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return daggers.ApplyGofmt(ctx, c)
	})
}
