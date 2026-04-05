package ctr

import (
	"errors"
	"fmt"
	"net/netip"
	"slices"
	"strconv"
	"strings"

	"github.com/aptd3v/containerkit/config"
	"github.com/aptd3v/containerkit/errdefs"
	"github.com/aptd3v/containerkit/fields"
	"github.com/moby/moby/api/types/blkiodev"
	"github.com/moby/moby/api/types/container"
	dtypes "github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/api/types/network"
)

type hostSetters struct {
	ref   *Container
	Mount *mountSetters
}

func newHostSetter(ctr *Container) *hostSetters {
	return &hostSetters{ref: ctr, Mount: newMountSetter(ctr)}
}
func newHostError(field fields.Field, err error) error {
	return errdefs.NewHostConfigError(field, err)
}

func (h *hostSetters) MountPoint(setters ...config.SetMountConfig) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		mounts := make([]dtypes.Mount, 0, len(setters))
		mount := dtypes.Mount{}
		for _, set := range setters {
			if set == nil {
				continue
			}
			if err := set(&mount); err != nil {
				return newHostError(fields.Mounts, err)
			}
		}
		mounts = append(mounts, mount)
		opt.Mounts = append(opt.Mounts, mounts...)
		return nil
	}
}

// AutoRemove sets the auto remove flag for the container in the host configuration.
// parameters:
//   - autoRemove: the auto remove flag
func (h *hostSetters) AutoRemove(autoRemove bool) config.SetHostConfig {
	return wrapVoid(func(cfg *container.HostConfig) {
		cfg.AutoRemove = autoRemove
	})
}

// MemoryLimit sets the memory limit for the container in the host configuration.
// parameters:
//   - memory: the memory limit in bytes
func (h *hostSetters) MemoryLimit(memory int64) config.SetHostConfig {
	return wrapVoid(func(cfg *container.HostConfig) {
		cfg.Memory = memory
	})
}

// PortBindings appends port mappings between the host and the container in the host configuration.
// parameters:
//   - hostPort: the port on the host
//   - containerPort: the port on the container
//   - hostIP: the IP address of the host
func (h *hostSetters) PortBindings(hostPort, containerPort, hostIP string) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if hostPort == "" {
			return newHostError(fields.PortBindings, errors.New("empty host port"))
		}
		if containerPort == "" {
			return newHostError(fields.PortBindings, errors.New("empty container port"))
		}
		if hostIP == "" {
			return newHostError(fields.PortBindings, errors.New("empty host IP"))
		}
		cPort, err := network.ParsePort(containerPort)
		if err != nil {
			return newHostError(fields.PortBindings, err)
		}
		hostPort, err := network.ParsePort(hostPort)
		if err != nil {
			return newHostError(fields.PortBindings, err)
		}

		if opt.PortBindings == nil {
			opt.PortBindings = make(network.PortMap)
		}
		hostIPAddr, err := netip.ParseAddr(hostIP)
		if err != nil {
			return newHostError(fields.PortBindings, err)
		}

		opt.PortBindings[cPort] = []network.PortBinding{
			{
				HostIP:   hostIPAddr,
				HostPort: hostPort.String(),
			},
		}
		return nil
	}
}

// DNSLookups appends a DNS server to the host configuration for the container.
// parameters:
//   - dnsLookups: the DNS server to add
//
// note: The strings can be in dotted decimal ("192.0.2.1"),
// IPv6 ("2001:db8::68"), or IPv6 with a scoped addressing zone ("fe80::1cc0:3e8c:119f:c2e1%ens18").
func (h *hostSetters) DNSLookups(dnsLookups ...string) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if len(dnsLookups) == 0 {
			return newHostError(fields.DNSLookups, errors.New("empty DNS lookups"))
		}
		for _, d := range dnsLookups {
			addr, err := netip.ParseAddr(d)
			if err != nil {
				return newHostError(fields.DNSLookups, err)
			}
			opt.DNS = append(opt.DNS, addr)
		}
		return nil
	}
}

// DNSOption appends a DNS option to the host configuration for the container.
// parameters:
//   - dnsOption: the DNS option to add
//
// note: A key-value pair representing a DNS option and its value.
// See your operating system's documentation for resolv.conf for valid options.
func (h *hostSetters) DNSOptions(dnsOption ...string) config.SetHostConfig {
	return wrapVoid(func(opt *container.HostConfig) {
		opt.DNSOptions = append(opt.DNSOptions, dnsOption...)
	})
}

// DNSSearch appends a DNS search domain to search non-fully qualified hostnames.
// parameters:
//   - search: the DNS search domain to add
func (h *hostSetters) DNSSearches(search ...string) config.SetHostConfig {
	return wrapVoid(func(opt *container.HostConfig) {
		opt.DNSSearch = append(opt.DNSSearch, search...)
	})
}

// ExtraHosts appends extra hosts to the host configuration for the container.
// parameters:
//   - extraHosts: the extra hosts to add
func (h *hostSetters) ExtraHosts(extraHosts ...string) config.SetHostConfig {
	return wrapVoid(func(opt *container.HostConfig) {
		opt.ExtraHosts = append(opt.ExtraHosts, extraHosts...)
	})
}

// AddedGroups appends supplementary groups to the host configuration for the container.
// parameters:
//   - group: the group to add
func (h *hostSetters) AddedGroups(group ...string) config.SetHostConfig {
	return wrapVoid(func(opt *container.HostConfig) {
		opt.GroupAdd = append(opt.GroupAdd, group...)
	})
}

