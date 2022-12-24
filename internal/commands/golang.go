package commands

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"mvdan.cc/sh/v3/syntax"

	"mvdan.cc/sh/v3/expand"

	"mvdan.cc/sh/v3/interp"

	"mvdan.cc/sh/v3/shell"

	"github.com/eunomie/dague/config"
	"github.com/eunomie/dague/daggers"
	"github.com/eunomie/dague/types"
)

// GoDeps is a command to download go modules.
func (l *List) goDeps(ctx context.Context, _ []string, conf *config.Dague, _ map[string]interface{}) error {
	return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
		daggers.GoDeps(c)
		return nil
	})
}

// GoMod is a command to run go mod tidy and export go.mod and go.sum files.
func (l *List) goMod(ctx context.Context, _ []string, conf *config.Dague, _ map[string]interface{}) error {
	return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
		if err := daggers.ExportGoMod(ctx, c); err != nil {
			return err
		}
		return nil
	})
}

// GoTest is a command running Go tests.
func (l *List) goTest(ctx context.Context, _ []string, conf *config.Dague, _ map[string]interface{}) error {
	return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
		return daggers.RunGoTests(ctx, c)
	})
}

// GoDoc is a command generating Go documentation into readme.md files.
func (l *List) goDoc(ctx context.Context, _ []string, conf *config.Dague, opts map[string]interface{}) error {
	check := false
	if v, ok := opts["check"]; ok {
		if b, ok := v.(bool); ok {
			check = b
		}
	}
	return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
		if check {
			return daggers.CheckGoDoc(ctx, c)
		}
		return daggers.GoDoc(ctx, c)
	})
}

// GoBuild is a command to build a Go binary based on the local architecture.
func (l *List) goBuild(ctx context.Context, args []string, conf *config.Dague, _ map[string]interface{}) error {
	var targetName string

	if len(args) == 0 {
		var targetNames []string
		for _, t := range conf.Go.Build.Targets {
			targetNames = append(targetNames, t.Name)
		}
		qs := []*survey.Question{
			{
				Name: "target",
				Prompt: &survey.Select{
					Message: "Choose the target to build:",
					Options: targetNames,
				},
			},
		}
		answer := struct {
			Target string
		}{}
		if err := survey.Ask(qs, &answer); err != nil {
			return fmt.Errorf("could not select the target to build")
		}
		targetName = answer.Target
	} else {
		targetName = args[0]
	}

	var target config.Target
	var ok bool
	for _, t := range conf.Go.Build.Targets {
		if t.Name == targetName {
			target = t
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("could not find the target %q to build", targetName)
	}

	env := map[string]string{}

	for k, v := range target.Env {
		if strings.HasPrefix(v, "$ ") {
			shellCmd := strings.TrimPrefix(v, "$ ")
			value, err := interpretShell(ctx, shellCmd, env)
			if err != nil {
				return err
			}
			env[k] = value
		} else {
			env[k] = v
		}
	}

	var buildFlags []string
	if target.Ldflags != "" {
		flags, err := shell.Expand(target.Ldflags, func(s string) string {
			return env[s]
		})
		if err != nil {
			return err
		}
		buildFlags = append(buildFlags, "-ldflags="+flags)
	}
	return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
		if target.Type == "local" {
			return daggers.LocalBuild(ctx, c, types.LocalBuildOpts{
				BuildOpts: types.BuildOpts{
					Dir:        target.Out,
					In:         target.Path,
					EnvVars:    env,
					BuildFlags: buildFlags,
				},
				Out: filepath.Base(target.Path),
			})
		}
		var platforms []types.Platform
		for _, p := range target.Platforms {
			t := strings.SplitN(p, "/", 2)
			platforms = append(platforms, types.Platform{OS: t[0], Arch: t[1]})
		}
		return daggers.CrossBuild(ctx, c, types.CrossBuildOpts{
			BuildOpts: types.BuildOpts{
				Dir:        target.Out,
				In:         target.Path,
				EnvVars:    env,
				BuildFlags: buildFlags,
			},
			OutFileFormat: filepath.Base(target.Path) + "_%s_%s",
			Platforms:     platforms,
		})
	})
}

func interpretShell(ctx context.Context, cmd string, env map[string]string) (string, error) {
	script, err := syntax.NewParser().Parse(strings.NewReader(cmd), "")
	if err != nil {
		return "", err
	}

	out := bytes.NewBufferString("")

	pairs := os.Environ()
	for k, v := range env {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
	}
	runner, err := interp.New(interp.Env(expand.ListEnviron(pairs...)), interp.StdIO(nil, out, out))
	if err != nil {
		return "", err
	}

	if err = runner.Run(ctx, script); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}