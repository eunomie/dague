package commands

import (
	"context"
	"fmt"

	"github.com/eunomie/dague/config"
	"github.com/eunomie/dague/daggers"
)

func (l *List) goLint(ctx context.Context, _ []string, conf *config.Dague, _ map[string]interface{}) error {
	return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
		if conf.Go.Lint.Govulncheck.Enable {
			err := daggers.GoVulnCheck(ctx, c)
			if err != nil {
				return err
			}
		}
		if conf.Go.Lint.Golangci.Enable {
			return daggers.GolangCILint(ctx, c)
		}
		return nil
	})
}

func (l *List) goLintGovuln(ctx context.Context, _ []string, conf *config.Dague, _ map[string]interface{}) error {
	return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
		if !conf.Go.Lint.Govulncheck.Enable {
			return fmt.Errorf("govulncheck must be enabled")
		}
		return daggers.GoVulnCheck(ctx, c)
	})
}

func (l *List) goLintGolangCILint(ctx context.Context, _ []string, conf *config.Dague, _ map[string]interface{}) error {
	return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
		if !conf.Go.Lint.Golangci.Enable {
			return fmt.Errorf("golangci-lint must be enabled")
		}
		return daggers.GolangCILint(ctx, c)
	})
}
