package fields

type Field uint8

const (
	Unknown Field = iota
	//Base Fields

	Hostname
	Domainname
	User
	AttachStdin
	AttachStdout
	AttachStderr
	ExposedPorts
	Tty
	OpenStdin
	StdinOnce
	HealthConfig
	Env
	Cmd
	ArgsEscaped
	Image
	Volumes
	WorkingDir
	Entrypoint
	NetworkDisabled
	OnBuild
	Labels
	StopSignal
	StopTimeout
	Shell

	//Host Fields

	Memory
	Binds
	ContainerIDFile
	LogConfig
	NetworkMode
	PortBindings
	RestartPolicy
	AutoRemove
	VolumeDriver
	VolumesFrom
	ConsoleSize
	CapAdd
	CapDrop
	CgroupnsMode
	DNSLookups
	DNSOptions
	DNSSearch
	ExtraHosts
	GroupAdd
	IpcMode
	Cgroup
	Links
	OomScoreAdj
	PidMode
	Privileged
	PublishAllPorts
	ReadonlyRootfs
	SecurityOpt
	StorageOpt
	Tmpfs
	UTSMode
	UsernsMode
	ShmSize
	Sysctls
	Runtime
	Isolation
	Resources
	CPUShares
	Mounts
	MaskedPaths
	ReadonlyPaths
	Init
	Devices
	CPUPeriod
	CPUQuota
	CpusetCpus
	CpusetMems
	Ulimits
	CPURealtimeRuntime
	CPURealtimePeriod
	MemorySwappiness
	PidsLimit
	BlkioWeight
	BlkioDeviceReadBps
	BlkioDeviceWriteBps
	BlkioDeviceReadIOps
	BlkioDeviceWriteIOps
	DeviceCgroupRules
	CgroupParent

	//Network Fields

	Endpoint
	EndpointIPAMConfig
	IPv4Address
	IPv6Address
	LinkLocalIPs
	MacAddress
	Gateway
	IPAddress
	IPv6Gateway
	GlobalIPv6Address
	GlobalIPv6PrefixLen

	//Platform Fields
	Platform
	Architecture
	OS
	OSVersion
	OSFeatures
	Variant
)

func (id Field) String() string {
	switch id {
	case Unknown:
		return "unknown"
	case Hostname:
		return "hostname"
	case Domainname:
		return "domainname"
	case User:
		return "user"
	case AttachStdin:
		return "attach_stdin"
	case AttachStdout:
		return "attach_stdout"
	case AttachStderr:
		return "attach_stderr"
	case ExposedPorts:
		return "exposed_ports"
	case Tty:
		return "tty"
	case OpenStdin:
		return "open_stdin"
	case StdinOnce:
		return "stdin_once"
	case Env:
		return "env"
	case Cmd:
		return "cmd"
	case HealthConfig:
		return "healthconfig"
	case ArgsEscaped:
		return "args_escaped"
	case Image:
		return "image"
	case Volumes:
		return "volumes"
	case WorkingDir:
		return "working_dir"
	case Entrypoint:
		return "entrypoint"
	case NetworkDisabled:
		return "network_disabled"
	case OnBuild:
		return "on_build"
	case Labels:
		return "labels"
	case StopSignal:
		return "stop_signal"
	case StopTimeout:
		return "stop_timeout"
	case Shell:
		return "shell"
	case Memory:
		return "memory"
	case Binds:
		return "binds"
	case UTSMode:
		return "uts_mode"
	case ContainerIDFile:
		return "container_id_file"
	case LogConfig:
		return "log_config"
	case NetworkMode:
		return "network_mode"
	case PortBindings:
		return "port_bindings"
	case RestartPolicy:
		return "restart_policy"
	case AutoRemove:
		return "auto_remove"
	case VolumeDriver:
		return "volume_driver"
	case VolumesFrom:
		return "volumes_from"
	case ConsoleSize:
		return "console_size"
	case Devices:
		return "devices"
	case CapAdd:
		return "cap_add"
	case CapDrop:
		return "cap_drop"
	case CgroupnsMode:
		return "cgroupns_mode"
	case DNSLookups:
		return "dns_lookups"
	case DNSOptions:
		return "dns_options"
	case DNSSearch:
		return "dns_search"
	case ExtraHosts:
		return "extra_hosts"
	case GroupAdd:
		return "group_add"
	case IpcMode:
		return "ipc_mode"
	case Cgroup:
		return "cgroup"
	case Links:
		return "links"
	case OomScoreAdj:
		return "oom_score_adj"
	case CPUShares:
		return "cpu_shares"
	case CPUPeriod:
		return "cpu_period"
	case PidMode:
		return "pid_mode"
	case Privileged:
		return "privileged"
	case PublishAllPorts:
		return "publish_all_ports"
	case ReadonlyRootfs:
		return "readonly_rootfs"
	case SecurityOpt:
		return "security_opt"
	case StorageOpt:
		return "storage_opt"
	case Tmpfs:
		return "tmpfs"
	case Sysctls:
		return "sysctls"
	case Runtime:
		return "runtime"
	case Isolation:
		return "isolation"
	case Resources:
		return "resources"
	case CPUQuota:
		return "cpu_quota"
	case Mounts:
		return "mounts"
	case MaskedPaths:
		return "masked_paths"
	case ReadonlyPaths:
		return "readonly_paths"
	case CpusetCpus:
		return "cpuset_cpus"
	case CpusetMems:
		return "cpuset_mems"
	case CPURealtimeRuntime:
		return "cpu_realtime_runtime"
	case CPURealtimePeriod:
		return "cpu_realtime_period"
	case Ulimits:
		return "ulimits"
	case MemorySwappiness:
		return "memory_swappiness"
	case PidsLimit:
		return "pids_limit"
	case BlkioWeight:
		return "blkio_weight"
	case BlkioDeviceReadBps:
		return "blkio_device_read_bps"
	case BlkioDeviceWriteBps:
		return "blkio_device_write_bps"
	case BlkioDeviceReadIOps:
		return "blkio_device_read_iops"
	case BlkioDeviceWriteIOps:
		return "blkio_device_write_iops"
	case DeviceCgroupRules:
		return "device_cgroup_rules"
	case CgroupParent:
		return "cgroup_parent"
	case Endpoint:
		return "endpoint"
	case EndpointIPAMConfig:
		return "endpoint_ipam_config"
	case IPv4Address:
		return "ipv4_address"
	case IPv6Address:
		return "ipv6_address"
	case LinkLocalIPs:
		return "link_local_ips"
	case MacAddress:
		return "mac_address"
	case Gateway:
		return "gateway"
	case IPAddress:
		return "ip_address"
	case IPv6Gateway:
		return "ipv6_gateway"
	case GlobalIPv6Address:
		return "global_ipv6_address"
	case GlobalIPv6PrefixLen:
		return "global_ipv6_prefix_len"
	case Platform:
		return "platform"
	case Architecture:
		return "architecture"
	case OS:
		return "os"
	case OSVersion:
		return "os_version"
	case OSFeatures:
		return "os_features"
	case Variant:
		return "variant"
	default:
		return "unknown"
	}
}
