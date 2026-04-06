package cli

import (
	"context"

	"github.com/aptd3v/containerkit/config"
	"github.com/moby/moby/client"
)

type ContainerAttachResult = client.ContainerAttachResult

func (c *Cli) ContainerAttach(ctx context.Context, containerID string, options ...config.SetContainerAttachOptions) (ContainerAttachResult, error) {

	opt := client.ContainerAttachOptions{}
	for _, option := range options {
		if err := option(&opt); err != nil {
			return ContainerAttachResult{}, err
		}
	}
	return c.client.ContainerAttach(ctx, containerID, opt)
}

type attachSetters struct {
	ref *Cli
}

// Stream sets the stream flag for the container attach options.
func (a *attachSetters) Stream(stream bool) config.SetContainerAttachOptions {
	return wrapVoid(func(opt *client.ContainerAttachOptions) {
		opt.Stream = stream
	})
}

// Stdin sets the stdin flag for the container attach options.
func (a *attachSetters) Stdin(stdin bool) config.SetContainerAttachOptions {
	return wrapVoid(func(o *client.ContainerAttachOptions) {
		o.Stdin = true
	})
}

// Stdout sets the stdout flag for the container attach options.
func (a *attachSetters) Stdout(stdout bool) config.SetContainerAttachOptions {
	return wrapVoid(func(o *client.ContainerAttachOptions) {
		o.Stdout = stdout
	})
}

// Stderr sets the stderr flag for the container attach options.
func (a *attachSetters) Stderr(stderr bool) config.SetContainerAttachOptions {
	return wrapVoid(func(o *client.ContainerAttachOptions) {
		o.Stderr = stderr
	})
}

// DetachKeys sets the detach keys for the container attach options.
func (a *attachSetters) DetachKeys(detachKeys string) config.SetContainerAttachOptions {
	return wrapVoid(func(o *client.ContainerAttachOptions) {
		o.DetachKeys = detachKeys
	})
}

// Logs sets the logs flag for the container attach options.
func (a *attachSetters) Logs(logs bool) config.SetContainerAttachOptions {
	return wrapVoid(func(o *client.ContainerAttachOptions) {
		o.Logs = logs
	})
}
