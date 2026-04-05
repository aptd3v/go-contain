package config

import (
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/api/types/network"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// SetBaseConfig is a type alias for [SetConfig] [container.Config]
type SetBaseConfig = SetConfig[container.Config]

// SetHostConfig is a type alias for [SetConfig] [container.HostConfig]
type SetHostConfig = SetConfig[container.HostConfig]

// SetNetworkConfig is a type alias for [SetConfig] [network.NetworkingConfig]
type SetNetworkConfig = SetConfig[network.NetworkingConfig]

// SetPlatformConfig is a type alias for [SetConfig] [ocispec.Platform]
type SetPlatformConfig = SetConfig[ocispec.Platform]

// SetHealthConfig is a type alias for [SetConfig] [container.HealthConfig]
type SetHealthConfig = SetConfig[container.HealthConfig]

// SetMountConfig is a type alias for [SetConfig] [mount.Mount]
type SetMountConfig = SetConfig[mount.Mount]

// SetEndpointConfig is a type alias for [SetConfig] [network.EndpointConfig]
type SetEndpointConfig = SetConfig[network.EndpointSettings]

// SetEndpointIPAMConfig is a type alias for [SetConfig] [network.EndpointIPAMConfig]
type SetEndpointIPAMConfig = SetConfig[network.EndpointIPAMConfig]