// Binds appends a bind mount to the host configuration for the container.
// parameters:
//   - binds: the bind mount to add. Format is "source:target:mode". Paths may contain colons (e.g. Windows).
//     Mode may be single or comma-separated: ro, rw, readonly, z, Z; consistent, cached, delegated (macOS);
//     rshared, shared, rslave, slave, rprivate, private (propagation). E.g. "/etc/vector/vector.yml:/etc/vector/vector.yml:ro,z".
func (h *hostSetters) Binds(binds ...string) config.SetHostConfig {
	return func(opt *container.HostConfig) error {

		// this is in case user likes to use the setter multiple
		// times across multiple calls, check a single source of truth
		check := make([]string, 0, len(binds)+len(opt.Binds))
		check = append(check, binds...)
		check = append(check, opt.Binds...)
		if err := validateMounts(check); err != nil {
			return newHostError(fields.Binds, err)
		}
		opt.Binds = append(opt.Binds, binds...)
		return nil
	}
}

// validBindModeOptions is the set of bind mount mode tokens accepted by Docker/Compose.
// - ro, rw, readonly: read-only / read-write (readonly is alias for ro)
// - z, Z: SELinux labeling (shared / private)
// - consistent, cached, delegated: consistency (e.g. Docker Desktop macOS)
// - rshared, shared, rslave, slave, rprivate, private: bind propagation
var validBindModeOptions = map[string]struct{}{
	"ro": {}, "rw": {}, "readonly": {},
	"z": {}, "Z": {},
	"consistent": {}, "cached": {}, "delegated": {},
	"rshared": {}, "shared": {}, "rslave": {}, "slave": {}, "rprivate": {}, "private": {},
}

func validBindMode(mode string) bool {
	mode = strings.TrimSpace(mode)
	if mode == "" {
		return false
	}
	parts := strings.SplitSeq(mode, ",")
	for part := range parts {
		part = strings.TrimSpace(part)
		if _, ok := validBindModeOptions[part]; !ok {
			return false
		}
	}
	return true
}

// parseBindSpec parses one bind spec "source:target:mode", supporting paths that contain colons
// (e.g. Windows "C:\data:/app:ro") by using the segment after the last ":" as mode and
// splitting source and target on ":/" so container paths (starting  /) are unambiguous.
func parseBindSpec(spec string) (source, target, mode string, ok bool) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return "", "", "", false
	}
	last := strings.LastIndex(spec, ":")
	if last == -1 {
		return "", "", "", false
	}
	mode = strings.TrimSpace(spec[last+1:])
	rest := spec[:last]
	// Container path starts  /; split on ":/" so Windows host paths (e.g. C:\data) stay intact.
	sep := ":/"
	idx := strings.Index(rest, sep)
	if idx == -1 {
		return "", "", "", false
	}
	source = strings.TrimSpace(rest[:idx])
	target = strings.TrimSpace(rest[idx+len(sep):])
	return source, target, mode, true
}

// validateMounts validates mount specifications in the format "/source/path:/target/path:mode"
// or "source:target:mode" (paths may contain colons, e.g. Windows "C:\\data:/app:ro").
// Mode may be a single option or comma-separated (e.g. "ro", "rw,z", "ro,cached", "delegated").
func validateMounts(mounts []string) error {
	var errMsgs []string
	seenTargets := make(map[string]bool)

	for i, mount := range mounts {
		if mount = strings.TrimSpace(mount); mount == "" {
			errMsgs = append(errMsgs, fmt.Sprintf("binds[%d]: empty mount specification", i))
			continue
		}

		sourcePath, targetPath, mode, ok := parseBindSpec(mount)
		if !ok {
			errMsgs = append(errMsgs, fmt.Sprintf("binds[%d]: invalid format '%s' (must be 'source:target:mode')", i, mount))
			continue
		}
		if sourcePath == "" {
			errMsgs = append(errMsgs, fmt.Sprintf("binds[%d]: empty source path", i))
			continue
		}
		if targetPath == "" {
			errMsgs = append(errMsgs, fmt.Sprintf("binds[%d]: empty target path", i))
			continue
		}
		if seenTargets[targetPath] {
			errMsgs = append(errMsgs, fmt.Sprintf("binds[%d]: duplicate target path '%s'", i, targetPath))
			continue
		}
		seenTargets[targetPath] = true
		if !validBindMode(mode) {
			errMsgs = append(errMsgs, fmt.Sprintf("binds[%d]: invalid mode '%s' (e.g. 'ro', 'rw', 'z', 'Z', 'ro,z', 'ro,cached')", i, mode))
			continue
		}
	}

	if len(errMsgs) > 0 {
		return errors.New(strings.Join(errMsgs, "\n"))
	}
	return nil
}

// UTSModeHost sets the UTS (Unix Timesharing System) namespace mode to host.
//
// note: Docker disallows combining the hostname and domainname flags with uts=host.
// This is to prevent containers running in the host's UTS namespace from attempting
// to change the hosts configuration.
func (h *hostSetters) UTSModeHost() config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if opt.UTSMode.IsHost() {
			return newHostError(fields.UTSMode, errors.New("UTS mode is already set to host"))
		}
		if h.ref.BaseConfig.Hostname != "" {
			return newHostError(fields.UTSMode, errors.New("hostname is set, cannot use UTSModeHost"))
		}
		if h.ref.BaseConfig.Domainname != "" {
			return newHostError(fields.UTSMode, errors.New("domainname is set, cannot use UTSModeHost"))
		}
		opt.UTSMode = container.UTSMode("host")
		return nil
	}
}

