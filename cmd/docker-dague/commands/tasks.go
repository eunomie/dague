package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"

	"github.com/eunomie/dague/config"
)

func TaskCommands(conf *config.Dague) []*cobra.Command {
	return []*cobra.Command{
		Task(conf),
	}
}

func Task(conf *config.Dague) *cobra.Command {
	return &cobra.Command{
		Use:   "task [TASK]",
		Short: "Run tasks",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			var taskName string
			if len(args) == 0 {
				var taskNames []string
				for k := range conf.Tasks {
					taskNames = append(taskNames, k)
				}
				answer := struct{ Task string }{}
				err := survey.Ask([]*survey.Question{
					{
						Name: "task",
						Prompt: &survey.Select{
							Message: "Choose the task to run:",
							Options: taskNames,
						},
					},
				}, &answer)
				if err != nil {
					return err
				}
				taskName = answer.Task
			} else {
				taskName = args[0]
			}

			task, ok := conf.Tasks[taskName]
			if !ok {
				return fmt.Errorf("could not find the task %q to run", taskName)
			}

			return runShell(ctx, task)
		},
	}
}

func runShell(ctx context.Context, cmd string) error {
	script, err := syntax.NewParser().Parse(strings.NewReader(cmd), "")
	if err != nil {
		return err
	}

	runner, err := interp.New(interp.Env(expand.ListEnviron(os.Environ()...)), interp.StdIO(nil, os.Stdout, os.Stderr))
	if err != nil {
		return err
	}

	return runner.Run(ctx, script)
}
