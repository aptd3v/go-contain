package main

import (
	"fmt"

	"github.com/aptd3v/containerkit/container"
)

func main() {
	c, with := container.New()
	c.WithBaseConfig(
		with.Hostname("banana"),
		with.HealthConfig(
			with.Health.Test("CMD", "/bin/bash", "-c", "echo", "hello"),
			with.Health.Test("CMD", "/bin/bash", "-c", "echo", "hello"),
		),
	)

	fmt.Println(c.BaseConfig.Hostname, c.BaseConfig.Healthcheck.Test)
}
