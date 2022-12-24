package commands

import (
	"context"
	"fmt"

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
	l.register("go:fmt", goFmt)
	l.register("go:fmt:print", goFmtPrint)
	l.register("go:fmt:write", goFmtWrite)
	l.register("go:fmt:imports", goFmtImports)

	l.register("go:lint", goLint)
	l.register("go:lint:govuln", goLintGovuln)
	l.register("go:lint:golangci", goLintGolangCILint)

	l.register("go:deps", goDeps)
	l.register("go:mod", goMod)
	l.register("go:test", goTest)
	l.register("go:doc", goDoc)
	l.register("go:build", goBuild)

	l.register("task", task)
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
	return r
}
