// Package network provides functions to set the network configuration for a project
package network

import (
	"fmt"

	"github.com/aptd3v/go-contain/pkg/create/config/sc/network/pool"
	"github.com/aptd3v/go-contain/pkg/create/errdefs"
	"github.com/compose-spec/compose-go/v2/types"
)

// SetNetworkProjectConfig is a function that sets the network configuration for a project
type SetNetworkProjectConfig func(*types.NetworkConfig) error

// WithDriver sets the driver for the network
// parameters:
//   - driver: the driver for the network
func WithDriver(driver string) SetNetworkProjectConfig {
	return func(opt *types.NetworkConfig) error {
		opt.Driver = driver
		return nil
	}
}

// WithDriverOptions sets the driver options for the network
// parameters:
//   - key: the key of the driver option
//   - value: the value of the driver option
func WithDriverOptions(key, value string) SetNetworkProjectConfig {
	return func(opt *types.NetworkConfig) error {
		if opt.DriverOpts == nil {
			opt.DriverOpts = make(map[string]string)
		}
		opt.DriverOpts[key] = value
		return nil
	}
}

// WithIpam sets the ipam driver for the network
// parameters:
//   - driver: the driver for the ipam
func WithIpamDriver(driver string) SetNetworkProjectConfig {
	return func(opt *types.NetworkConfig) error {
		opt.Ipam.Driver = driver
		return nil
	}
}

// WithIpamPool sets the ipam pool for the network
// parameters:
//   - setters: the setters for the ipam pool
func WithIpamPool(setters ...pool.SetIpamPoolProjectConfig) SetNetworkProjectConfig {
	return func(opt *types.NetworkConfig) error {
		if len(setters) == 0 {
			return nil
		}
		if opt.Ipam.Config == nil {
			opt.Ipam.Config = make([]*types.IPAMPool, 0)
		}
		ipamPool := types.IPAMPool{}
		for _, pool := range setters {
			if pool == nil {
				continue
			}
			if err := pool(&ipamPool); err != nil {
				return err
			}
		}
		opt.Ipam.Config = append(opt.Ipam.Config, &ipamPool)
		return nil
	}
}

// WithInternal sets the network to be internal
// parameters:
//   - internal: the internal flag for the network
func WithInternal() SetNetworkProjectConfig {
	return func(opt *types.NetworkConfig) error {
		opt.Internal = true
		return nil
	}
}

// WithAttachable sets the network to be attachable
// parameters:
//   - attachable: the attachable flag for the network
func WithAttachable() SetNetworkProjectConfig {
	return func(opt *types.NetworkConfig) error {
		opt.Attachable = true
		return nil
	}
}

// WithEnableIPv6 sets the network to enable ipv6
// parameters:
//   - enableIPv6: the enable ipv6 flag for the network
func WithEnableIPv6() SetNetworkProjectConfig {
	return func(opt *types.NetworkConfig) error {
		enableIPv6 := true
		opt.EnableIPv6 = &enableIPv6
		return nil
	}
}

// WithLabel appends a label to the network
// parameters:
//   - key: the key of the label
//   - value: the value of the label
func WithLabel(key, value string) SetNetworkProjectConfig {
	return func(opt *types.NetworkConfig) error {
		if opt.Labels == nil {
			opt.Labels = make(map[string]string)
		}
		opt.Labels[key] = value
		return nil
	}
}

// Fail is a function that returns a setter that always returns the given error
//
// note: this is useful for when you want to fail the network config
// and append the error to the service config error collection
func Fail(err error) SetNetworkProjectConfig {
	return func(opt *types.NetworkConfig) error {
		return errdefs.NewServiceConfigError("network", err.Error())
	}
}

// Failf is a function that returns a setter that always returns the given error
//
// note: this is useful for when you want to fail the network config
// and append the error to the service config error collection
func Failf(stringFormat string, args ...any) SetNetworkProjectConfig {
	return func(opt *types.NetworkConfig) error {
		return errdefs.NewServiceConfigError("network", fmt.Sprintf(stringFormat, args...))
	}
}
