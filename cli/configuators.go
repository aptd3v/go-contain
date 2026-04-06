package cli

import "github.com/aptd3v/containerkit/config"

// CliConfigurator is a collection of setters for the overall cli configuration.
type CliConfigurator interface {
	WithAttachOptions(options ...config.SetContainerAttachOptions) CliConfigurator
}
