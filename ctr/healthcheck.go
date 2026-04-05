package ctr

import (
	"fmt"
	"time"

	"github.com/aptd3v/containerkit/config"
	"github.com/aptd3v/containerkit/fields"
	"github.com/moby/moby/api/types/container"
)

type health struct{ ref *Container }

func newHealthCheckConfigurator(ctr *Container) *health {
	return &health{ref: ctr}
}

// Test sets the test of the health check
func (h *health) Test(args ...string) config.SetHealthConfig {
	return wrapVoid(func(cfg *container.HealthConfig) { cfg.Test = append(cfg.Test, args...) })
}

// Interval sets the interval of the health check
func (h *health) Interval(d time.Duration) config.SetHealthConfig {
	return wrapVoid(func(cfg *container.HealthConfig) { cfg.Interval = d })
}

// Timeout sets the timeout of the health check
func (h *health) Timeout(d time.Duration) config.SetHealthConfig {
	return wrapVoid(func(cfg *container.HealthConfig) { cfg.Timeout = d })
}

// Retries sets the retries of the health check
func (h *health) Retries(n int) config.SetHealthConfig {
	return wrapVoid(func(cfg *container.HealthConfig) { cfg.Retries = n })
}

// StartPeriod sets the start period of the health check
func (h *health) StartPeriod(d time.Duration) config.SetHealthConfig {
	return wrapVoid(func(cfg *container.HealthConfig) { cfg.StartPeriod = d })
}

// Fail is a function that returns an error
//
// note: this is useful for when you want to fail the health check config
// and append the error to the base config error collection
func (h *health) Fail(field fields.Field, err error) config.SetHealthConfig {
	return func(cfg *container.HealthConfig) error {
		return newBaseError(field, err)
	}
}

// Failf is a function that returns an error
//
// note: this is useful for when you want to fail the health check config
// and append the error to the base config error collection
func (h *health) Failf(field fields.Field, stringFormat string, args ...any) config.SetHealthConfig {
	return func(cfg *container.HealthConfig) error {
		return newBaseError(field, fmt.Errorf(stringFormat, args...))
	}
}
