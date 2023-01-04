package daggers

import (
	"context"
	"strings"

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
		Sources(c).WithExec(formatWrite(formatter)),
		"*.go",
		"./",
	)
}

func ApplyFormatAndImports(ctx context.Context, c *Client, formatter string, locals []string) error {
	return dague.ExportFilePattern(
		ctx,
		Sources(c).WithExec(goImports(locals)).WithExec(formatWrite(formatter)),
		"*.go",
		"./",
	)
}

func goImports(locals []string) []string {
	args := []string{"goimports", "-w", "-format-only"}
	if len(locals) > 0 {
		args = append(args, "-local", strings.Join(locals, ","))
	}
	args = append(args, ".")
	return args
}

func formatWrite(formatter string) []string {
	return []string{formatter, "-w", "."}
}

func formatPrint(formatter string) []string {
	return []string{formatter, "-d", "-e", "."}
}

func GoImports(ctx context.Context, c *Client, locals []string) error {
	return dague.ExportFilePattern(
		ctx,
		SourcesNoDeps(c).WithExec(goImports(locals)),
		"*.go",
		"./",
	)
}