// UserNSModeHost sets the user namespace mode for the container in the host configuration to host.
//
// note: If you enable user namespaces on the daemon, all containers are started with user namespaces
// enabled by default. To disable user namespace remapping for a specific container
func (h *hostSetters) UserNSModeHost() config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if opt.UsernsMode.IsHost() {
			return newHostError(fields.UsernsMode, errors.New("user namespace mode is already set to host"))
		}
		opt.UsernsMode = container.UsernsMode("host")
		return nil
	}
}

// ShmSize sets the size of the shared memory file system (/dev/shm) for the container in the host configuration to the specified size.
func (h *hostSetters) ShmSize(size int64) config.SetHostConfig {
	return wrapVoid(func(opt *container.HostConfig) {
		opt.ShmSize = size
	})
}

// Runtime sets the OCI runtime for the container in the host configuration.
// parameters:
//   - runtime: the runtime to use, e.g., "runc", "nvidia", or custom runtimes installed on the host.
func (h *hostSetters) Runtime(runtime string) config.SetHostConfig {
	return wrapVoid(func(opt *container.HostConfig) {
		opt.Runtime = runtime
	})
}

// ConsoleSize sets the console size for the container in the host configuration.
// parameters:
//   - height: the height of the console
//   - width: the width of the console
func (h *hostSetters) ConsoleSize(height uint, width uint) config.SetHostConfig {
	return wrapVoid(func(opt *container.HostConfig) {
		opt.ConsoleSize = [2]uint{height, width}
	})
}

// Isolation sets the isolation mode for the container in the host configuration.
// Parameters:
//   - isolation: the isolation mode to use. Valid values typically include "default", "process", and "hyperv".
//
// Note: This function applies isolation settings only on Windows environments.
func (h *hostSetters) Isolation(isolation string) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		switch isolation {
		case "", "default", "process", "hyperv":
			opt.Isolation = container.Isolation(isolation)
			return nil
		default:
			return newHostError(fields.Isolation, errors.New("invalid isolation mode"))
		}
	}
}

// CPUCount sets the number of CPUs allocated to the container in the host configuration.
//
// note: The number of CPUs must be at least 1. only available on Windows environments.
func (h *hostSetters) CPUCount(count int64) config.SetHostConfig {
	return wrapVoid(func(opt *container.HostConfig) {
		opt.CPUCount = count
	})
}

// CPUPercent sets the CPU percentage limit for the container.
//
// This option specifies the maximum amount of CPU the container can use as a percentage
// (e.g., 50 means the container can use up to 50% of one CPU core).
//
// note: The CPU percentage must be between 0 and 100. only available on Windows environments.
func (h *hostSetters) CPUPercent(percent int64) config.SetHostConfig {
	return wrapVoid(func(opt *container.HostConfig) {
		opt.CPUPercent = percent
	})
}

// ReadonlyPaths appends a list of paths to be marked as read-only inside the container.
//
// These paths are mounted as read-only from the host and are inaccessible for writing
// by any process in the container.
//
// Parameters:
//   - paths: one or more absolute paths to mark as read-only
//
// Returns an error if any path is empty or not absolute.
func (h *hostSetters) ReadonlyPaths(paths ...string) config.SetHostConfig {
	return func(opt *container.HostConfig) error {

		for i, path := range paths {
			// check duplicates
			if slices.Contains(opt.ReadonlyPaths, path) {
				return newHostError(fields.ReadonlyPaths, fmt.Errorf("path[%d]: duplicate path: %q", i, path))
			}
			if strings.TrimSpace(path) == "" {
				return newHostError(fields.ReadonlyPaths, fmt.Errorf("path[%d]: cannot be empty", i))
			}
			if !strings.HasPrefix(path, "/") {
				return newHostError(fields.ReadonlyPaths, fmt.Errorf("path[%d]: must be absolute: %q", i, path))
			}
		}
		opt.ReadonlyPaths = append(opt.ReadonlyPaths, paths...)
		return nil
	}
}

// MaskedPaths appends a list of paths to be masked inside the container.
//
// Masked paths are mounted as read-only and inaccessible from in the container,
// overriding the default masked paths (e.g., "/proc/kcore", "/proc/latency_stats").
//
// Parameters:
//   - paths: one or more absolute paths to mask
//
// Returns an error if any path is empty or not absolute.
func (h *hostSetters) MaskedPaths(paths ...string) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		for i, path := range paths {
			if strings.TrimSpace(path) == "" {
				return newHostError(fields.MaskedPaths, fmt.Errorf("path[%d]: cannot be empty", i))
			}
			if !strings.HasPrefix(path, "/") {
				return newHostError(fields.MaskedPaths, fmt.Errorf("path must be absolute: %q", path))
			}
		}
		opt.MaskedPaths = append(opt.MaskedPaths, paths...)
		return nil
	}
}

