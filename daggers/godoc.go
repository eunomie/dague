package daggers

import (
	"context"

	"github.com/eunomie/dague"
)

func GoDoc(ctx context.Context, c *Client) error {
	return dague.ExportFilePattern(
		ctx,
		SourcesNoDeps(c).WithExec([]string{"gomarkdoc", "-u", "-e", "-o", "{{.Dir}}/README.md", "./..."}),
		"*.md",
		".",
	)
}

func CheckGoDoc(ctx context.Context, c *Client) error {
	return dague.Exec(
		ctx,
		SourcesNoDeps(c),
		[]string{"gomarkdoc", "-c", "-u", "-e", "-o", "{{.Dir}}/README.md", "./..."},
	)
}
