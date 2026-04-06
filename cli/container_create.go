package cli

import (
	"context"

	"github.com/aptd3v/containerkit/ctr"
	"github.com/moby/moby/client"
)

type ContainerCreateResult = client.ContainerCreateResult

func (c *Cli) ContainerCreate(ctx context.Context, ctr *ctr.Container) (ContainerCreateResult, error) {
	opt := client.ContainerCreateOptions{
		Config:           ctr.BaseConfig,
		HostConfig:       ctr.HostConfig,
		NetworkingConfig: ctr.NetworkConfig,
		Platform:         ctr.PlatformConfig,
		Name:             ctr.Name,
	}
	return c.client.ContainerCreate(ctx, opt)
}
