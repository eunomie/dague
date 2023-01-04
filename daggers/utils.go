package daggers

import (
	"context"

	"dagger.io/dagger"

	"github.com/eunomie/dague"
	"github.com/eunomie/dague/config"
)

func RunInDagger(ctx context.Context, conf *config.Dague, do func(*Client) error) error {
	return dague.RunInDagger(ctx, func(client *dagger.Client) error {
		c := NewClient(client, conf)
		return do(c)
	})
}

func applyBase(cont *dagger.Container, c *dagger.Client, conf *config.Dague) *dagger.Container {
	for host, guest := range conf.Go.Image.Mounts {
		cont = cont.WithMountedDirectory(guest, c.Host().Directory(host))
	}
	for k, v := range conf.Go.Image.Env {
		cont = cont.WithEnvVariable(k, v)
	}

	for _, cache := range conf.Go.Image.Caches {
		cacheVolume := c.CacheVolume(cache.Target)
		cont = cont.WithMountedCache(cache.Target, cacheVolume)
	}

	return cont
}
