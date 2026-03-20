package config

import (
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type ContainerConfiguator[
	Base container.Config,
	Host container.HostConfig,
	Network network.NetworkingConfig,
	Platform ocispec.Platform,
] interface {
	WithBaseConfig(setters ...SetConfig[container.Config]) ContainerConfiguator[Base, Host, Network, Platform]
	WithHostConfig(setters ...SetConfig[container.HostConfig]) ContainerConfiguator[Base, Host, Network, Platform]
	WithNetworkConfig(setters ...SetConfig[network.NetworkingConfig]) ContainerConfiguator[Base, Host, Network, Platform]
	WithPlatformConfig(setters ...SetConfig[ocispec.Platform]) ContainerConfiguator[Base, Host, Network, Platform]
	Errors() []error
	ErrorsEach(func(err error) error) error
}
