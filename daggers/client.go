package daggers

import (
	"dagger.io/dagger"

	"github.com/eunomie/dague/config"
)

type Client struct {
	Dagger *dagger.Client
	Config *config.Dague
}

func NewClient(c *dagger.Client, conf *config.Dague) *Client {
	return &Client{
		Dagger: c,
		Config: conf,
	}
}
