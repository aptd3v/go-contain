package cli

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/aptd3v/containerkit/config"
	"github.com/moby/moby/client"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

type Setters struct {
	ContainerAttach *attachSetters
	ImagePull       *imagePullSetters
}

type Cli struct {
	client *client.Client
	errs   []error
	Set    *Setters
}

func wrapVoid[T any](fn func(*T)) config.SetConfig[T] {
	return func(cfg *T) error {
		if fn == nil {
			return nil
		}
		fn(cfg)
		return nil
	}
}

func createSetters() *Setters {
	return &Setters{
		ContainerAttach: &attachSetters{},
		ImagePull:       &imagePullSetters{},
	}
}

func NewCli(setters ...config.SetClientOption) (c *Cli, with *Setters, err error) {
	opt := []client.Opt{}
	for _, setter := range setters {
		if setter == nil {
			continue
		}
		opt = append(opt, setter())
	}
	cli, err := client.New(opt...)
	if err != nil {
		return nil, nil, err
	}
	with = createSetters()
	return &Cli{client: cli, Set: with}, with, nil
}

// WithHTTPClient overrides the client's HTTP client with the specified one.
func WithHTTPClient(httpClient *http.Client) config.SetClientOption {
	return func() client.Opt {
		return client.WithHTTPClient(httpClient)
	}
}

// FromEnv configures the client with values from environment variables. It
// is the equivalent of using the [WithTLSClientConfigFromEnv], [WithHostFromEnv],
// and [WithVersionFromEnv] options.
//
// FromEnv uses the following environment variables:
//
//   - DOCKER_HOST ([EnvOverrideHost]) to set the URL to the docker server.
//   - DOCKER_API_VERSION ([EnvOverrideAPIVersion]) to set the version of the
//     API to use, leave empty for latest.
//   - DOCKER_CERT_PATH ([EnvOverrideCertPath]) to specify the directory from
//     which to load the TLS certificates ("ca.pem", "cert.pem", "key.pem').
//   - DOCKER_TLS_VERIFY ([EnvTLSVerify]) to enable or disable TLS verification
//     (off by default).
func FromEnv() config.SetClientOption {
	return func() client.Opt {
		return client.FromEnv
	}
}

// WithDialContext applies the dialer to the client transport. This can be
// used to set the Timeout and KeepAlive settings of the client. It returns
// an error if the client does not have a [http.Transport] configured.
func WithDialContext(dialContext func(ctx context.Context, network, addr string) (net.Conn, error)) config.SetClientOption {
	return func() client.Opt {
		return client.WithDialContext(dialContext)
	}
}

// WithHost overrides the client host with the specified one.
func WithHost(host string) config.SetClientOption {
	return func() client.Opt {
		return client.WithHost(host)
	}
}

// WithHostFromEnv overrides the client host with the host specified in the
// DOCKER_HOST ([EnvOverrideHost]) environment variable. If DOCKER_HOST is not set,
// or set to an empty value, the host is not modified.
func WithHostFromEnv() config.SetClientOption {
	return func() client.Opt {
		return client.WithHostFromEnv()
	}
}

// WithTimeout configures the time limit for requests made by the HTTP client.
func WithTimeout(timeout time.Duration) config.SetClientOption {
	return func() client.Opt {
		return client.WithTimeout(timeout)
	}
}

// WithUserAgent configures the User-Agent header to use for HTTP requests.
// It overrides any User-Agent set in headers. When set to an empty string,
// the User-Agent header is removed, and no header is sent.
func WithUserAgent(userAgent string) config.SetClientOption {
	return func() client.Opt {
		return client.WithUserAgent(userAgent)
	}
}

// WithHTTPHeaders appends custom HTTP headers to the client's default headers.
// It does not allow for built-in headers (such as "User-Agent", if set) to
// be overridden. Also see [WithUserAgent].
func WithHTTPHeaders(headers map[string]string) config.SetClientOption {
	return func() client.Opt {
		return client.WithHTTPHeaders(headers)
	}
}

// WithScheme overrides the client scheme with the specified one.
func WithScheme(scheme string) config.SetClientOption {
	return func() client.Opt {
		return client.WithScheme(scheme)
	}
}

// WithTLSClientConfig applies a TLS config to the client transport.
func WithTLSClientConfig(cacertPath, certPath, keyPath string) config.SetClientOption {
	return func() client.Opt {
		return client.WithTLSClientConfig(cacertPath, certPath, keyPath)

	}
}

// WithTLSClientConfigFromEnv configures the client's TLS settings with the
// settings in the DOCKER_CERT_PATH ([EnvOverrideCertPath]) and DOCKER_TLS_VERIFY
// ([EnvTLSVerify]) environment variables. If DOCKER_CERT_PATH is not set or empty,
// TLS configuration is not modified.
//
// WithTLSClientConfigFromEnv uses the following environment variables:
//
//   - DOCKER_CERT_PATH ([EnvOverrideCertPath]) to specify the directory from
//     which to load the TLS certificates ("ca.pem", "cert.pem", "key.pem").
//   - DOCKER_TLS_VERIFY ([EnvTLSVerify]) to enable or disable TLS verification
//     (off by default).
func WithTLSClientConfigFromEnv() config.SetClientOption {
	return func() client.Opt {
		return client.WithTLSClientConfigFromEnv()
	}
}

// WithTraceProvider sets the trace provider for the client.
// If this is not set then the global trace provider will be used.
func WithTraceProvider(provider trace.TracerProvider) config.SetClientOption {
	return func() client.Opt {
		return client.WithTraceProvider(provider)
	}
}

// WithTraceOptions sets tracing span options for the client.
func WithTraceOptions(opts ...otelhttp.Option) config.SetClientOption {
	return func() client.Opt {
		return client.WithTraceOptions(opts...)
	}
}