// NetworkMode sets the network mode for the container in the host configuration.
//
// Accepts standard Docker network modes such as:
//   - "bridge"    (default)
//   - "host"      (shares host's network stack)
//   - "none"      (no network)
//   - "container:<id|name>" (joins another container's network namespace)
//
// Returns an error if the mode is empty or not one of the supported formats.
func (h *hostSetters) NetworkMode(mode string) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		mode = strings.TrimSpace(mode)
		if mode == "" ||
			mode == "bridge" ||
			mode == "host" ||
			mode == "none" ||
			strings.HasPrefix(mode, "container:") {
			opt.NetworkMode = container.NetworkMode(mode)
			return nil
		}
		return newHostError(fields.NetworkMode, fmt.Errorf("invalid network mode: %q", mode))
	}
}

// VolumeDriver sets the volume driver for the container in the host configuration.
//
// The volume driver specifies the plugin or mechanism used to manage volumes mounted into the container.
// This option is typically only used  the Docker API directly, not  Docker Compose.
//
// Parameters:
//   - driver: the name of the volume driver to use.
//
// Returns an error if the driver name is empty or contains only whitespace.
func (h *hostSetters) VolumeDriver(driver string) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		driver = strings.TrimSpace(driver)
		if driver == "" {
			return newHostError(fields.VolumeDriver, errors.New("volume driver cannot be empty"))
		}
		opt.VolumeDriver = driver
		return nil
	}
}

// VolumesFrom appends a list of volumes to inherit from another container, specified in the form <container name>[:<ro|rw>].
// parameters:
//   - from: the container to inherit from
func (h *hostSetters) VolumesFrom(from string) config.SetHostConfig {
	return wrapVoid(func(opt *container.HostConfig) {
		opt.VolumesFrom = append(opt.VolumesFrom, from)
	})
}

// IpcMode sets the IPC (Inter-Process Communication) mode for the container.
//
// IPC mode controls how processes inside the container share memory and other IPC resources.
// Common valid values include:
//   - ""              (default — isolated IPC namespace)
//   - "host"          (use the host's IPC namespace)
//   - "private"       (use a private IPC namespace)
//   - "shareable"     (allows other containers to join this container's IPC namespace)
//   - "container:<name | id>" (join another container’s IPC namespace)
//
// Parameters:
//   - mode: the IPC mode to use.
//
// Returns an error if the mode is invalid or improperly formatted.
func (h *hostSetters) IpcMode(mode string) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if mode == "" || mode == "host" || mode == "private" || mode == "shareable" || strings.HasPrefix(mode, "container:") {
			opt.IpcMode = container.IpcMode(mode)
			return nil
		}
		return newHostError(fields.IpcMode, fmt.Errorf("invalid IPC mode: %q", mode))
	}
}

// Cgroup sets the cgroup namespace mode for the container in the host configuration.
//
// The cgroup namespace determines how the container is isolated in terms of resource control.
// Valid values include:
//   - ""           (default behavior — Docker decides)
//   - "host"       (container shares the host's cgroup namespace)
//   - "private"    (container gets its own cgroup namespace)
//   - "none"       (disable cgroup namespace, if supported)
//
// Parameters:
//   - cgroup: the cgroup namespace mode to use
//
// Returns an error if the cgroup mode is invalid.
func (h *hostSetters) Cgroup(cgroup string) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if cgroup == "" || cgroup == "host" || cgroup == "private" || cgroup == "none" {
			opt.Cgroup = container.CgroupSpec(cgroup)
			return nil
		}
		return newHostError(fields.Cgroup, fmt.Errorf("invalid cgroup mode: %q", cgroup))
	}
}

// OomScoreAdj sets the Out-Of-Memory (OOM) score adjustment for the container.
//
// The OOM score adjustment ranges from -1000 to 1000 and tells the kernel how likely
// it is to kill the container’s process when the system is under memory pressure.
// - A value closer to -1000 means the process is protected from OOM killing.
// - A value closer to 1000 makes the process more likely to be killed.
//
// Parameters:
//   - score: an integer between -1000 and 1000 (inclusive)
//
// Returns an error if the score is outside the valid range.
func (h *hostSetters) OomScoreAdj(score int) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if score < -1000 || score > 1000 {
			return newHostError(fields.OomScoreAdj, fmt.Errorf("value must be between -1000 and 1000"))
		}
		opt.OomScoreAdj = score
		return nil
	}
}

// OomKillDisable sets the OOM kill disable flag for the container in the host configuration.
// parameters:
//   - oomKillDisable: the OOM kill disable flag to use
func (h *hostSetters) OomKillDisable(disable bool) config.SetHostConfig {
	return wrapVoid(func(opt *container.HostConfig) {
		opt.OomKillDisable = &disable
	})
}

// PIDMode sets the PID mode for the container in the host configuration.
//
// PID mode controls the process ID namespace isolation of the container.
// Common values include "host", "container:<name|id>", or an empty string for private namespace.
//
// parameters:
//   - mode: the PID mode to use
func (h *hostSetters) PIDMode(mode string) config.SetHostConfig {
	return wrapVoid(func(opt *container.HostConfig) {
		opt.PidMode = container.PidMode(mode)
	})
}

// PublishAllPorts sets the publish all ports flag to true for the container
func (h *hostSetters) PublishAllPorts(publish bool) config.SetHostConfig {
	return wrapVoid(func(opt *container.HostConfig) {
		opt.PublishAllPorts = publish
	})
}

