package shell

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"mvdan.cc/sh/v3/shell"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func Interpret(ctx context.Context, cmd string, env map[string]string) (string, error) {
	out := bytes.NewBufferString("")

	if err := interpret(ctx, cmd, env, out, out); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

func Run(ctx context.Context, cmd string, env map[string]string) error {
	return interpret(ctx, cmd, env, os.Stdout, os.Stderr)
}

func interpret(ctx context.Context, cmd string, env map[string]string, outWriter, errWriter io.Writer) error {
	script, err := syntax.NewParser().Parse(strings.NewReader(cmd), "")
	if err != nil {
		return err
	}

	pairs := os.Environ()
	for k, v := range env {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
	}
	runner, err := interp.New(interp.Env(expand.ListEnviron(pairs...)), interp.StdIO(nil, outWriter, errWriter))
	if err != nil {
		return err
	}

	return runner.Run(ctx, script)
}

func Expand(s string, env map[string]string) (string, error) {
	return shell.Expand(s, func(name string) string {
		if env != nil {
			if v, ok := env[name]; ok {
				return v
			}
		}
		return os.Getenv(name)
	})
}
