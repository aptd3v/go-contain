package config

// SetConfig function sets configuration value for any type.
type SetConfig[T any] func(cfg *T) error