// ReadOnlyRootfs sets the readonly rootfs flag for the container in the host configuration.
// parameters:
//   - readonlyRootfs: the readonly rootfs flag to use
func (h *hostSetters) ReadonlyRootfs(readonly bool) config.SetHostConfig {
	return wrapVoid(func(opt *container.HostConfig) {
		opt.ReadonlyRootfs = readonly
	})
}

// SecurityOpts appends security options to the container's host configuration.
//
// These options customize security labels or settings used by Mandatory Access Control (MAC) systems
// like SELinux, AppArmor, or Seccomp. They allow fine-tuning container security contexts,
// such as specifying custom SELinux labels or disabling certain security profiles.
//
// Docker itself does minimal validation on these strings, so this function simply appends them.
// Users should provide valid options according to their security system's documentation.
//
// parameters:
//   - opts: a variadic list of security options to append (e.g., SELinux labels or profile overrides)
func (h *hostSetters) SecurityOpts(opts ...string) config.SetHostConfig {
	return wrapVoid(func(opt *container.HostConfig) {
		opt.SecurityOpt = append(opt.SecurityOpt, opts...)
	})
}

// StorageOpt appends a storage driver option for the container's host configuration.
// Storage options are key-value pairs passed to the container storage driver to configure
// specific behaviors like size limits, encryption, or performance settings.
//
// Since storage options are driver-specific and Docker itself does minimal validation,
// this function trims whitespace and requires a non-empty key to avoid invalid or
// meaningless options that might cause runtime errors or unexpected behavior.
//
// parameters:
//   - key: the storage option key (must be non-empty after trimming whitespace)
//   - value: the value for the storage option
func (h *hostSetters) StorageOpt(key, value string) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		key = strings.TrimSpace(key)
		if key == "" {
			return newHostError(fields.StorageOpt, errors.New("storage option key cannot be empty"))
		}
		if opt.StorageOpt == nil {
			opt.StorageOpt = make(map[string]string)
		}
		if _, ok := opt.StorageOpt[key]; ok {
			return newHostError(fields.StorageOpt, fmt.Errorf("storage option key already exists: %q", key))
		}
		opt.StorageOpt[key] = value
		return nil
	}
}

// Tmpfs appends to a map of tmpfs (mounts) used for the container
// parameters:
//   - key: the key of the tmpfs option
//   - value: the value of the tmpfs option
func (h *hostSetters) Tmpfs(key, value string) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if opt.Tmpfs == nil {
			opt.Tmpfs = make(map[string]string)
		}
		if _, ok := opt.Tmpfs[key]; ok {
			return newHostError(fields.Tmpfs, fmt.Errorf("tmpfs option key already exists: %q", key))
		}
		opt.Tmpfs[key] = value
		return nil
	}
}

// Privileged enables privileged mode for the container.
// Privileged mode grants the container extended Linux capabilities and access to devices.
func (h *hostSetters) Privileged(privileged bool) config.SetHostConfig {
	return wrapVoid(func(opt *container.HostConfig) {
		opt.Privileged = privileged
	})
}

// AddedDevice adds a device mapping to the container's host configuration.
//
// This allows the container to access a device from the host (e.g., /dev/snd, /dev/ttyUSB0).
//
// Parameters:
//   - hostPath: the path to the device on the host (e.g., "/dev/snd").
//   - containerPath: the path the device will be available at inside the container (e.g., "/dev/snd").
//   - permissions: cgroup permissions for the device ("r", "w", "m", or combinations like "rw").
//
// Returns an error if any parameter is empty or permissions contain invalid characters.
func (h *hostSetters) AddedDevice(hostPath string, containerPath string, permissions string) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if opt.Devices == nil {
			opt.Devices = make([]container.DeviceMapping, 0)
		}
		if strings.TrimSpace(hostPath) == "" {
			return newHostError(fields.Devices, errors.New("host device path cannot be empty"))
		}
		if strings.TrimSpace(containerPath) == "" {
			return newHostError(fields.Devices, errors.New("container device path cannot be empty"))
		}
		if !isValidDevicePermission(permissions) {
			return newHostError(fields.Devices, errors.New("invalid device permissions (must be combination of 'r', 'w', 'm')"))
		}
		opt.Devices = append(opt.Devices, container.DeviceMapping{
			PathOnHost:        hostPath,
			PathInContainer:   containerPath,
			CgroupPermissions: permissions,
		})
		return nil
	}
}

// isValidDevicePermission checks if permissions contain only valid cgroup characters: r (read), w (write), m (mknod).
func isValidDevicePermission(p string) bool {
	seen := map[rune]bool{}
	for _, r := range p {
		switch r {
		case 'r', 'w', 'm':
			if seen[r] {
				return false // disallow duplicates
			}
			seen[r] = true
		default:
			return false
		}
	}
	return len(p) > 0
}

// ContainerIDFile sets the path to a file where the container ID will be written after creation.
//
// After `client.ContainerCreate`, Docker will write the container's ID to this file.
// This is useful for external tooling or scripts that need to reference the container after it starts.
//
// Parameters:
//   - path: absolute or relative path to the container ID file.
//
// Returns an error if the path is empty or only whitespace.
func (h *hostSetters) ContainerIDFile(path string) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if strings.TrimSpace(path) == "" {
			return newHostError(fields.ContainerIDFile, errors.New("path cannot be empty"))
		}
		opt.ContainerIDFile = path
		return nil
	}
}

