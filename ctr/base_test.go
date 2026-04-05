package ctr

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/aptd3v/containerkit/config"
	"github.com/aptd3v/containerkit/errdefs"
	"github.com/aptd3v/containerkit/fields"
	"github.com/moby/moby/api/types/network"
	"github.com/stretchr/testify/assert"
)

var (
	nilStringSlice []string
	nilMapStruct   map[string]struct{}
	port80         = network.MustParsePort("80/tcp")
)

func TestBaseAssignments(t *testing.T) {
	_, with := New()
	tests := []struct {
		config   *BaseConfig
		setFn    config.SetBaseConfig
		field    string
		expected any
		wantErr  bool
		message  string
	}{
		{
			config: &BaseConfig{},
			setFn:  with.Base.HealthConfig(),
			field:  "Healthcheck",
			expected: &HealthConfig{
				Test: []string{"NONE"},
			},
			wantErr: false,
			message: "With HealthCheck Disabled",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.Env("TEST", "test"),
			field:    "Env",
			expected: []string{"TEST=test"},
			wantErr:  false,
			message:  "With Env",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.Env("", ""),
			field:    "Env",
			expected: nilStringSlice,
			wantErr:  true,
			message:  "With Empty Env Key",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.Env("TEST", ""),
			field:    "Env",
			expected: []string{"TEST=\"\""},
			wantErr:  false,
			message:  "With Empty Env Value",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.EnvMap(map[string]string{"TEST": "test"}),
			field:    "Env",
			expected: []string{"TEST=test"},
			wantErr:  false,
			message:  "With EnvMap",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.EnvMap(map[string]string{"": ""}),
			field:    "Env",
			expected: make([]string, 0),
			wantErr:  true,
			message:  "With Empty EnvMap",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.EnvMap(map[string]string{"TEST": ""}),
			field:    "Env",
			expected: []string{"TEST=\"\""},
			wantErr:  false,
			message:  "With EnvMap with empty value",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.ExposedPort("80/tcp"),
			field:    "ExposedPorts",
			expected: network.PortSet{port80: {}},
			wantErr:  false,
			message:  "With ExposedPort 80/tcp",
		},

		{
			config:   &BaseConfig{},
			setFn:    with.Base.ExposedPort("300000"),
			field:    "ExposedPorts",
			expected: network.PortSet{},
			wantErr:  true,
			message:  "With Invalid ExposedPort",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.ExposedPort("80"),
			field:    "ExposedPorts",
			expected: network.PortSet{port80: {}},
			wantErr:  false,
			message:  "With ExposedPort 80 with default protocol tcp",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.Hostname("test"),
			field:    "Hostname",
			expected: "test",
			wantErr:  false,
			message:  "With Hostname: test string",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.Hostname(""),
			field:    "Hostname",
			expected: "",
			wantErr:  true,
			message:  "With Hostname: empty string",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.Domainname("test"),
			field:    "Domainname",
			expected: "test",
			wantErr:  false,
			message:  "With Domainname: test string",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.Domainname(""),
			field:    "Domainname",
			expected: "",
			wantErr:  true,
			message:  "With Domainname: empty string string",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.Image("alpine:latest"),
			field:    "Image",
			expected: "alpine:latest",
			wantErr:  false,
			message:  "With Image: alpine:latest string",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.Image(""),
			field:    "Image",
			expected: "",
			wantErr:  true,
			message:  "With Image: empty string",
		},

		{
			config:   &BaseConfig{},
			setFn:    with.Base.Command("curl", "-f", "http://localhost:8080/health"),
			field:    "Cmd",
			expected: []string{"curl", "-f", "http://localhost:8080/health"},
			wantErr:  false,
			message:  "With Command: curl -f http://localhost:8080/health",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.Command("", ""),
			field:    "Cmd",
			expected: nilStringSlice,
			wantErr:  true,
			message:  "With Command: empty command",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.User("test"),
			field:    "User",
			expected: "test",
			wantErr:  false,
			message:  "WithUser ok",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.User(""),
			field:    "User",
			expected: "",
			wantErr:  true,
			message:  "WithUser error",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.AttachedStdin(true),
			field:    "AttachStdin",
			expected: true,
			wantErr:  false,
			message:  "WithAttachedStdin ok",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.AttachedStdin(false),
			field:    "AttachStdin",
			expected: false,
			wantErr:  false,
			message:  "WithAttachedStdin true",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.AttachedStdout(true),
			field:    "AttachStdout",
			expected: true,
			wantErr:  false,
			message:  "WithAttachedStdout true",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.AttachedStderr(true),
			field:    "AttachStderr",
			expected: true,
			wantErr:  false,
			message:  "WithAttachedStderr true",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.Tty(true),
			field:    "Tty",
			expected: true,
			wantErr:  false,
			message:  "WithTty true",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.StdinOpen(true),
			field:    "OpenStdin",
			expected: true,
			wantErr:  false,
			message:  "WithStdinOpen ok",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.StdinOpen(false),
			field:    "OpenStdin",
			expected: false,
			wantErr:  false,
			message:  "WithStdinOpen true",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.StdinOnce(true),
			field:    "StdinOnce",
			expected: true,
			wantErr:  false,
			message:  "WithStdinOnce true",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.EscapedArgs(true),
			field:    "ArgsEscaped",
			expected: true,
			wantErr:  false,
			message:  "WithEscapedArgs true",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.Volume("test"),
			field:    "Volumes",
			expected: map[string]struct{}{"test": {}},
			wantErr:  false,
			message:  "WithVolume test string",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.Volume(""),
			field:    "Volumes",
			expected: nilMapStruct,
			wantErr:  true,
			message:  "WithVolume empty string",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.WorkingDir("test"),
			field:    "WorkingDir",
			expected: "test",
			wantErr:  false,
			message:  "WithWorkingDir test string",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.WorkingDir(""),
			field:    "WorkingDir",
			expected: "",
			wantErr:  true,
			message:  "WithWorkingDir empty string",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.DisabledNetwork(true),
			field:    "NetworkDisabled",
			expected: true,
			wantErr:  false,
			message:  "WithDisabledNetwork ok",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.OnBuild("apk", "add", "--no-cache", "curl"),
			field:    "OnBuild",
			expected: []string{"apk", "add", "--no-cache", "curl"},
			wantErr:  false,
			message:  "WithOnBuild apk add --no-cache curl",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.OnBuild(),
			field:    "OnBuild",
			expected: nilStringSlice,
			wantErr:  true,
			message:  "WithOnBuild empty onbuild",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.Label("com.example.version", "1.0.0"),
			field:    "Labels",
			expected: map[string]string{"com.example.version": "1.0.0"},
			wantErr:  false,
			message:  "WithLabel com.example.version 1.0.0",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.Label("", "1.0.0"),
			field:    "Labels",
			expected: make(map[string]string),
			wantErr:  true,
			message:  "WithLabel empty label key",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.LabelMap(map[string]string{}),
			field:    "Labels",
			expected: map[string]string{},
			wantErr:  true,
			message:  "WithLabelMap empty map",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.LabelMap(map[string]string{"": "1.0.0"}),
			field:    "Labels",
			expected: map[string]string{},
			wantErr:  true,
			message:  "WithLabelMap empty label key",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.LabelMap(map[string]string{"com.example.version": "1.0.0"}),
			field:    "Labels",
			expected: map[string]string{"com.example.version": "1.0.0"},
			wantErr:  false,
			message:  "WithLabelMap com.example.version 1.0.0",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.Label("com.example.version", ""),
			field:    "Labels",
			expected: map[string]string{"com.example.version": ""},
			wantErr:  false,
			message:  "WithLabel com.example.version empty label value",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.StopSignal("SIGKILL"),
			field:    "StopSignal",
			expected: "SIGKILL",
			wantErr:  false,
			message:  "WithStopSignal SIGKILL",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.StopSignal(""),
			field:    "StopSignal",
			expected: "",
			wantErr:  true,
			message:  "WithStopSignal empty signal",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.Entrypoint("doker-entrypoint.sh"),
			field:    "Entrypoint",
			expected: []string{"doker-entrypoint.sh"},
			wantErr:  false,
			message:  "WithEntrypoint doker-entrypoint.sh",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.Entrypoint(),
			field:    "Entrypoint",
			expected: nilStringSlice,
			wantErr:  true,
			message:  "WithEntrypoint empty entrypoint",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.Shell("sh", "-c"),
			field:    "Shell",
			expected: []string{"sh", "-c"},
			wantErr:  false,
			message:  "WithShell sh -c",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.Shell(),
			field:    "Shell",
			expected: nilStringSlice,
			wantErr:  true,
			message:  "WithShell empty shell",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.StopTimeout(10 * time.Second),
			field:    "StopTimeout",
			expected: func() *int { v := 10; return &v }(),
			wantErr:  false,
			message:  "WithStopTimeout 10 seconds",
		},

		{
			config:   &BaseConfig{},
			setFn:    with.Base.Fail(fields.Cmd, errors.New("test")),
			field:    "",
			expected: nil,
			wantErr:  true,
			message:  "Fail ok",
		},
		{
			config:   &BaseConfig{},
			setFn:    with.Base.Failf(fields.Cmd, "test-%s", "test"),
			field:    "",
			expected: nil,
			wantErr:  true,
			message:  "Failf ok",
		},
		{
			config: &BaseConfig{},
			setFn: with.Base.HealthConfig(
				nil,
				with.Base.Health.Interval(10*time.Second),
				with.Base.Health.Test("CMD-SHELL", "curl -f http://localhost:8080/health || exit 1"),
				with.Base.Health.Timeout(10*time.Second),
				with.Base.Health.StartPeriod(10*time.Second),
				with.Base.Health.Retries(10),
			),
			field: "Healthcheck",
			expected: &HealthConfig{
				Interval:    10 * time.Second,
				Retries:     10,
				StartPeriod: 10 * time.Second,
				Test:        []string{"CMD-SHELL", "curl -f http://localhost:8080/health || exit 1"},
				Timeout:     10 * time.Second,
			},
			wantErr: false,
			message: "WithHealthCheck ok",
		},
		{
			config: &BaseConfig{},
			setFn: with.Base.HealthConfig(
				with.Base.Health.Fail(
					fields.HealthConfig,
					errors.New("test"),
				),
			),
			field: "Healthcheck",
			expected: &HealthConfig{
				Test: nilStringSlice,
			},
			wantErr: true,
			message: "WithHealthCheck fail",
		},
		{
			config: &BaseConfig{},
			setFn: with.Base.HealthConfig(
				with.Base.Health.Failf(
					fields.HealthConfig,
					"error-%s", "test",
				),
			),
			field: "Healthcheck",
			expected: &HealthConfig{
				Test: nilStringSlice,
			},
			wantErr: true,
			message: "WithHealthCheck failf",
		},

		{
			config: &BaseConfig{},
			setFn: with.Base.HealthConfig(
				with.Base.Health.Test(),
				with.Base.Health.Interval(10*time.Second),
				with.Base.Health.Timeout(10*time.Second),
				with.Base.Health.StartPeriod(10*time.Second),
				with.Base.Health.Retries(10),
			),
			field: "Healthcheck",
			expected: &HealthConfig{
				Test:        []string{"NONE"},
				Interval:    10 * time.Second,
				Timeout:     10 * time.Second,
				StartPeriod: 10 * time.Second,
				Retries:     10,
			},
			wantErr: false,
			message: "WithHealthCheck ok with empty test",
		},
	}

	for _, test := range tests {
		err := test.setFn(test.config)
		if test.wantErr {
			assert.Error(t, err)
			assert.True(t,
				errdefs.IsContainerConfigError(err),
				"expected container config error  for %s", test.message,
			)
		} else {
			assert.NoError(t, err)
		}
		if test.field != "" {
			assert.Equal(t,
				test.expected,
				reflect.ValueOf(*test.config).FieldByName(test.field).Interface(),
				test.message,
			)
		}
	}
}
func TestBaseConfigBasic(t *testing.T) {
	container, with := New()
	container.WithBaseConfig(
		with.Base.Hostname("com.example.container"),
		nil, // test nil setter no panic
		with.Base.Fail(fields.Unknown, errors.New("should append 1 error")),
		with.Base.HealthConfig(
			with.Base.Health.Test("CMD-SHELL", "curl -f http://localhost:8080/health || exit 1"),
			with.Base.Health.Interval(10*time.Second),
			with.Base.Health.Timeout(10*time.Second),
			with.Base.Health.StartPeriod(10*time.Second),
			with.Base.Health.Retries(10),
		),
	)

	assert.Equal(t, "com.example.container", container.BaseConfig.Hostname)
	assert.Equal(t, []string{"CMD-SHELL", "curl -f http://localhost:8080/health || exit 1"}, container.BaseConfig.Healthcheck.Test)
	assert.Equal(t, 10*time.Second, container.BaseConfig.Healthcheck.Interval)
	assert.Equal(t, 10*time.Second, container.BaseConfig.Healthcheck.Timeout)
	assert.Equal(t, 10*time.Second, container.BaseConfig.Healthcheck.StartPeriod)
	assert.Equal(t, 10, container.BaseConfig.Healthcheck.Retries)
	assert.Equal(t, 1, len(container.Errors()), "expected 1 error")

	container.ErrorsEach(func(err error) error {
		assert.Error(t, err)
		assert.True(t, errdefs.IsContainerConfigError(err), "expected container config error for %s", err.Error())
		return nil
	})
	//nil callback no panic
	err := container.ErrorsEach(nil)
	assert.NoError(t, err)

	container.ErrorsEach(func(err error) error {
		assert.Error(t, err)
		assert.True(t, errdefs.IsContainerConfigError(err), "expected container config error for %s", err.Error())
		return nil
	})
	// single error returned

	err = container.ErrorsEach(func(_ error) error {
		return errors.New("should be error")
	})
	assert.Error(t, err)

}
