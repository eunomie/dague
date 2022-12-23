package daggers

import (
	"context"
	"strings"

	"dagger.io/dagger"

	"github.com/eunomie/dague"
)

func PrintGoformatter(ctx context.Context, c *Client, formatter string) error {
	return dague.Exec(
		ctx,
		SourcesNoDeps(c),
		formatPrint(formatter),
	)
}

func ApplyGoformatter(ctx context.Context, c *Client, formatter string) error {
	return dague.ExportFilePattern(
		ctx,
		Sources(c).Exec(formatWrite(formatter)),
		"*.go",
		"./",
	)
}

func ApplyFormatAndImports(ctx context.Context, c *Client, formatter string, locals []string) error {
	return dague.ExportFilePattern(
		ctx,
		Sources(c).Exec(goImports(locals)).Exec(formatWrite(formatter)),
		"*.go",
		"./",
	)
}

func goImports(locals []string) dagger.ContainerExecOpts {
	args := []string{"goimports", "-w", "-format-only"}
	if len(locals) > 0 {
		args = append(args, "-local", strings.Join(locals, ","))
	}
	args = append(args, ".")
	return dagger.ContainerExecOpts{Args: args}
}

func formatWrite(formatter string) dagger.ContainerExecOpts {
	return dagger.ContainerExecOpts{
		Args: []string{formatter, "-w", "."},
	}
}

func formatPrint(formatter string) dagger.ContainerExecOpts {
	return dagger.ContainerExecOpts{
		Args: []string{formatter, "-d", "-e", "."},
	}
}

func GoImports(ctx context.Context, c *Client, locals []string) error {
	return dague.ExportFilePattern(
		ctx,
		SourcesNoDeps(c).Exec(goImports(locals)),
		"*.go",
		"./",
	)
}