// CPUShares sets the CPU shares (relative weight) for the container.
//
// CPU shares define the relative CPU time available to the container compared to others.
// For example, 1024 is the default (normal priority), 512 is half the CPU weight, 2048 is double.
//
// Parameters:
//   - shares: the number of CPU shares (must be a positive integer).
//
// Returns an error if the value is less than or equal to zero.
func (h *hostSetters) CPUShares(shares int64) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if shares <= 0 {
			return newHostError(fields.CPUShares, errors.New("CPU shares must be a positive integer"))
		}
		opt.CPUShares = shares
		return nil
	}
}

// CPUPeriod sets the CPU CFS (Completely Fair Scheduler) period in microseconds.
//
// Parameters:
//   - period: the CPU period in microseconds (must be between 1,000 and 1,000,000 microseconds)
//
// Returns an error if the input is out of range or invalid.
func (h *hostSetters) CPUPeriod(periodMS int64) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if periodMS < 1000 || periodMS > 1_000_000 {
			return newHostError(fields.CPUPeriod, errors.New("CPU period must be between 1,000 and 1,000,000 microseconds"))
		}
		opt.CPUPeriod = periodMS
		return nil
	}
}

// CPUQuota sets the CPU CFS (Completely Fair Scheduler) quota for the container.
//
// Parameters:
//   - quota: the CPU quota in microseconds (must be a positive integer)
func (h *hostSetters) CPUQuota(quotaMS int64) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if quotaMS < 0 {
			return newHostError(fields.CPUQuota, errors.New("CPU quota must be a positive integer"))
		}
		opt.CPUQuota = quotaMS
		return nil
	}
}

// CpusetCpus sets the CPUs in which execution is allowed
// parameters:
//   - cpus: the CPUs in which execution is allowed
//
// note: Limit the specific CPUs or cores a container can use.
// A comma-separated list or hyphen-separated range of CPUs a container can use,
// if you have more than one CPU. The first CPU is numbered 0. A valid value might
// be 0-3 (to use the first, second, third, and fourth CPU) or 1,3 Limit the specific
// CPUs or cores a container can use. A comma-separated list or hyphen-separated range of
// CPUs a container can use, if you have more than one CPU. The first CPU is numbered 0.
// A valid value might be 0-3 (to use the first, second, third, and fourth CPU) or 1,3
func (h *hostSetters) CpusetCpus(cpus string) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if cpus == "" {
			return newHostError(fields.CpusetCpus, errors.New("cpus is required"))
		}
		if err := validateCpuset(cpus); err != nil {
			return newHostError(fields.CpusetCpus, errors.New("invalid cpuset format"))
		}
		opt.CpusetCpus = cpus
		return nil
	}
}

// MemoryReservation sets the memory soft limit
// parameters:
//   - memory: the memory soft limit
//
// note: the memory limit can be specified as an int64 in bytes
func (h *hostSetters) MemoryReservation(memory int64) config.SetHostConfig {
	return wrapVoid(func(opt *container.HostConfig) {
		opt.MemoryReservation = memory
	})
}

// MemorySwap sets the total memory limit (memory + swap)
// parameters:
//   - memorySwap: the total memory limit (memory + swap)
//
// note: the memory limit can be specified as an int64 in bytes
func (h *hostSetters) MemorySwap(memorySwap int64) config.SetHostConfig {
	return wrapVoid(func(opt *container.HostConfig) {
		opt.MemorySwap = memorySwap
	})
}

// Ulimits adds a user resource limit (ulimit) to the container's host configuration.
//
// Ulimits define resource constraints for processes running inside the container,
// such as the maximum number of open files ("nofile") or processes ("nproc").
//
// Parameters:
//   - name: the name of the resource to limit (e.g., "nofile", "nproc").
//   - soft: the soft limit, which is the value enforced for running processes.
//   - hard: the hard limit, which is the maximum value to which the soft limit can be raised.
//
// Returns an error if:
//   - the name is empty,
//   - any limit is negative,
//   - or the soft limit is greater than the hard limit.
func (h *hostSetters) Ulimits(name string, soft, hard int64) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if name == "" {
			return newHostError(fields.Ulimits, errors.New("ulimit name cannot be empty"))
		}
		if soft < 0 || hard < 0 {
			return newHostError(fields.Ulimits, errors.New("ulimit values must be non-negative"))
		}
		if soft > hard {
			return newHostError(fields.Ulimits, fmt.Errorf("soft limit (%d) cannot be greater than hard limit (%d)", soft, hard))
		}
		if opt.Ulimits == nil {
			opt.Ulimits = make([]*container.Ulimit, 0)
		}
		opt.Ulimits = append(opt.Ulimits, &container.Ulimit{
			Name: name,
			Soft: soft,
			Hard: hard,
		})
		return nil
	}
}

// Init sets the init flag for the container
// parameters:
//   - init: the init flag to use
//
// Note: use to run a custom init inside the container, if null, use the daemon's configured settings
func (h *hostSetters) Init(init bool) config.SetHostConfig {
	return wrapVoid(func(opt *container.HostConfig) {
		opt.Init = &init
	})
}

