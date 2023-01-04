package commands

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/eunomie/dague/internal/shell"

	"github.com/AlecAivazis/survey/v2"
	"github.com/eunomie/dague/config"
)

func (l *List) task(ctx context.Context, args []string, conf *config.Dague, _ map[string]interface{}) error {
	var taskName string
	if len(args) == 0 {
		var taskNames []string
		for k := range conf.Tasks {
			taskNames = append(taskNames, k)
		}
		sort.Strings(taskNames)
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

	for _, dep := range task.Deps {
		cmd := strings.Split(dep, " ")
		name, args := cmd[0], cmd[1:]
		if err := l.Run(name)(ctx, args, conf, nil); err != nil {
			return err
		}
	}

	if task.Cmds == "" {
		return nil
	}
	return shell.Run(ctx, task.Cmds, conf.Vars)
}
