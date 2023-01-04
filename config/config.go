package config

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/eunomie/dague/internal/shell"

	"gopkg.in/yaml.v2"
)

type (
	Dague struct {
		Vars  map[string]string `yaml:"vars"`
		Go    Go                `yaml:"go"`
		Tasks Tasks             `yaml:"tasks"`
	}

	Go struct {
		Image  Image           `yaml:"image"`
		AppDir string          `yaml:"appDir"`
		Fmt    Fmt             `yaml:"fmt"`
		Lint   Lint            `yaml:"lint"`
		Build  Build           `yaml:"build"`
		Exec   map[string]Exec `yaml:"exec"`
	}

	Image struct {
		Src         string            `yaml:"src"`
		AptPackages []string          `yaml:"aptPackages"`
		ApkPackages []string          `yaml:"apkPackages"`
		GoPackages  []string          `yaml:"goPackages"`
		Mounts      map[string]string `yaml:"mounts"`
		Env         map[string]string `yaml:"env"`
		Caches      []Cache           `yaml:"caches"`
	}

	Cache struct {
		Target string `yaml:"target"`
	}

	Fmt struct {
		Formatter string    `yaml:"formatter"`
		Goimports Goimports `yaml:"goimports"`
	}

	Goimports struct {
		Locals []string `yaml:"locals"`
	}

	Lint struct {
		Govulncheck Govulncheck `yaml:"govulncheck"`
		Golangci    Golangci    `yaml:"golangci"`
	}

	Govulncheck struct {
		Enable bool `yaml:"enable"`
	}

	Golangci struct {
		Enable bool   `yaml:"enable"`
		Image  string `yaml:"image"`
	}

	Build struct {
		Targets map[string]Target `yaml:"targets"`
	}

	Target struct {
		Path      string            `yaml:"path"`
		Out       string            `yaml:"out"`
		Env       map[string]string `yaml:"env"`
		Ldflags   string            `yaml:"ldflags"`
		Platforms []string          `yaml:"platforms,omitempty"`
	}

	Exec struct {
		Deps   []string `yaml:"deps"`
		Cmds   string   `yaml:"cmds"`
		Export Export   `yaml:"export"`
	}

	Export struct {
		Pattern string `yaml:"pattern"`
		Path    string `yaml:"path"`
	}

	Tasks map[string]Task

	Task struct {
		Deps []string `yaml:"deps"`
		Cmds string   `yaml:"cmds"`
	}
)

const (
	defaultConfigFile = ".dague.yml"
)

//go:embed .dague.default.yml
var defaults []byte

func Load(ctx context.Context) (Dague, error) {
	configData, err := os.ReadFile(defaultConfigFile)
	if err != nil {
		return Dague{}, fmt.Errorf("could not read .dague.yml config file: %w", err)
	}

	merged, err := YAML([][]byte{defaults, configData}, false)
	if err != nil {
		return Dague{}, fmt.Errorf("could not merge .dague.yml with defaults: %w", err)
	}

	var dague Dague
	err = yaml.Unmarshal(merged.Bytes(), &dague)
	if err != nil {
		return Dague{}, fmt.Errorf("could not parse .dague.yml config file: %w", err)
	}

	// post step to expand vars
	vars := map[string]string{}
	var res string
	for k, v := range dague.Vars {
		if strings.HasPrefix(v, "shell ") {
			res, err = shell.Interpret(ctx, strings.TrimPrefix(v, "shell"), nil)
		} else {
			res, err = shell.Expand(v, nil)
		}
		if err != nil {
			return Dague{}, err
		}
		vars[k] = res
	}
	dague.Vars = vars

	// post step to expand mount path for the base image
	mounts := map[string]string{}
	for k, v := range dague.Go.Image.Mounts {
		expandedK, err := shell.Expand(k, dague.Vars)
		if err != nil {
			return Dague{}, err
		}
		mounts[expandedK] = v
	}
	dague.Go.Image.Mounts = mounts

	return dague, nil
}

func (d *Dague) VarsDup() map[string]string {
	vars := map[string]string{}
	for k, v := range d.Vars {
		vars[k] = v
	}
	return vars
}
