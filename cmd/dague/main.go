package main

import (
	"github.com/alecthomas/kong"
	"github.com/eunomie/dague/kong/gofumpt"
	"github.com/eunomie/dague/kong/golang"
	"github.com/eunomie/dague/kong/lint"
)

type (
	CLI struct {
		Lint    lint.Lint       `cmd:""`
		Gofumpt gofumpt.Gofumpt `cmd:""`
		Go      golang.Golang   `cmd:""`
	}
)

func main() {
	var cli CLI
	ctx := kong.Parse(&cli, kong.UsageOnError())
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
