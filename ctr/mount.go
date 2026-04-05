package ctr

import (
	"os"

	"github.com/aptd3v/containerkit/config"
	dtypes "github.com/moby/moby/api/types/mount"
)

// Propagation represents the propagation of a mount.
type Propagation string

const (
	// PropagationRPrivate RPRIVATE
	PropagationRPrivate Propagation = "rprivate"
	// PropagationPrivate PRIVATE
	PropagationPrivate Propagation = "private"
	// PropagationRShared RSHARED
	PropagationRShared Propagation = "rshared"
	// PropagationShared SHARED
	PropagationShared Propagation = "shared"
	// PropagationRSlave RSLAVE
	PropagationRSlave Propagation = "rslave"
	// PropagationSlave SLAVE
	PropagationSlave Propagation = "slave"
)

/*
MountType is constant for the type of mount

	// TypeBind is the type for mounting a host directory
	"bind"
	// TypeVolume is the type for remote storage volumes
	"volume"
	// TypeTmpfs is the type for mounting tmpfs
	"tmpfs"
	// TypeNamedPipe is the type for mounting Windows named pipes
	"npipe"
*/
type MountType string

const (
	// TypeBind is the type for mounting a host directory
	MountTypeBind MountType = "bind"
	// TypeVolume is the type for remote storage volumes
	MountTypeVolume MountType = "volume"
	// TypeTmpfs is the type for mounting tmpfs
	MountTypeTmpfs MountType = "tmpfs"
	// TypeNamedPipe is the type for mounting Windows named pipes
	MountTypeNamedPipe MountType = "npipe"
)

/*
Consistency is constant for the consistency of the mount

	// ConsistencyFull guarantees bind mount-like consistency
	"consistent"
	// ConsistencyCached mounts can cache read data and FS structure
	"cached"
	// ConsistencyDelegated mounts can cache read and written data and structure
	"delegated"
	// ConsistencyDefault provides "consistent" behavior unless overridden
	"default"
*/
type Consistency string

const (
	// ConsistencyFull guarantees bind mount-like consistency
	ConsistencyFull Consistency = "consistent"
	// ConsistencyCached mounts can cache read data and FS structure
	ConsistencyCached Consistency = "cached"
	// ConsistencyDelegated mounts can cache read and written data and structure
	ConsistencyDelegated Consistency = "delegated"
	// ConsistencyDefault provides "consistent" behavior unless overridden
	ConsistencyDefault Consistency = "default"
)

type mountSetters struct{ ref *Container }

func (m *mountSetters) Type(mountType MountType) config.SetMountConfig {
	return func(opt *dtypes.Mount) error {
		opt.Type = dtypes.Type(mountType)
		return nil
	}
}

func newMountSetter(ctr *Container) *mountSetters {
	return &mountSetters{ref: ctr}
}

// Source sets the mount source
// parameters:
//   - source: the source of the mount
func (m *mountSetters) Source(source string) config.SetMountConfig {
	return func(opt *dtypes.Mount) error {
		opt.Source = source
		return nil
	}
}

// Target sets the mount target
// parameters:
//   - target: the target of the mount
func (m *mountSetters) Target(target string) config.SetMountConfig {
	return func(opt *dtypes.Mount) error {
		opt.Target = target
		return nil
	}
}

// ReadOnly sets the mount read only, (attempts recursive read-only if possible)
func (m *mountSetters) ReadOnly(readonly bool) config.SetMountConfig {
	return func(opt *dtypes.Mount) error {
		opt.ReadOnly = readonly
		return nil
	}
}

// Consistency sets the consistency of the mount
// parameters:
//   - consistency: the consistency of the mount (Consistency represents the consistency requirements of a mount)
func (m *mountSetters) Consistency(consistency Consistency) config.SetMountConfig {
	return func(opt *dtypes.Mount) error {
		opt.Consistency = dtypes.Consistency(consistency)
		return nil
	}
}

// TmpfsSizeBytes sets the size of the tmpfs mount in bytes
// parameters:
//   - size: the size of the tmpfs mount in bytes
//
// This will be converted to an operating system specific value
// depending on the host. For example, on linux, it will be converted to
// use a 'k', 'm' or 'g' syntax. BSD, though not widely supported with
// docker, uses a straight byte value.
//
// Percentages are not supported.
func (m *mountSetters) TmpfsSizeBytes(size int64) config.SetMountConfig {
	return wrapVoid(func(opt *dtypes.Mount) {
		opt.TmpfsOptions = ensureNotNil(&opt.TmpfsOptions)
		opt.TmpfsOptions.SizeBytes = size
	})
}

// TmpfsMode sets the mode of the tmpfs mount
// parameters:
//   - mode: the mode of the tmpfs mount upon creation
func (m *mountSetters) TmpfsMode(mode os.FileMode) config.SetMountConfig {
	return wrapVoid(func(opt *dtypes.Mount) {
		opt.TmpfsOptions = ensureNotNil(&opt.TmpfsOptions)
		opt.TmpfsOptions.Mode = mode
	})
}

