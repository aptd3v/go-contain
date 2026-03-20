package config

type SetConfig[T any] func(cfg *T) error

type SetInnerConfig[Inner any, Setters SetConfig[Inner]] func(setters ...Setters) error
