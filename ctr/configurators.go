package ctr

import (
	"github.com/aptd3v/containerkit/config"
)

// ContainerConfigurator is a collection of setters for the overall container configuration.
type ContainerConfigurator interface {
	// WithBaseConfig sets the base configuration for the container
	// Parameters:
	//   - setters: setters for the base configuration
	WithBaseConfig(setters ...config.SetBaseConfig) ContainerConfigurator
	// WithHostConfig sets the host configuration for the container
	// Parameters:
	//   - setters: setters for the host configuration
	WithHostConfig(setters ...config.SetHostConfig) ContainerConfigurator
	// WithNetworkConfig sets the network configuration for the container
	// Parameters:
	//   - setters: setters for the network configuration
	WithNetworkConfig(setters ...config.SetNetworkConfig) ContainerConfigurator
	// WithPlatformConfig sets the platform configuration for the container
	// Parameters:
	//   - setters: setters for the platform configuration
	WithPlatformConfig(setters ...config.SetPlatformConfig) ContainerConfigurator
	// Errors returns the errors for the container configuration
	// Returns:
	//   - errors: errors for the container configuration
	Errors() []error
	// ErrorsEach iterates over the errors for the container configuration
	// Parameters:
	//   - cb: callback function
	ErrorsEach(cb func(err error) error) error
}
