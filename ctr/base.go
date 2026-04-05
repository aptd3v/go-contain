package ctr

import (
	"errors"
	"fmt"
	"time"

	"github.com/aptd3v/containerkit/config"
	"github.com/aptd3v/containerkit/errdefs"
	"github.com/aptd3v/containerkit/fields"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
)

type baseSetters struct {
	ref    *Container
	Health *health
}

func newBaseSetter(ctr *Container) *baseSetters {
	return &baseSetters{ref: ctr, Health: newHealthCheckConfigurator(ctr)}
}
func newBaseError(field fields.Field, err error) error {
	return errdefs.NewBaseConfigError(field, err)
}

// Hostname sets the hostname of the container
func (b *baseSetters) Hostname(name string) config.SetBaseConfig {
	return func(cfg *container.Config) error {
		if name == "" {
			return newBaseError(fields.Hostname, fmt.Errorf("hostname is empty '%s'", name))
		}
		if b.ref.HostConfig.UTSMode.IsHost() {
			return newBaseError(fields.Hostname, errors.New("UTS mode is set to host, cannot set hostname"))
		}
		cfg.Hostname = name
		return nil
	}
}

// Entrypoint appends the entrypoint to the container configuration
// Parameters:
//   - entrypoint: entrypoint and its arguments
func (b *baseSetters) Entrypoint(entrypoint ...string) config.SetBaseConfig {
	return func(cfg *container.Config) error {
		if len(entrypoint) == 0 {
			return newBaseError(fields.Entrypoint, fmt.Errorf("entrypoint is empty '%v'", entrypoint))
		}
		cfg.Entrypoint = append(cfg.Entrypoint, entrypoint...)
		return nil
	}
}

// Domainname sets the domain name of the container
func (b *baseSetters) Domainname(domainname string) config.SetBaseConfig {
	return func(cfg *container.Config) error {
		if domainname == "" {
			return newBaseError(fields.Domainname, fmt.Errorf("domainname is empty '%s'", domainname))
		}
		if b.ref.HostConfig.UTSMode.IsHost() {
			return newBaseError(fields.Domainname, errors.New("UTS mode is set to host, cannot set domainname"))
		}
		cfg.Domainname = domainname
		return nil
	}
}
func (b *baseSetters) HealthConfig(setters ...config.SetConfig[container.HealthConfig]) config.SetBaseConfig {
	return func(cfg *container.Config) error {
		cfg.Healthcheck = ensureNotNil(&cfg.Healthcheck)
		for _, set := range setters {
			if set == nil {
				continue
			}
			if err := set(cfg.Healthcheck); err != nil {
				return newBaseError(fields.HealthConfig, err)
			}
		}
		// if no test is set, set it to NONE
		if cfg.Healthcheck.Test == nil {
			cfg.Healthcheck.Test = []string{"NONE"}
		}
		return nil
	}
}

// Env appends an environment variable and its value to the container configuration
// Parameters:
//   - key: environment variable name
//   - value: environment variable value
func (b *baseSetters) Env(key string, value string) config.SetBaseConfig {
	return func(cfg *container.Config) error {
		if key == "" {
			return newBaseError(fields.Env, fmt.Errorf("key is empty '%s'", key))
		}
		if value == "" {
			return wrapVoid(func(cfg *container.Config) {
				cfg.Env = append(cfg.Env, fmt.Sprintf("%s=\"%s\"", key, value))
			})(cfg)
		}
		cfg.Env = append(cfg.Env, fmt.Sprintf("%s=%s", key, value))
		return nil
	}
}

// EnvMap appends a map of environment variables to the container configuration
// Parameters:
//   - env: map of environment variables
func (b *baseSetters) EnvMap(env map[string]string) config.SetBaseConfig {
	return func(cfg *container.Config) error {
		if cfg.Env == nil {
			cfg.Env = make([]string, 0)
		}
		for key, value := range env {
			if err := b.Env(key, value)(cfg); err != nil {
				return err
			}
		}
		return nil
	}
}

