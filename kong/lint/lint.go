package lint

import (
	"context"

	"dagger.io/dagger"
	"github.com/eunomie/dague"
	"github.com/eunomie/dague/daggers"
)

type (
	// Lint regroups a default set of linters to analyse source code.
	Lint struct {
		Govuln GovulnCmd `cmd:"" group:"lint" help:"checks vulnerabilities in Go code"`
	}

	// GovulnCmd runs govulncheck on Go source code.
	GovulnCmd struct{}
)

// Run govulncheck in dagger
func (GovulnCmd) Run() error {
	ctx := context.Background()
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return daggers.GoVulnCheck(ctx, c)
	})
}
