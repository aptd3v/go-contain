package ctr

import (
	"net/netip"

	"github.com/aptd3v/containerkit/config"
	"github.com/aptd3v/containerkit/fields"
	"github.com/moby/moby/api/types/network"
)

type ipam struct {
	ref *Container
}

func newIPAMSetter(ctr *Container) *ipam {
	return &ipam{ref: ctr}
}

// IPv4Address sets the IPv4 address for the IPAM config.
// If the address is already set, it will be overwritten.
// Parameter:
//   - address: the IPv4 address to be used for the network
func (i *ipam) IPv4Address(address string) config.SetEndpointIPAMConfig {
	return func(opt *network.EndpointIPAMConfig) error {
		addr, err := netip.ParseAddr(address)
		if err != nil {
			return newNetworkError(fields.IPv4Address, err)
		}
		opt.IPv4Address = addr
		return nil
	}
}

// IPv6Address sets the IPv6 address for the IPAM config.
// If the address is already set, it will be overwritten.
// Parameter:
//   - address: the IPv6 address to be used for the network
func (i *ipam) IPv6Address(address string) config.SetEndpointIPAMConfig {
	return func(opt *network.EndpointIPAMConfig) error {
		addr, err := netip.ParseAddr(address)
		if err != nil {
			return newNetworkError(fields.IPv6Address, err)
		}
		opt.IPv6Address = addr
		return nil
	}
}

// WithLinkLocalIPs appends the link local IPs for the IPAM config
// Parameter:
//   - ips: the link local IPs to be used for the network
func (i *ipam) LinkLocalIPs(ips ...string) config.SetEndpointIPAMConfig {
	return func(opt *network.EndpointIPAMConfig) error {
		if opt.LinkLocalIPs == nil {
			opt.LinkLocalIPs = make([]netip.Addr, 0)
		}
		for _, ip := range ips {
			addr, err := netip.ParseAddr(ip)
			if err != nil {
				return newNetworkError(fields.LinkLocalIPs, err)
			}
			opt.LinkLocalIPs = append(opt.LinkLocalIPs, addr)
		}
		return nil
	}
}
