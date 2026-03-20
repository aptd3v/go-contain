package container

import (
	"github.com/aptd3v/containerkit/container/config"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type Container struct {
	config.ContainerConfiguator[container.Config, container.HostConfig, network.NetworkingConfig, ocispec.Platform]
	errs           []error
	BaseConfig     *container.Config
	Hostconfig     *container.HostConfig
	NetworkConfig  *network.NetworkingConfig
	PlatformConfig *ocispec.Platform
	Setters        *setters
}

//	func (c *Container) WithHostname(name string) config.SetConfigVoid[container.Config] {
//		return func(cfg *container.Config) { c.BaseConfig.Hostname = name }
//	}
//
//	func (c *Container) WithHealthConfig(setters ...config.SetConfig[container.HealthConfig]) error {
//		return applyInnerObjectError(c.BaseConfig.Healthcheck, setters...)
//	}

func (c *Container) WithBaseConfig(setters ...config.SetConfig[container.Config]) *Container {
	if c.BaseConfig == nil {
		c.BaseConfig = &container.Config{}
	}
	for _, set := range setters {
		if set == nil {
			continue
		}
		if err := set(c.BaseConfig); err != nil {
			c.errs = append(c.errs, err)
		}
	}

	return c
}
func (c *Container) WithHostConfig(setters ...config.SetConfig[container.HostConfig]) *Container {
	if c.Hostconfig == nil {
		c.Hostconfig = &container.HostConfig{}
	}
	for _, set := range setters {
		if set == nil {
			continue
		}
		if err := set(c.Hostconfig); err != nil {
			c.errs = append(c.errs, err)
		}
	}
	return c
}
func (c *Container) WithNetworkConfig(setters ...config.SetConfig[network.NetworkingConfig]) *Container {
	if c.NetworkConfig == nil {
		c.NetworkConfig = &network.NetworkingConfig{}
	}
	for _, set := range setters {
		if set == nil {
			continue
		}
		if err := set(c.NetworkConfig); err != nil {
			c.errs = append(c.errs, err)
		}
	}
	return c
}
func (c *Container) WithPlatformConfig(setters ...config.SetConfig[ocispec.Platform]) *Container {
	if c.PlatformConfig == nil {
		c.PlatformConfig = &ocispec.Platform{}
	}
	for _, set := range setters {
		if set == nil {
			continue
		}
		if err := set(c.PlatformConfig); err != nil {
			c.errs = append(c.errs, err)
		}
	}
	return c
}
func (c *Container) Errors() []error {
	return c.errs
}
func (c *Container) ErrorsEach(cb func(err error) error) error {
	if cb == nil {
		return nil
	}
	for _, e := range c.errs {
		if err := cb(e); err != nil {
			return err
		}
	}
	return nil
}

// Setters
type setters struct {
	Hostname func(name string) config.SetConfig[container.Config]

	//inner
	Health *health

	HealthConfig func(setters ...config.SetConfig[container.HealthConfig]) config.SetConfig[container.Config]
}

func (c *Container) With() *setters {
	if c.Setters != nil {
		return c.Setters
	}
	c.Setters = new(setters)
	c.Setters.Hostname = func(name string) config.SetConfig[container.Config] {
		return wrapVoid(func(cfg *container.Config) { cfg.Hostname = name })
	}

	c.setHealthcheck()
	c.Setters.HealthConfig = func(setters ...config.SetConfig[container.HealthConfig]) config.SetConfig[container.Config] {
		return func(cfg *container.Config) error {
			health := &container.HealthConfig{}
			applyInnerObjectError(c.errs, health, setters...)
			c.BaseConfig.Healthcheck = health
			return nil
		}
	}

	return c.Setters
}

func New() (ctr *Container, with *setters) {
	ctr = &Container{}
	return ctr, ctr.With()
}

func applyInnerObjectError[Object any, Setters config.SetConfig[Object]](errs []error, obj *Object, setters ...Setters) {
	for _, set := range setters {
		if set == nil {
			continue
		}
		if err := set(obj); err != nil {
			errs = append(errs, err)
		}

	}
}

func wrapVoid[T any](fn func(*T)) config.SetConfig[T] {
	return func(cfg *T) error {
		if fn == nil {
			return nil
		}
		fn(cfg)
		return nil
	}
}
