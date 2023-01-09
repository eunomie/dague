package daggers

import (
	"context"
	"fmt"
	"strings"

	"github.com/eunomie/dague"
)

func PrintGoformatter(ctx context.Context, c *Client, formatter string) error {
	out, err := SourcesNoDeps(c).WithExec(formatPrint(formatter)).Stdout(ctx)
	if err != nil {
		return err
	}
	if out != "" {
		return fmt.Errorf("")
	}
	return nil
}

func PrintFormatAndImports(ctx context.Context, c *Client, formatter string, locals []string) error {
	outFmt, err := SourcesNoDeps(c).WithExec(formatPrint(formatter)).Stdout(ctx)
	if err != nil {
		return err
	}
	outImports, err := SourcesNoDeps(c).WithExec(goImportsPrint(locals)).Stdout(ctx)
	if err != nil {
		return err
	}
	if outFmt != "" || outImports != "" {
		return fmt.Errorf("")
	}
	return nil
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
		Sources(c).WithExec(goImportsWrite(locals)).WithExec(formatWrite(formatter)),
		"*.go",
		"./",
	)
}

func goImportsWrite(locals []string) []string {
	args := []string{"goimports", "-w", "-format-only"}
	if len(locals) > 0 {
		args = append(args, "-local", strings.Join(locals, ","))
	}
	args = append(args, ".")
	return args
}

func goImportsPrint(locals []string) []string {
	args := []string{"goimports", "-d", "-e", "-format-only"}
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

func GoImportsWrite(ctx context.Context, c *Client, locals []string) error {
	return dague.ExportFilePattern(
		ctx,
		SourcesNoDeps(c).WithExec(goImportsWrite(locals)),
		"*.go",
		"./",
	)
}

func GoImportsPrint(ctx context.Context, c *Client, locals []string) error {
	return dague.ExportFilePattern(
		ctx,
		SourcesNoDeps(c).WithExec(goImportsPrint(locals)),
		"*.go",
		"./",
	)
}
