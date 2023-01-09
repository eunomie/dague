package commands

import (
	"context"

	"github.com/eunomie/dague/config"
	"github.com/eunomie/dague/daggers"
)

func (l *List) goFmt(ctx context.Context, _ []string, conf *config.Dague, opts map[string]interface{}) error {
	return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
		check := false
		if v, ok := opts["check"]; ok {
			if b, ok := v.(bool); ok {
				check = b
			}
		}
		if check {
			return daggers.PrintFormatAndImports(ctx, c, conf.Go.Fmt.Formatter, conf.Go.Fmt.Goimports.Locals)
		}
		return daggers.ApplyFormatAndImports(ctx, c, conf.Go.Fmt.Formatter, conf.Go.Fmt.Goimports.Locals)
	})
}

func (l *List) goFmtPrint(ctx context.Context, _ []string, conf *config.Dague, _ map[string]interface{}) error {
	return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
		return daggers.PrintGoformatter(ctx, c, conf.Go.Fmt.Formatter)
	})
}

func (l *List) goFmtWrite(ctx context.Context, _ []string, conf *config.Dague, _ map[string]interface{}) error {
	return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
		return daggers.ApplyGoformatter(ctx, c, conf.Go.Fmt.Formatter)
	})
}

func (l *List) goImportsWrite(ctx context.Context, _ []string, conf *config.Dague, _ map[string]interface{}) error {
	return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
		return daggers.GoImportsWrite(ctx, c, conf.Go.Fmt.Goimports.Locals)
	})
}

func (l *List) goImportsPrint(ctx context.Context, _ []string, conf *config.Dague, _ map[string]interface{}) error {
	return daggers.RunInDagger(ctx, conf, func(c *daggers.Client) error {
		return daggers.GoImportsPrint(ctx, c, conf.Go.Fmt.Goimports.Locals)
	})
}
