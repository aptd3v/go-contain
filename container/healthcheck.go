package container

import (
	"github.com/aptd3v/containerkit/container/config"
	"github.com/moby/moby/api/types/container"
)

type health struct {
	Test func(args ...string) config.SetConfig[container.HealthConfig]
}

func (c *Container) setHealthcheck() {
	if c.Setters.Health == nil {
		c.Setters.Health = &health{}
	}
	c.Setters.Health.Test = func(args ...string) config.SetConfig[container.HealthConfig] {
		return wrapVoid(func(cfg *container.HealthConfig) { cfg.Test = append(cfg.Test, args...) })
	}
}