// TmpfsFlag sets the flag of the tmpfs mount
// parameters:
//   - flag: the flag to set on the tmpfs mount
//
// Example:
//   - TmpfsFlag("exec")
func (m *mountSetters) TmpfsFlag(flag string) config.SetMountConfig {
	return wrapVoid(func(opt *dtypes.Mount) {
		opt.TmpfsOptions = ensureNotNil(&opt.TmpfsOptions)
		if opt.TmpfsOptions.Options == nil {
			opt.TmpfsOptions.Options = make([][]string, 0)
		}
		opt.TmpfsOptions.Options = append(opt.TmpfsOptions.Options, []string{flag})
	})
}

// TmpfsKeyValue sets the key value pair of the tmpfs mount
// parameters:
//   - key: the key of the key value pair
//   - value: the value of the key value pair
//
// Example:
//   - TmpfsKeyValue("uid", "1000")
//   - TmpfsKeyValue("gid", "1000")
func (m *mountSetters) TmpfsKeyValue(key string, value string) config.SetMountConfig {
	return wrapVoid(func(opt *dtypes.Mount) {
		opt.TmpfsOptions = ensureNotNil(&opt.TmpfsOptions)
		opt.TmpfsOptions.Options = append(opt.TmpfsOptions.Options, []string{key, value})
	})
}

// VolumeNoCopy sets the no copy flag of the volume mount to true
func (m *mountSetters) VolumeNoCopy(noCopy bool) config.SetMountConfig {
	return wrapVoid(func(opt *dtypes.Mount) {
		opt.VolumeOptions = ensureNotNil(&opt.VolumeOptions)
		opt.VolumeOptions.NoCopy = noCopy
	})
}

// VolumeLabel sets the label of the volume mount
// parameters:
//   - key: the key of the label
//   - value: the value of the label
func (m *mountSetters) VolumeLabel(key string, value string) config.SetMountConfig {
	return wrapVoid(func(opt *dtypes.Mount) {
		opt.VolumeOptions = ensureNotNil(&opt.VolumeOptions)
		if opt.VolumeOptions.Labels == nil {
			opt.VolumeOptions.Labels = make(map[string]string)
		}
		opt.VolumeOptions.Labels[key] = value
	})
}

// VolumeSubPath sets the sub path of the volume mount
// parameters:
//   - subPath: the sub path of the volume will mount the sub directory of the volume
func (m *mountSetters) VolumeSubPath(subPath string) config.SetMountConfig {
	return wrapVoid(func(opt *dtypes.Mount) {
		opt.VolumeOptions = ensureNotNil(&opt.VolumeOptions)
		opt.VolumeOptions.Subpath = subPath
	})
}

// VolumeDriver sets the driver of the volume mount
// parameters:
//   - driver: the driver of the volume mount
//   - options: the options of the volume mount passed directly to the linux kernel mount -o
//   - device: the device of the volume mount
func (m *mountSetters) VolumeDriver(driver string, options string, device string) config.SetMountConfig {
	return wrapVoid(func(opt *dtypes.Mount) {
		opt.VolumeOptions = ensureNotNil(&opt.VolumeOptions)
		opt.VolumeOptions.DriverConfig = ensureNotNil(&opt.VolumeOptions.DriverConfig)
		opt.VolumeOptions.DriverConfig.Name = driver
		opt.VolumeOptions.DriverConfig.Options = map[string]string{
			"o":      options, // options passed directly to the linux kernel mount -o
			"device": device,  // device of the volume mount
		}
	})
}

// BindPropagation sets the propagation of the bind mount
// parameters:
//   - propagation: the propagation of the bind mount
func (m *mountSetters) BindPropagation(propagation Propagation) config.SetMountConfig {
	return wrapVoid(func(opt *dtypes.Mount) {
		opt.BindOptions = ensureNotNil(&opt.BindOptions)
		opt.BindOptions.Propagation = dtypes.Propagation(propagation)
	})
}

// BindNonRecursive sets the non recursive flag of the bind mount to true
func (m *mountSetters) BindNonRecursive(nonRecursive bool) config.SetMountConfig {
	return wrapVoid(func(opt *dtypes.Mount) {
		opt.BindOptions = ensureNotNil(&opt.BindOptions)
		opt.BindOptions.NonRecursive = nonRecursive
	})
}

// BindReadOnlyNonRecursive sets the read only non recursive flag of the bind mount to true
func (m *mountSetters) BindReadOnlyNonRecursive(nonRecursive bool) config.SetMountConfig {
	return wrapVoid(func(opt *dtypes.Mount) {
		opt.BindOptions = ensureNotNil(&opt.BindOptions)
		opt.BindOptions.ReadOnlyNonRecursive = nonRecursive
	})
}

// BindReadOnlyForceRecursive sets the read only force recursive flag of the bind mount to true
// ReadOnlyForceRecursive raises an error if the mount cannot be made recursively read-only.
func (m *mountSetters) BindReadOnlyForceRecursive(force bool) config.SetMountConfig {
	return wrapVoid(func(opt *dtypes.Mount) {
		opt.BindOptions = ensureNotNil(&opt.BindOptions)
		opt.BindOptions.ReadOnlyForceRecursive = force
	})
}

// BindCreateMountpoint creates a mountpoint for the bind mount
func (m *mountSetters) BindCreateMountpoint(createMountpoint bool) config.SetMountConfig {
	return wrapVoid(func(opt *dtypes.Mount) {
		opt.BindOptions = ensureNotNil(&opt.BindOptions)
		opt.BindOptions.CreateMountpoint = createMountpoint
	})
}
