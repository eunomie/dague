package config

import (
	"fmt"
	"os"

	"github.com/ghodss/yaml"
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
	Target struct {
		Name      string            `yaml:"name"`
		Type      string            `yaml:"type"`
		Path      string            `yaml:"path"`
		Out       string            `yaml:"out"`
		Env       map[string]string `yaml:"env"`
		Ldflags   string            `yaml:"ldflags"`
		Platforms []string          `yaml:"platforms,omitempty"`
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

// go:embed .dague.default.yml
var defaults []byte

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
