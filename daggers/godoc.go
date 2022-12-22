package daggers

import (
	"context"

	"dagger.io/dagger"

	"github.com/eunomie/dague"
)

func GoDoc(ctx context.Context, c *Client) error {
	return dague.ExportFilePattern(
		ctx,
		SourcesNoDeps(c).Exec(dagger.ContainerExecOpts{
			Args: []string{"gomarkdoc", "-u", "-e", "-o", "{{.Dir}}/README.md", "./..."},
		}),
		"*.md",
		".",
	)
}

func CheckGoDoc(ctx context.Context, c *Client) error {
	return dague.Exec(
		ctx,
		SourcesNoDeps(c),
		dagger.ContainerExecOpts{
			Args: []string{"gomarkdoc", "-c", "-u", "-e", "-o", "{{.Dir}}/README.md", "./..."},
		},
	)
}
