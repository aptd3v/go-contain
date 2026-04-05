package errdefs

import "errors"

// IsContainerConfigError returns true if the error is a [ErrContainerConfig]
func IsContainerConfigError(err error) bool {
	return errors.Is(err, ErrContainerConfig)
}

// IsHostConfigError returns true if the error is a [ErrHostConfig]
func IsHostConfigError(err error) bool {
	return errors.Is(err, ErrHostConfig)
}

// IsNetworkConfigError returns true if the error is a [ErrNetworkConfig]
func IsNetworkConfigError(err error) bool {
	return errors.Is(err, ErrNetworkConfig)
}

// IsPlatformConfigError returns true if the error is a [ErrPlatformConfig]
func IsPlatformConfigError(err error) bool {
	return errors.Is(err, ErrPlatformConfig)
}
