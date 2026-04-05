package ctr

import (
	"github.com/aptd3v/containerkit/config"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// BaseConfig is a type alias for moby type [container.Config]
type BaseConfig = container.Config

// HostConfig is a type alias for moby type [container.HostConfig]
type HostConfig = container.HostConfig

// NetworkingConfig is a type alias for moby type [network.NetworkingConfig]
type NetworkingConfig = network.NetworkingConfig

// PlatformConfig is a type alias for opencontainers type [ocispec.Platform]
type PlatformConfig = ocispec.Platform

// HealthConfig is a type alias for moby type [container.HealthConfig]
type HealthConfig = container.HealthConfig

// Container is holds the container configuration as well as its setters.
// It holds an error collection for the container configuration. that is
// collected by the various setters.
//
//	ctr, with := ctr.New()
//	ctr.WithBaseConfig(
//		with.Base.Hostname("hostname"),
//		with.Base.HealthConfig(
//			with.Base.Health.Test("CMD-SHELL", "curl -f http://localhost:8080/health || exit 1"),
//		),
//	)
//	len(ctr.Errors()) // 0
type Container struct {
	// errs is the collection of errors for the container configuration.
	errs []error
	// BaseConfig is the basic container related configuration options.
	BaseConfig *BaseConfig
	// HostConfig is the host related configuration options.
	HostConfig *HostConfig
	// NetworkConfig is the network related configuration options.
	NetworkConfig *NetworkingConfig
	// PlatformConfig is the platform related configuration options.
	PlatformConfig *PlatformConfig
	// Set is a collection of setters for the container configuration.
	// It is the second return value of [New]. Typically you would do
	// 		ctr, with := ctr.New()
	//		ctr.WithBaseConfig(with.Hostname("hostname"))
	Set *Setters
}

func (c *Container) WithBaseConfig(setters ...config.SetBaseConfig) ContainerConfigurator {
	c.BaseConfig = ensureNotNil(&c.BaseConfig)
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
func (c *Container) WithHostConfig(setters ...config.SetHostConfig) ContainerConfigurator {
	c.HostConfig = ensureNotNil(&c.HostConfig)
	for _, set := range setters {
		if set == nil {
			continue
		}
		if err := set(c.HostConfig); err != nil {
			c.errs = append(c.errs, err)
		}
	}
	return c
}
func (c *Container) WithNetworkConfig(setters ...config.SetNetworkConfig) ContainerConfigurator {
	c.NetworkConfig = ensureNotNil(&c.NetworkConfig)
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
func (c *Container) WithPlatformConfig(setters ...config.SetPlatformConfig) ContainerConfigurator {
	c.PlatformConfig = ensureNotNil(&c.PlatformConfig)
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

// Setters contains various setters acting like various namespaces for the container configuration.
type Setters struct {
	// Base is a collection of setters for the container base configuration.
	Base *baseSetters
	// Host is a collection of setters for the host configuration.
	Host *hostSetters
	// Network is a collection of setters for the network configuration.
	Network *networkSetters
	// Platform is a collection of setters for the platform configuration.
	Platform *platformSetters
}

// SettersExt is to extend the setters with additional setters.
type SettersExt[T any] struct {
	*Setters
	Ext T
}

func ExtendSetters[T any](c *Container, override T) (with *SettersExt[T]) {
	if c == nil {
		c = &Container{
			BaseConfig:     &BaseConfig{},
			HostConfig:     &HostConfig{},
			NetworkConfig:  &NetworkingConfig{},
			PlatformConfig: &PlatformConfig{},
		}
	}
	if c.Set == nil {
		c.createSetters()
	}
	with = &SettersExt[T]{
		Setters: c.Set,
		Ext:     override,
	}
	return with
}

func New() (c *Container, with *Setters) {
	c = &Container{
		BaseConfig:     &BaseConfig{},
		HostConfig:     &HostConfig{},
		NetworkConfig:  &NetworkingConfig{},
		PlatformConfig: &PlatformConfig{},
	}
	return c, c.createSetters()
}

// CreateSetters creates a new 'with' Setters for the container configuration.
// It is safe to call this method multiple times.
func (c *Container) createSetters() *Setters {
	if c.Set == nil {
		c.Set = &Setters{
			Base:     newBaseSetter(c),
			Host:     newHostSetter(c),
			Network:  newNetworkSetter(c),
			Platform: newPlatformSetter(c),
		}
	}
	return c.Set
}

// wrapVoid helper function wraps setter function that will never error,
// like simple string assignment.
func wrapVoid[T any](fn func(*T)) config.SetConfig[T] {
	return func(cfg *T) error {
		if fn == nil {
			return nil
		}
		fn(cfg)
		return nil
	}
}

func ensureNotNil[T any](ptr **T) *T {
	if *ptr == nil {
		*ptr = new(T)
	}
	return *ptr
}
