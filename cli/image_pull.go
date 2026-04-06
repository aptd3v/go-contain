package cli

import (
	"context"

	"github.com/aptd3v/containerkit/config"
	"github.com/moby/moby/client"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type ImagePullResponse = client.ImagePullResponse
type imagePullSetters struct {
	ref *Cli
}

// RegistryAuth sets the registry auth for the image pull options.
func (i *imagePullSetters) RegistryAuth(registryAuth string) config.SetImagePullOptions {
	return wrapVoid(func(opt *client.ImagePullOptions) {
		opt.RegistryAuth = registryAuth
	})
}

// PullAllPlatforms sets the pull all platforms flag for the image pull options.
// Parameter:
//   - all: if true, all platforms will be pulled
func (i *imagePullSetters) PullAllPlatforms(all bool) config.SetImagePullOptions {
	return wrapVoid(func(opt *client.ImagePullOptions) {
		opt.All = all
	})
}

// Platforms sets the platforms for the image pull options.
func (i *imagePullSetters) Platforms(platform ...ocispec.Platform) config.SetImagePullOptions {
	return wrapVoid(func(opt *client.ImagePullOptions) {
		opt.Platforms = append(opt.Platforms, platform...)
	})
}

// Platform appends the platform from the container configuration.
// Parameter:
//   - container: the container to get the platform from
func (i *imagePullSetters) Platform(setters ...config.SetPlatformConfig) config.SetImagePullOptions {
	return func(opt *client.ImagePullOptions) error {
		platform := ocispec.Platform{}
		for _, setter := range setters {
			if setter == nil {
				continue
			}
			if err := setter(&platform); err != nil {
				return err
			}
		}
		opt.Platforms = append(opt.Platforms, platform)
		return nil
	}
}

func (c *Cli) ImagePull(ctx context.Context, image string, options ...config.SetImagePullOptions) (ImagePullResponse, error) {
	opt := client.ImagePullOptions{}
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(&opt); err != nil {
			return nil, err
		}
	}
	return c.client.ImagePull(ctx, image, opt)
}
