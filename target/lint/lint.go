package lint

import (
	"context"

	"dagger.io/dagger"
	"github.com/eunomie/dague"
	"github.com/eunomie/dague/stages"

	"github.com/magefile/mage/mg"
)

type Lint mg.Namespace

// Govuln checks vulnerabilities in Go code
func (Lint) Govuln(ctx context.Context) error {
	return dague.RunInDagger(ctx, func(c *dagger.Client) error {
		return stages.GoVulnCheck(ctx, c)
	})
}
