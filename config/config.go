package config

import (
	"fmt"
	"os"

	"github.com/ghodss/yaml"
)

var (
	// BuildImage is the base Golang docker image used to run all the tools.
	BuildImage = "golang:1.19.4-alpine3.17"
	// AppDir is the default local to the container folder where to copy/mount sources.
	AppDir = "/go/src"
)

type (
	Dague struct {
		Go Go `yaml:"go"`
	}
	Goimports struct {
		Locals []string `yaml:"locals"`
	}
	Fmt struct {
		Formatter string    `yaml:"formatter"`
		Goimports Goimports `yaml:"goimports"`
	}
	Govulncheck struct {
		Enable bool `yaml:"enable"`
	}
	Golangci struct {
		Enable bool   `yaml:"enable"`
		Image  string `yaml:"image"`
	}
	Lint struct {
		Govulncheck Govulncheck `yaml:"govulncheck"`
		Golangci    Golangci    `yaml:"golangci"`
	}
	Env struct {
		CGOENABLED int `yaml:"CGO_ENABLED"`
	}
	Target struct {
		Name      string   `yaml:"name"`
		Type      string   `yaml:"type"`
		Path      string   `yaml:"path"`
		Out       string   `yaml:"out"`
		Env       Env      `yaml:"env"`
		Ldflags   string   `yaml:"ldflags"`
		Platforms []string `yaml:"platforms,omitempty"`
	}
	Build struct {
		Targets []Target `yaml:"targets"`
	}
	Go struct {
		Image  string `yaml:"image"`
		AppDir string `yaml:"appDir"`
		Fmt    Fmt    `yaml:"fmt"`
		Lint   Lint   `yaml:"lint"`
		Build  Build  `yaml:"build"`
	}
)

const (
	defaultConfigFile = ".dague.yml"
)

var (
	// go:embed .dague.default.yml
	defaults []byte
)

func Load() (Dague, error) {
	var dague Dague
	err := yaml.Unmarshal(defaults, &dague)
	if err != nil {
		return Dague{}, fmt.Errorf("could not read default configuration: %w", err)
	}

	configData, err := os.ReadFile(defaultConfigFile)
	if err != nil {
		return Dague{}, fmt.Errorf("could not read .dague.yml config file: %w", err)
	}

	err = yaml.Unmarshal(configData, &dague)
	if err != nil {
		return Dague{}, fmt.Errorf("could not parse .dague.yml config file: %w", err)
	}

	return dague, nil
}
