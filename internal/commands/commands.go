package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/eunomie/dague/internal/ui"

	"github.com/eunomie/dague/config"
)

type (
	Runnable func(ctx context.Context, args []string, conf *config.Dague, opts map[string]interface{}) error
	List     struct {
		cmds map[string]Runnable
	}
)

func notImplemented(_ context.Context, _ []string, _ *config.Dague, _ map[string]interface{}) error {
	return fmt.Errorf("not implemented")
}

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

func (l *List) Run(name string) Runnable {
	r, ok := l.cmds[name]
	if !ok {
		return notImplemented
	}

	ui.Blue.Fprintf(os.Stderr, "[%s]\n", name)
	return r
}
