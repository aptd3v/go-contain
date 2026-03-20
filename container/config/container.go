package config

import (
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type ContainerConfig interface {
	*container.Config | *container.HostConfig | *network.NetworkingConfig | *ocispec.Platform
}
