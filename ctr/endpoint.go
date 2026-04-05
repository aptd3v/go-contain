package ctr

import (
	"net"
	"net/netip"

	"github.com/aptd3v/containerkit/config"
	"github.com/aptd3v/containerkit/fields"
	"github.com/moby/moby/api/types/network"
)

type endpointSetters struct {
	ref  *Container
	IPAM *ipam
}

func newEndpointSetter(ctr *Container) *endpointSetters {
	return &endpointSetters{ref: ctr, IPAM: newIPAMSetter(ctr)}
}

// IPAMConfig sets the IPAM configuration for the endpoint
// Parameters:
//   - setters: setters for the IPAM configuration
func (e *endpointSetters) IPAMConfig(setters ...config.SetEndpointIPAMConfig) config.SetEndpointConfig {
	return func(cfg *network.EndpointSettings) error {
		ipamConfig := network.EndpointIPAMConfig{}
		for _, setter := range setters {
			if setter == nil {
				continue
			}
			if err := setter(&ipamConfig); err != nil {
				return newNetworkError(fields.EndpointIPAMConfig, err)
			}
		}
		cfg.IPAMConfig = &ipamConfig
		return nil
	}
}

// Links appends the links for the endpoint
// Parameters:
//   - links: the links to be used for the endpoint
func (e *endpointSetters) Links(links ...string) config.SetEndpointConfig {
	return wrapVoid(func(cfg *network.EndpointSettings) {
		cfg.Links = append(cfg.Links, links...)
	})
}
func (e *endpointSetters) Aliases(aliases ...string) config.SetEndpointConfig {
	return wrapVoid(func(cfg *network.EndpointSettings) {
		cfg.Aliases = append(cfg.Aliases, aliases...)
	})
}

// MacAddress sets the MAC address for the endpoint
// Parameter:
//   - mac: MAC address to be used for the endpoint
func (e *endpointSetters) MacAddress(mac string) config.SetEndpointConfig {
	return func(cfg *network.EndpointSettings) error {
		hwAddr, err := net.ParseMAC(mac)
		if err != nil {
			return newNetworkError(fields.MacAddress, err)
		}
		cfg.MacAddress = network.HardwareAddr(hwAddr)
		return nil
	}
}

// DriverOptions sets the driver options for the endpoint
// Parameter:
//   - key: the key of the driver option
//   - value: the value of the driver option
func (e *endpointSetters) DriverOptions(key, value string) config.SetEndpointConfig {
	return func(options *network.EndpointSettings) error {
		if options.DriverOpts == nil {
			options.DriverOpts = make(map[string]string)
		}
		options.DriverOpts[key] = value
		return nil
	}
}

// Gateway sets the gateway for the endpoint
// Parameter:
//   - gateway: the gateway to be used for the endpoint
func (e *endpointSetters) Gateway(gateway string) config.SetEndpointConfig {
	return func(cfg *network.EndpointSettings) error {
		gatewayAddr, err := netip.ParseAddr(gateway)
		if err != nil {
			return newNetworkError(fields.Gateway, err)
		}
		cfg.Gateway = gatewayAddr
		return nil
	}
}

// GwPriority sets the gateway priority for the endpoint
// Parameter:
//   - priority: the priority to be used for the endpoint
func (e *endpointSetters) GwPriority(priority int) config.SetEndpointConfig {
	return wrapVoid(func(cfg *network.EndpointSettings) {
		cfg.GwPriority = priority
	})
}

// NetworkID sets the network ID for the endpoint
// Parameter:
//   - networkID: the network ID to be used for the endpoint
func (e *endpointSetters) NetworkID(networkID string) config.SetEndpointConfig {
	return wrapVoid(func(cfg *network.EndpointSettings) {
		cfg.NetworkID = networkID
	})
}

// EndpointID sets the endpoint ID for the endpoint
// Parameter:
//   - endpointID: the endpoint ID to be used for the endpoint
func (e *endpointSetters) EndpointID(endpointID string) config.SetEndpointConfig {
	return wrapVoid(func(cfg *network.EndpointSettings) {
		cfg.EndpointID = endpointID
	})
}

// IPAddress sets the IP address for the endpoint
// Parameter:
//   - ipAddress: the IP address to be used for the endpoint
func (e *endpointSetters) IPAddress(ipAddress string) config.SetEndpointConfig {
	return func(cfg *network.EndpointSettings) error {
		ipAddressAddr, err := netip.ParseAddr(ipAddress)
		if err != nil {
			return newNetworkError(fields.IPAddress, err)
		}
		cfg.IPAddress = ipAddressAddr
		return nil
	}
}

// IPPrefixLen sets the IP prefix length for the endpoint
// Parameter:
//   - ipPrefixLen: the IP prefix length to be used for the endpoint
func (e *endpointSetters) IPPrefixLen(ipPrefixLen int) config.SetEndpointConfig {
	return wrapVoid(func(cfg *network.EndpointSettings) {
		cfg.IPPrefixLen = ipPrefixLen
	})
}

// IPv6Gateway sets the IPv6 gateway for the endpoint
// Parameter:
//   - ipv6Gateway: the IPv6 gateway to be used for the endpoint
func (e *endpointSetters) IPv6Gateway(ipv6Gateway string) config.SetEndpointConfig {
	return func(cfg *network.EndpointSettings) error {
		ipv6GatewayAddr, err := netip.ParseAddr(ipv6Gateway)
		if err != nil {
			return newNetworkError(fields.IPv6Gateway, err)
		}
		cfg.IPv6Gateway = ipv6GatewayAddr
		return nil
	}
}

// GlobalIPv6Address sets the global IPv6 address for the endpoint
// Parameter:
//   - globalIPv6Address: the global IPv6 address to be used for the endpoint
func (e *endpointSetters) GlobalIPv6Address(globalIPv6Address string) config.SetEndpointConfig {
	return func(cfg *network.EndpointSettings) error {
		globalIPv6AddressAddr, err := netip.ParseAddr(globalIPv6Address)
		if err != nil {
			return newNetworkError(fields.GlobalIPv6Address, err)
		}
		cfg.GlobalIPv6Address = globalIPv6AddressAddr
		return nil
	}
}

// GlobalIPv6PrefixLen sets the global IPv6 prefix length for the endpoint
// Parameter:
//   - globalIPv6PrefixLen: the global IPv6 prefix length to be used for the endpoint
func (e *endpointSetters) GlobalIPv6PrefixLen(globalIPv6PrefixLen int) config.SetEndpointConfig {
	return wrapVoid(func(cfg *network.EndpointSettings) {
		cfg.GlobalIPv6PrefixLen = globalIPv6PrefixLen
	})
}

// DNSNames sets the DNS names for the endpoint
// Parameter:
//   - dnsNames: the DNS names to be used for the endpoint
func (e *endpointSetters) DNSNames(dnsNames ...string) config.SetEndpointConfig {
	return wrapVoid(func(cfg *network.EndpointSettings) {
		cfg.DNSNames = dnsNames
	})
}
