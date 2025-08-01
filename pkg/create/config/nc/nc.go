// Package nc provides the options for the network config.
package nc

import (
	"fmt"

	"github.com/aptd3v/go-contain/pkg/create"
	"github.com/aptd3v/go-contain/pkg/create/config/nc/endpoint"
	"github.com/aptd3v/go-contain/pkg/create/errdefs"
	"github.com/docker/docker/api/types/network"
)

func WithEndpoint(name string, setters ...endpoint.SetEndpointConfig) create.SetNetworkConfig {
	return func(options *network.NetworkingConfig) error {
		if options.EndpointsConfig == nil {
			options.EndpointsConfig = make(map[string]*network.EndpointSettings)
		}
		if options.EndpointsConfig[name] == nil {
			options.EndpointsConfig[name] = &network.EndpointSettings{}
		}
		for _, set := range setters {
			if set != nil {
				if err := set(options.EndpointsConfig[name]); err != nil {
					return errdefs.NewNetworkConfigError("endpoint", fmt.Sprintf("failed to set endpoint: %s", err))
				}
			}
		}
		return nil
	}
}

// Fail is a function that returns an error
//
// note: this is useful for when you want to fail the network config
// and append the error to the network config error collection
func Fail(err error) create.SetNetworkConfig {
	return func(options *network.NetworkingConfig) error {
		return errdefs.NewNetworkConfigError("network_config", err.Error())
	}
}

// Failf is a function that returns an error
//
// note: this is useful for when you want to fail the network config
// and append the error to the network config error collection
func Failf(stringFormat string, args ...any) create.SetNetworkConfig {
	return func(options *network.NetworkingConfig) error {
		return errdefs.NewNetworkConfigError("network_config", fmt.Sprintf(stringFormat, args...))
	}
}