// ExposedPort appends a port to be exposed from the container
// Parameter:
//   - port: port number to be exposed from the container (e.g., "8080/tcp")
//
// note: default to tcp if no protocol is specified
func (b *baseSetters) ExposedPort(port string) config.SetBaseConfig {
	return func(cfg *container.Config) error {
		if cfg.ExposedPorts == nil {
			cfg.ExposedPorts = make(network.PortSet)
		}
		p, err := network.ParsePort(port)
		if err != nil {
			return newBaseError(fields.ExposedPorts, err)
		}
		cfg.ExposedPorts[p] = struct{}{}
		return nil
	}
}

// Image sets the image to use for the container
func (b *baseSetters) Image(image string) config.SetBaseConfig {
	return func(cfg *container.Config) error {
		if image == "" {
			return newBaseError(fields.Image, fmt.Errorf("image is empty '%s'", image))
		}
		cfg.Image = image
		return nil
	}
}

// Command appends the command to be run in the container
// Parameters:
//   - cmd: command
//   - args: command and its arguments
func (b *baseSetters) Command(cmd string, args ...string) config.SetBaseConfig {
	return func(cfg *container.Config) error {
		if len(cmd) == 0 {
			return newBaseError(fields.Cmd, fmt.Errorf("command is empty '%v'", cmd))
		}
		cfg.Cmd = append([]string{cmd}, args...)
		return nil
	}
}

// User sets the user that commands are run as inside the container
func (b *baseSetters) User(user string) config.SetBaseConfig {
	return func(cfg *container.Config) error {
		if user == "" {
			return newBaseError(fields.User, fmt.Errorf("user is empty '%s'", user))
		}
		cfg.User = user
		return nil
	}
}

// AttachedStdin enables attaching to container's standard input
func (b *baseSetters) AttachedStdin(attach bool) config.SetBaseConfig {
	return wrapVoid(func(cfg *container.Config) {
		cfg.AttachStdin = attach
	})
}

// AttachedStdout enables attaching to container's standard output
func (b *baseSetters) AttachedStdout(attach bool) config.SetBaseConfig {
	return wrapVoid(func(config *container.Config) {
		config.AttachStdout = attach
	})
}

// AttachedStderr enables attaching to container's standard error
func (b *baseSetters) AttachedStderr(attach bool) config.SetBaseConfig {
	return wrapVoid(func(cfg *container.Config) {
		cfg.AttachStderr = attach
	})
}

// Tty allocates a pseudo-TTY for the container
func (b *baseSetters) Tty(allocate bool) config.SetBaseConfig {
	return wrapVoid(func(cfg *container.Config) {
		cfg.Tty = allocate
	})
}

// StdinOpen keeps STDIN open even if not attached
func (b *baseSetters) StdinOpen(open bool) config.SetBaseConfig {
	return wrapVoid(func(cfg *container.Config) {
		cfg.OpenStdin = open
	})
}

// StdinOnce closes STDIN after the first attach
func (b *baseSetters) StdinOnce(once bool) config.SetBaseConfig {
	return wrapVoid(func(cfg *container.Config) {
		cfg.StdinOnce = once
	})
}

// EscapedArgs indicates that command arguments are already escaped
func (b *baseSetters) EscapedArgs(escaped bool) config.SetBaseConfig {
	return wrapVoid(func(cfg *container.Config) {
		cfg.ArgsEscaped = escaped
	})
}

// Volume appends a  short hand volume mount point to the container
// Parameter:
//   - volume: path where the volume should be mounted
func (b *baseSetters) Volume(volume string) config.SetBaseConfig {
	return func(cfg *container.Config) error {
		if volume == "" {
			return newBaseError(fields.Volumes, fmt.Errorf("volume is empty '%s'", volume))
		}
		if cfg.Volumes == nil {
			cfg.Volumes = make(map[string]struct{})
		}
		cfg.Volumes[volume] = struct{}{}
		return nil
	}
}