// CPURealtimePeriod sets the CPU real-time period in microseconds.
// This controls the scheduling period for real-time tasks in the container.
//
// Accepts either:
//   - an int value representing microseconds directly (e.g. 100000 for 100ms), or
//   - a duration string (e.g. "1ms", "1s", "750ms").
//
// Valid values range from 1,000 to 1,000,000 microseconds (1ms to 1s), inclusive.
func (h *hostSetters) CPURealtimePeriod(periodMS int64) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if periodMS < 1000 || periodMS > 1_000_000 {
			return newHostError(fields.CPURealtimePeriod, errors.New("CPU real-time period must be between 1,000 and 1,000,000 microseconds"))
		}
		opt.CPURealtimePeriod = periodMS
		return nil
	}
}

// CPURealtimeRuntime sets the CPU real-time runtime in microseconds.
// A value of -1 disables the limit (default).
// parameters:
//   - runtime: the CPU real-time runtime in microseconds
func (h *hostSetters) CPURealtimeRuntime(runtime int64) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if runtime < -1 {
			return newHostError(fields.CPURealtimeRuntime, errors.New("runtime must be -1 or >= 0"))
		}
		opt.CPURealtimeRuntime = runtime
		return nil
	}
}

// CpusetMems sets the memory nodes in which execution is allowed.
// Only effective on NUMA systems.
// parameters:
//   - mems: the memory nodes in which execution is allowed
func (h *hostSetters) CpusetMems(mems string) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if mems == "" {
			return newHostError(fields.CpusetMems, errors.New("mems is required"))
		}
		if err := validateCpuset(mems); err != nil {
			return newHostError(fields.CpusetMems, errors.New("invalid cpuset format"))
		}
		opt.CpusetMems = mems
		return nil
	}
}
func validateCpuset(mems string) error {
	for r := range strings.SplitSeq(mems, ",") {
		r = strings.TrimSpace(r)
		if r == "" {
			return fmt.Errorf("empty cpuset range")
		}
		parts := strings.Split(r, "-")
		if len(parts) == 1 {
			if _, err := strconv.Atoi(parts[0]); err != nil {
				return fmt.Errorf("invalid cpuset value: %s", r)
			}
		} else if len(parts) == 2 {
			start, err1 := strconv.Atoi(parts[0])
			end, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil || start > end {
				return fmt.Errorf("invalid cpuset range: %s", r)
			}
		} else {
			return fmt.Errorf("invalid cpuset format: %s", r)
		}
	}
	return nil
}

// MemorySwappiness tunes container memory swappiness (0 to 100).
// - A value of 0 turns off anonymous page swapping.
// - A value of 100 sets the host's swappiness value.
// - Values between 0 and 100 modify the swappiness level accordingly.
// parameters:
//   - swappiness: the swappiness level
func (h *hostSetters) MemorySwappiness(swappiness int64) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if swappiness < 0 || swappiness > 100 {
			return newHostError(fields.MemorySwappiness, errors.New("swappiness must be between 0 and 100"))
		}
		opt.MemorySwappiness = &swappiness
		return nil
	}
}

// PidsLimit sets the container's PIDs limit.
// parameters:
//   - limit: the PIDs limit
func (h *hostSetters) PidsLimit(limit int64) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if limit < -1 || limit == 0 {
			return newHostError(fields.PidsLimit, errors.New("limit must be -1 (unlimited) or a positive integer"))
		}
		opt.PidsLimit = &limit
		return nil
	}
}

// BlkioWeight sets the block IO weight (relative weight) for the container.
// Weight is a value between 10 and 1000.
// parameters:
//   - weight: the block IO weight
func (h *hostSetters) BlkioWeight(weight uint16) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if weight < 10 || weight > 1000 {
			return newHostError(fields.BlkioWeight, errors.New("weight must be between 10 and 1000"))
		}
		opt.BlkioWeight = weight
		return nil
	}
}

// BlkioDeviceReadBps appends a block IO read bandwidth throttle limit
// for a specific device to the container's host configuration.
// It limits the read rate (in bytes per second) on the specified device in the container.
//
// Parameters:
//   - devicePath: the device path (e.g., "/dev/sda")
//   - rate: the maximum read rate limit, either as an int (bytes per second)
//     or a human-readable string  units (e.g., "10MiB", "500KiB").
//
// Returns an error if the string rate cannot be parsed.
func (h *hostSetters) BlkioDeviceReadBps(devicePath string, rate uint64) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if devicePath == "" {
			return newHostError(fields.BlkioDeviceReadBps, errors.New("device path cannot be empty"))
		}
		opt.BlkioDeviceReadBps = append(opt.BlkioDeviceReadBps, &blkiodev.ThrottleDevice{
			Path: devicePath,
			Rate: rate,
		})
		return nil
	}
}

// BlkioDeviceWriteBps appends a block IO write bandwidth throttle limit
// for a specific device to the container's host configuration.
// It limits the write rate (in bytes per second) on the specified device in the container.
//
// Parameters:
//   - devicePath: the device path (e.g., "/dev/sda")
//   - rate: the maximum write rate limit, either as an int (bytes per second)
//     or a human-readable string  units (e.g., "10MiB", "500KiB").
//
// Returns an error if the string rate cannot be parsed.
func (h *hostSetters) BlkioDeviceWriteBps(devicePath string, rate uint64) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if devicePath == "" {
			return newHostError(fields.BlkioDeviceWriteBps, errors.New("device path cannot be empty"))
		}

		opt.BlkioDeviceWriteBps = append(opt.BlkioDeviceWriteBps, &blkiodev.ThrottleDevice{
			Path: devicePath,
			Rate: rate,
		})
		return nil
	}
}

