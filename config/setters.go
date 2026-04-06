package config

import "github.com/moby/moby/client"

// SetConfig function sets configuration value for any type.
type SetConfig[T any] func(cfg *T) error

// SetClientOption function sets configuration value for client.Opt.
type SetClientOption func() client.Opt