// WorkingDir sets the working directory for commands to run in
func (b *baseSetters) WorkingDir(dir string) config.SetBaseConfig {
	return func(cfg *container.Config) error {
		if dir == "" {
			return newBaseError(fields.WorkingDir, fmt.Errorf("working directory is empty '%s'", dir))
		}
		cfg.WorkingDir = dir
		return nil
	}
}

// DisabledNetwork disables networking for the container
func (b *baseSetters) DisabledNetwork(disabled bool) config.SetBaseConfig {
	return wrapVoid(func(cfg *container.Config) {
		cfg.NetworkDisabled = disabled
	})
}

// OnBuild appends ONBUILD metadata that will trigger when the image is used as a base image
func (b *baseSetters) OnBuild(args ...string) config.SetBaseConfig {
	return func(cfg *container.Config) error {
		if len(args) == 0 {
			return newBaseError(fields.OnBuild, fmt.Errorf("onbuild is empty '%v'", args))
		}
		cfg.OnBuild = append(cfg.OnBuild, args...)
		return nil
	}
}

// LabelMap appends a map of labels to the container
// Parameters:
//   - labels: map of labels
func (b *baseSetters) LabelMap(labels map[string]string) config.SetBaseConfig {
	return func(cfg *container.Config) error {
		if cfg.Labels == nil {
			cfg.Labels = make(map[string]string)
		}
		if len(labels) == 0 {
			return newBaseError(fields.Labels, fmt.Errorf("map is empty '%v'", labels))
		}
		for key, value := range labels {
			if err := b.Label(key, value)(cfg); err != nil {
				return err
			}
		}
		return nil
	}
}

// Label appends a label to the container
// Parameters:
//   - label: label key
//   - value: label value
func (b *baseSetters) Label(label, value string) config.SetBaseConfig {
	return func(cfg *container.Config) error {
		if cfg.Labels == nil {
			cfg.Labels = make(map[string]string)
		}
		if label == "" {
			return newBaseError(fields.Labels, fmt.Errorf("label is empty '%s'", label))
		}
		cfg.Labels[label] = value
		return nil
	}
}

// StopSignal sets the signal that will be used to stop the container
func (b *baseSetters) StopSignal(sig string) config.SetBaseConfig {
	return func(cfg *container.Config) error {
		if sig == "" {
			return newBaseError(fields.StopSignal, fmt.Errorf("signal is empty '%s'", sig))
		}
		cfg.StopSignal = sig
		return nil
	}
}

// Shell sets the shell for shell-form of RUN, CMD, ENTRYPOINT
func (b *baseSetters) Shell(shell ...string) config.SetBaseConfig {
	return func(cfg *container.Config) error {
		if len(shell) == 0 {
			return newBaseError(fields.Shell, fmt.Errorf("shell is empty '%v'", shell))
		}
		cfg.Shell = append(cfg.Shell, shell...)
		return nil
	}
}

// StopTimeout sets the timeout to stop the container
func (b *baseSetters) StopTimeout(timeout time.Duration) config.SetBaseConfig {
	return wrapVoid(func(cfg *container.Config) {
		v := int(timeout.Seconds())
		cfg.StopTimeout = &v
	})
}

// Fail is a function that returns an error
//
// note: this is useful for when you want to fail the container config
// and append the error to the container config error collection
func (b *baseSetters) Fail(field fields.Field, err error) config.SetBaseConfig {
	return func(cfg *container.Config) error {
		return newBaseError(field, err)
	}
}

// Failf is a function that returns an error
//
// note: this is useful for when you want to fail the container config
// and append the error to the container config error collection
func (b *baseSetters) Failf(field fields.Field, stringFormat string, args ...any) config.SetBaseConfig {
	return func(cfg *container.Config) error {
		return newBaseError(field, fmt.Errorf(stringFormat, args...))
	}
}