// BlkioDeviceReadIOps appends a block IO read operations per second (IOPS) limit
// for a specific device to the container's host configuration.
//
// Parameters:
//   - devicePath: the device path (e.g., "/dev/sda")
//   - rate: the maximum read IOPS limit
//
// Returns an error if the string rate cannot be parsed.
func (h *hostSetters) BlkioDeviceReadIOps(devicePath string, rate uint64) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if devicePath == "" {
			return newHostError(fields.BlkioDeviceReadIOps, errors.New("device path cannot be empty"))
		}

		opt.BlkioDeviceReadIOps = append(opt.BlkioDeviceReadIOps, &blkiodev.ThrottleDevice{
			Path: devicePath,
			Rate: rate,
		})
		return nil
	}
}

// BlkioDeviceWriteIOps appends a block IO write operations per second (IOPS) limit
// for a specific device to the container's host configuration.
//
// Parameters:
//   - devicePath: the device path (e.g., "/dev/sda")
//   - rate: the maximum write IOPS limit
//
// Returns an error if the string rate cannot be parsed.
func (h *hostSetters) BlkioDeviceWriteIOps(devicePath string, rate uint64) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if opt.BlkioDeviceWriteIOps == nil {
			opt.BlkioDeviceWriteIOps = make([]*blkiodev.ThrottleDevice, 0)
		}
		if devicePath == "" {
			return newHostError(fields.BlkioDeviceWriteIOps, errors.New("device path cannot be empty"))
		}

		opt.BlkioDeviceWriteIOps = append(opt.BlkioDeviceWriteIOps, &blkiodev.ThrottleDevice{
			Path: devicePath,
			Rate: rate,
		})
		return nil
	}
}

// Sysctls adds appends a sysctl key-value pair in the container's host configuration.
// Sysctls allow tuning of kernel parameters inside the container.
//
// Parameters:
//   - key: the sysctl parameter name (e.g., "net.ipv4.ip_forward")
//   - value: the sysctl parameter value (e.g., "1")
//
// Note: The effectiveness of sysctls depends on the container runtime and host kernel support.
func (h *hostSetters) Sysctls(key, value string) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		if opt.Sysctls == nil {
			opt.Sysctls = make(map[string]string)
		}
		if key == "" {
			return newHostError(fields.Sysctls, errors.New("key cannot be empty"))
		}
		if _, ok := opt.Sysctls[key]; ok {
			return newHostError(fields.Sysctls, fmt.Errorf("sysctl key already exists: %q", key))
		}
		opt.Sysctls[key] = value
		return nil
	}
}

// DeviceCgroupRules appends device cgroup rules to the container's host configuration.
// Device cgroup rules control access to devices inside the container.
//
// Parameters:
//   - rules: one or more device cgroup rule strings (e.g., "c 1:3 rwm").
//
// Note: Rules must follow the device cgroup format accepted by the Linux kernel.
func (h *hostSetters) DeviceCgroupRules(rules ...string) config.SetHostConfig {
	return wrapVoid(func(opt *container.HostConfig) {
		opt.DeviceCgroupRules = append(opt.DeviceCgroupRules, rules...)
	})
}

// CgroupParent sets the cgroup parent for the container in the host configuration.
// This determines the parent cgroup under which the container's cgroup will be created.
//
// Parameters:
//   - parent: the cgroup parent path (e.g., "docker", "system.slice")
func (h *hostSetters) CgroupParent(parent string) config.SetHostConfig {
	return wrapVoid(func(opt *container.HostConfig) {
		opt.CgroupParent = parent
	})
}

// DeviceRequest adds a device request to the container's host configuration.
// Device requests specify special device access requirements, such as GPUs.
//
// Parameters:
//   - driver: the device driver name (e.g., "nvidia")
//   - count: the number of devices to request (use -1 for all available)
//   - deviceIDs: specific device IDs to request (empty for any)
//   - capabilities: list of capability sets required (e.g., [][]string{{"gpu"}, {"compute"}})
//
// Note: Device requests are commonly used for GPU access and require appropriate drivers.
func (h *hostSetters) DeviceRequest(driver string, count int, deviceIDs []string, capabilities [][]string) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		opt.DeviceRequests = append(opt.DeviceRequests, container.DeviceRequest{
			Driver:       driver,
			Count:        count,
			DeviceIDs:    deviceIDs,
			Capabilities: capabilities,
		})
		return nil
	}
}

// LogDriver sets the log driver and its options for the container.
//
// Parameters:
//   - driver: the log driver to use (e.g., "json-file", "syslog", "fluentd")
//   - config: a map of driver-specific options to configure the log driver
//
// Note: The supported log drivers and options depend on the container runtime and host configuration.
func (h *hostSetters) LogDriver(driver string, config map[string]string) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		opt.LogConfig = container.LogConfig{
			Type:   driver,
			Config: config,
		}
		return nil
	}
}

// Fail is a function that returns an error
//
// note: this is useful for when you want to fail the host config
// and append the error to the host config error collection
func (h *hostSetters) Fail(field fields.Field, err error) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		return newHostError(field, err)
	}
}

// Failf is a function that returns an error
//
// note: this is useful for when you want to fail the host config
// and append the error to the host config error collection
func (h *hostSetters) Failf(field fields.Field, stringFormat string, args ...any) config.SetHostConfig {
	return func(opt *container.HostConfig) error {
		return newHostError(field, fmt.Errorf(stringFormat, args...))
	}
}
