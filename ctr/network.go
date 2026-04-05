package ctr

import (
	"fmt"

	"github.com/aptd3v/containerkit/config"
	"github.com/aptd3v/containerkit/errdefs"
	"github.com/aptd3v/containerkit/fields"
	"github.com/moby/moby/api/types/network"
)

type networkSetters struct {
	ref      *Container
	Endpoint *endpointSetters
}

func newNetworkSetter(ctr *Container) *networkSetters {
	return &networkSetters{ref: ctr, Endpoint: newEndpointSetter(ctr)}
}

func newNetworkError(field fields.Field, err error) error {
	return errdefs.NewNetworkConfigError(field, err)
}

// Endpoint sets the endpoint configuration for the network
// Parameters:
//   - name: the name of the endpoint
//   - setters: setters for the endpoint configuration
func (n *networkSetters) EndpointConfig(name string, setters ...config.SetEndpointConfig) config.SetNetworkConfig {
	return func(cfg *network.NetworkingConfig) error {
		if cfg.EndpointsConfig == nil {
			cfg.EndpointsConfig = make(map[string]*network.EndpointSettings)
		}
		if cfg.EndpointsConfig[name] == nil {
			cfg.EndpointsConfig[name] = &network.EndpointSettings{}
		}
		endpoint := network.EndpointSettings{}
		for _, setter := range setters {
			if setter == nil {
				continue
			}
			if err := setter(&endpoint); err != nil {
				return newNetworkError(fields.Endpoint, err)
			}
		}
		cfg.EndpointsConfig[name] = &endpoint
		return nil
	}
}

// Fail is a function that returns an error
//
// note: this is useful for when you want to fail the network config
// and append the error to the network config error collection
func (n *networkSetters) Fail(field fields.Field, err error) config.SetNetworkConfig {
	return func(cfg *network.NetworkingConfig) error {
		return newNetworkError(field, err)
	}
}

// Failf is a function that returns an error
//
// note: this is useful for when you want to fail the network config
// and append the error to the network config error collection
func (n *networkSetters) Failf(field fields.Field, stringFormat string, args ...any) config.SetNetworkConfig {
	return func(cfg *network.NetworkingConfig) error {
		return newNetworkError(field, fmt.Errorf(stringFormat, args...))
	}
}
