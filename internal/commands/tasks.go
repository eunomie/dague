package commands

import (
	"context"
	"fmt"

	"github.com/eunomie/dague/internal/ui"

	"github.com/eunomie/dague/internal/shell"

	"github.com/eunomie/dague/config"
)

func (l *List) task(ctx context.Context, args []string, conf *config.Dague, _ map[string]interface{}) error {
	var taskName string
	if len(args) == 0 {
		var taskNames []string
		for k := range conf.Tasks {
			taskNames = append(taskNames, k)
		}
		if len(taskNames) == 0 {
			taskName = taskNames[0]
		} else {
			selected, err := ui.Select("Choose the task to run:", taskNames)
			if err != nil {
				return fmt.Errorf("could not select the task to run: %w", err)
			}
			taskName = selected
		}
	} else {
		taskName = args[0]
	}

	task, ok := conf.Tasks[taskName]
	if !ok {
		return fmt.Errorf("could not find the task %q to run", taskName)
	}

	if err := l.RunDeps(ctx, task.Deps, conf); err != nil {
		return err
	}

	if task.Cmds == "" {
		return nil
	}
	return shell.Run(ctx, task.Cmds, conf.Vars)
}
