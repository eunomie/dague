package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/eunomie/dague/internal/ui"

	"github.com/eunomie/dague/config"
)

type (
	Runnable func(ctx context.Context, args []string, conf *config.Dague, opts map[string]interface{}) error
	List     struct {
		cmds map[string]Runnable
	}
)

func NewList() *List {
	l := &List{cmds: map[string]Runnable{}}
	l.register("go:fmt", l.goFmt)
	l.register("go:fmt:print", l.goFmtPrint)
	l.register("go:fmt:write", l.goFmtWrite)
	l.register("go:fmt:imports", l.goFmtImports)

	l.register("go:lint", l.goLint)
	l.register("go:lint:govuln", l.goLintGovuln)
	l.register("go:lint:golangci", l.goLintGolangCILint)

	l.register("go:mod", l.goMod)
	l.register("go:mod:download", l.goModDownload)
	l.register("go:test", l.goTest)
	l.register("go:doc", l.goDoc)
	l.register("go:build", l.goBuild)

	l.register("go:exec", l.goExec)

	l.register("task", l.task)
	return l
}

func (l *List) register(name string, runnable Runnable) {
	l.cmds[name] = runnable
}

func (l *List) Run(ctx context.Context, name string, args []string, conf *config.Dague, opts map[string]interface{}) error {
	r, ok := l.cmds[name]
	if !ok {
		return fmt.Errorf("not implemented")
	}

	_, _ = ui.Blue.Fprintf(os.Stderr, "%s %s\n", name, strings.Join(args, " "))

	return r(ctx, args, conf, opts)
}

func (l *List) RunDeps(ctx context.Context, deps []string, conf *config.Dague) error {
	for _, dep := range deps {
		cmd := strings.Split(dep, " ")
		name, args := cmd[0], cmd[1:]
		r, ok := l.cmds[name]
		if !ok {
			return fmt.Errorf("not implemented")
		}

		_, _ = ui.Purple.Fprintf(os.Stderr, "[-->] %s %s\n", name, strings.Join(args, " "))

		err := r(ctx, args, conf, nil)
		if err != nil {
			return err
		}

		_, _ = ui.Purple.Fprintf(os.Stderr, "[<--] %s %s\n", name, strings.Join(args, " "))
	}
	return nil
}
