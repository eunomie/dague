package main

import (
	"github.com/eunomie/dague/kong"
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
	kong.Run(&CLI{})
}
