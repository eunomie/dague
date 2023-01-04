package commands

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/eunomie/dague/internal/shell"

	"dagger.io/dagger"

	"github.com/eunomie/dague"

	"github.com/AlecAivazis/survey/v2"

	"github.com/eunomie/dague/config"
	"github.com/eunomie/dague/daggers"
	"github.com/eunomie/dague/types"
)

// goModDownload is a command to download go modules.
func (l *List) goModDownload(ctx context.Context, _ []string, conf *config.Dague, _ map[string]interface{}) error {
	return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
		daggers.GoDeps(c)
		return nil
	})
}

// goMod is a command to run go mod tidy and export go.mod and go.sum files.
func (l *List) goMod(ctx context.Context, _ []string, conf *config.Dague, _ map[string]interface{}) error {
	return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
		if err := daggers.ExportGoMod(ctx, c); err != nil {
			return err
		}
		return nil
	})
}

// goTest is a command running Go tests.
func (l *List) goTest(ctx context.Context, _ []string, conf *config.Dague, _ map[string]interface{}) error {
	return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
		return daggers.RunGoTests(ctx, c)
	})
}

// goDoc is a command generating Go documentation into readme.md files.
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

func (l *List) goExec(ctx context.Context, args []string, conf *config.Dague, _ map[string]interface{}) error {
	var execName string
	if len(args) == 0 {
		var execNames []string
		for k := range conf.Go.Exec {
			execNames = append(execNames, k)
		}
		sort.Strings(execNames)
		answer := struct {
			Exec string
		}{}
		if err := survey.Ask([]*survey.Question{
			{
				Name: "exec",
				Prompt: &survey.Select{
					Message: "Choose the task to run inside build container:",
					Options: execNames,
				},
			},
		}, &answer); err != nil {
			return fmt.Errorf("could not select the target to run: %w", err)
		}
		execName = answer.Exec
	} else {
		execName = args[0]
	}

	exec, ok := conf.Go.Exec[execName]
	if !ok {
		return fmt.Errorf("could not find the target %q to run", execName)
	}

	for _, dep := range exec.Deps {
		cmd := strings.Split(dep, " ")
		name, args := cmd[0], cmd[1:]
		if err := l.Run(name)(ctx, args, conf, nil); err != nil {
			return err
		}
	}

	return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
		execOpts := dagger.ContainerExecOpts{
			Args: []string{"sh", "-c", exec.Cmds},
		}
		if exec.Export.Path != "" && exec.Export.Pattern != "" {
			return dague.ExportFilePattern(ctx, daggers.Sources(c).Exec(execOpts), exec.Export.Pattern, exec.Export.Path)
		}
		return dague.Exec(ctx, daggers.Sources(c), execOpts)
	})
}

// goBuild is a command to build a Go binary based on the local architecture.
func (l *List) goBuild(ctx context.Context, args []string, conf *config.Dague, _ map[string]interface{}) error {
	var targetName string

	if len(args) == 0 {
		var targetNames []string
		for k := range conf.Go.Build.Targets {
			targetNames = append(targetNames, k)
		}
		sort.Strings(targetNames)
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
			return fmt.Errorf("could not select the target to build: %w", err)
		}
		targetName = answer.Target
	} else {
		targetName = args[0]
	}

	target, ok := conf.Go.Build.Targets[targetName]
	if !ok {
		return fmt.Errorf("could not find the target %q to build", targetName)
	}

	env := conf.VarsDup()

	for k, v := range target.Env {
		if strings.HasPrefix(v, "shell ") {
			shellCmd := strings.TrimPrefix(v, "shell ")
			value, err := shell.Interpret(ctx, shellCmd, env)
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
		flags, err := shell.Expand(target.Ldflags, env)
		if err != nil {
			return err
		}
		buildFlags = append(buildFlags, "-ldflags="+flags)
	}
	return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
		out := target.Out
		if out == "" {
			out = "./dist"
		}
		if len(target.Platforms) == 0 {
			// if platforms is not defined then we admit it's a local build
			return daggers.LocalBuild(ctx, c, types.LocalBuildOpts{
				BuildOpts: types.BuildOpts{
					Dir:        out,
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
				Dir:        out,
				In:         target.Path,
				EnvVars:    env,
				BuildFlags: buildFlags,
			},
			OutFileFormat: filepath.Base(target.Path) + "_%s_%s",
			Platforms:     platforms,
		})
	})
}
