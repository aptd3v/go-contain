package errdefs

import (
	"errors"
	"fmt"

	"github.com/aptd3v/containerkit/fields"
)

var (
	// ErrContainerConfig is an error related to a container config
	ErrContainerConfig = errors.New("container config error")
	// ErrHostConfig is an error related to a host config
	ErrHostConfig = errors.New("host config error")
	// ErrNetworkConfig is an error related to a network config
	ErrNetworkConfig = errors.New("network config error")
	// ErrPlatformConfig is an error related to a platform config
	ErrPlatformConfig = errors.New("platform config error")
)

func formatError(field fields.Field, message string) string {
	return fmt.Sprintf("error setting %s: %s", field, message)
}

type containerConfigError struct {
	field   fields.Field
	message string
}

func (e *containerConfigError) Error() string {
	return formatError(e.field, e.message)
}

func (e *containerConfigError) Unwrap() error {
	return ErrContainerConfig
}

type hostConfigError struct {
	field   fields.Field
	message string
}

func (e *hostConfigError) Error() string {
	return formatError(e.field, e.message)
}

func (e *hostConfigError) Unwrap() error {
	return ErrHostConfig
}

type networkConfigError struct {
	field   fields.Field
	message string
}

func (e *networkConfigError) Error() string {
	return formatError(e.field, e.message)
}

func (e *networkConfigError) Unwrap() error {
	return ErrNetworkConfig
}

type platformConfigError struct {
	field   fields.Field
	message string
}

func (e *platformConfigError) Error() string {
	return formatError(e.field, e.message)
}

func (e *platformConfigError) Unwrap() error {
	return ErrPlatformConfig
}

// NewBaseConfigError creates a new error related to a container config
func NewBaseConfigError(field fields.Field, err error) error {
	if err == nil {
		return nil
	}
	return &containerConfigError{
		field:   field,
		message: err.Error(),
	}
}

// NewHostConfigError creates a new error related to a host config
func NewHostConfigError(field fields.Field, err error) error {
	if err == nil {
		return nil
	}
	return &hostConfigError{
		field:   field,
		message: err.Error(),
	}
}

// NewNetworkConfigError creates a new error related to a network config
func NewNetworkConfigError(field fields.Field, err error) error {
	if err == nil {
		return nil
	}
	return &networkConfigError{
		field:   field,
		message: err.Error(),
	}
}

// NewPlatformConfigError creates a new error related to a platform config
func NewPlatformConfigError(field fields.Field, err error) error {
	if err == nil {
		return nil
	}
	return &platformConfigError{
		field:   field,
		message: err.Error(),
	}
}
