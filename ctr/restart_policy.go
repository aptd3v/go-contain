package ctr

import (
	"github.com/aptd3v/containerkit/config"
	"github.com/moby/moby/api/types/container"
)

type RestartPolicy string

const (
	RestartPolicyNo            RestartPolicy = "no"
	RestartPolicyOnFailure     RestartPolicy = "on-failure"
	RestartPolicyAlways        RestartPolicy = "always"
	RestartPolicyUnlessStopped RestartPolicy = "unless-stopped"
	RestartPolicyDisabled      RestartPolicy = "disabled"
)

// RestartPolicy adds a restart policy to the host configuration.
// parameters:
//   - mode: the restart policy to use
//   - maxRetryCount: the maximum number of retries before giving up
func (h *hostSetters) RestartPolicy(mode RestartPolicy, maxRetryCount int) config.SetHostConfig {
	var pm container.RestartPolicyMode
	switch mode {
	case RestartPolicyNo:
		pm = container.RestartPolicyDisabled
	case RestartPolicyOnFailure:
		pm = container.RestartPolicyOnFailure
	case RestartPolicyAlways:
		pm = container.RestartPolicyAlways
	case RestartPolicyUnlessStopped:
		pm = container.RestartPolicyUnlessStopped
	case RestartPolicyDisabled:
		pm = container.RestartPolicyDisabled
	default:
		pm = container.RestartPolicyDisabled
	}
	return func(opt *container.HostConfig) error {
		opt.RestartPolicy = container.RestartPolicy{
			Name:              pm,
			MaximumRetryCount: maxRetryCount,
		}
		return nil
	}
}

// RestartPolicyAlways sets the restart policy to always
func (h *hostSetters) RestartPolicyAlways() config.SetHostConfig {
	return h.RestartPolicy(RestartPolicyAlways, 0)
}

// RestartPolicyOnFailure sets the restart policy to on-failure
// parameters:
//   - maxRetryCount: the maximum number of retries before giving up
func (h *hostSetters) RestartPolicyOnFailure(maxRetryCount int) config.SetHostConfig {
	return h.RestartPolicy(RestartPolicyOnFailure, maxRetryCount)
}

// RestartPolicyUnlessStopped sets the restart policy to unless-stopped
func (h *hostSetters) RestartPolicyUnlessStopped() config.SetHostConfig {
	return h.RestartPolicy(RestartPolicyUnlessStopped, 0)
}

// RestartPolicyNever sets the restart policy to no
func (h *hostSetters) RestartPolicyNever() config.SetHostConfig {
	return h.RestartPolicy(RestartPolicyNo, 0)
}
