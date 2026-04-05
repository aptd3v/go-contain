package ctr

import (
	"fmt"

	"github.com/aptd3v/containerkit/config"
	"github.com/aptd3v/containerkit/errdefs"
	"github.com/aptd3v/containerkit/fields"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type platformSetters struct {
	ref *Container
}

func newPlatformError(field fields.Field, err error) error {
	return errdefs.NewPlatformConfigError(field, err)
}
func newPlatformSetter(ctr *Container) *platformSetters {
	return &platformSetters{ref: ctr}
}

// Architecture sets the architecture for the platform
// Parameter:
//   - architecture: the architecture to be used for the platform
//
// Architecture field specifies the CPU architecture, for example `amd64` or `ppc64le`.
func (p *platformSetters) Architecture(architecture string) config.SetPlatformConfig {
	return wrapVoid(func(options *ocispec.Platform) {
		options.Architecture = architecture
	})
}

// OS sets the OS for the platform
// Parameter:
//   - OS: the OS to be used for the platform
//
// OS specifies the operating system, for example `linux` or `windows`.
func (p *platformSetters) OS(OS string) config.SetPlatformConfig {
	return wrapVoid(func(options *ocispec.Platform) {
		options.OS = OS
	})
}

// OSVersion sets the OS version for the platform
// Parameter:
//   - OSVersion: the OS version to be used for the platform
//
// OSVersion is an optional field specifying the operating system version, for example on Windows `10.0.14393.1066`.
func (p *platformSetters) OSVersion(OSVersion string) config.SetPlatformConfig {
	return wrapVoid(func(options *ocispec.Platform) {
		options.OSVersion = OSVersion
	})
}

// OSFeatures sets the OS features for the platform
// Parameter:
//   - OSFeatures: the OS features to be used for the platform
//
// OSFeatures is an optional field specifying an array of strings, each listing a required OS feature (for example on Windows `win32k`).
func (p *platformSetters) OSFeatures(OSFeatures ...string) config.SetPlatformConfig {
	return wrapVoid(func(options *ocispec.Platform) {
		if options.OSFeatures == nil {
			options.OSFeatures = make([]string, 0)
		}
		options.OSFeatures = append(options.OSFeatures, OSFeatures...)
	})
}

// Variant sets the variant for the platform
// Parameter:
//   - variant: the variant to be used for the platform
//
// Variant is an optional field specifying a variant of the CPU, for example `v7` to specify ARMv7 when architecture is `arm`.
func (p *platformSetters) Variant(variant string) config.SetPlatformConfig {
	return wrapVoid(func(options *ocispec.Platform) {
		options.Variant = variant
	})
}

// Fail is a function that returns an error
//
// note: this is useful for when you want to fail the platform config
// and append the error to the platform config error collection
func (p *platformSetters) Fail(field fields.Field, err error) config.SetPlatformConfig {
	return func(options *ocispec.Platform) error {
		return newPlatformError(field, err)
	}
}

// Failf is a function that returns an error
//
// note: this is useful for when you want to fail the platform config
// and append the error to the platform config error collection
func (p *platformSetters) Failf(field fields.Field, stringFormat string, args ...any) config.SetPlatformConfig {
	return func(options *ocispec.Platform) error {
		return newPlatformError(field, fmt.Errorf(stringFormat, args...))
	}
}
