// package compose is a wrapper around the docker compose cli.
// it provides a programmatic interface for managing docker compose projects.
// it is not a complete 1:1 mapping of the docker compose cli, but provides a programmatic interface for managing
// docker compose projects in Go, making it useful for automation, tooling, or embedding docker compose behavior.
package compose

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/aptd3v/go-contain/pkg/create"
	"github.com/compose-spec/compose-go/v2/types"
)

type compose struct {
	project *create.Project
}

func NewCompose(project *create.Project) *compose {
	return &compose{
		project: project,
	}
}

// SetComposeUpOption is a function that sets a ComposeUpOptions
type SetComposeUpOption func(*ComposeUpOptions) error

// SetComposeDownOption is a function that sets a ComposeDownOptions
type SetComposeDownOption func(*ComposeDownOptions) error

// SetComposeLogsOption is a function that sets a ComposeLogsOptions
type SetComposeLogsOption func(*ComposeLogsOptions) error

// SetComposeKillOption is a function that sets a ComposeKillOptions
type SetComposeKillOption func(*ComposeKillOptions) error

// Events is a function that runs the docker compose events command
// it returns a channel of Events and an error if the command fails
// it also returns nil if the context is canceled
func (c *compose) Events(ctx context.Context, service string) (<-chan Events, <-chan error, error) {
	eventsCh := make(chan Events, 1)
	errCh := make(chan error, 1)

	writer := newEventsWriter(ctx, eventsCh, errCh)
	cmd, err := c.command(ctx, writer, []string{"events", "--json", service})
	if err != nil {
		return nil, nil, NewComposeEventsError(err)
	}

	go func() {
		defer close(eventsCh)
		defer close(errCh)

		if err := handleContextCancellation(ctx, cmd.Run()); err != nil {
			errCh <- NewComposeEventsError(err)
		}
	}()

	return eventsCh, errCh, nil
}

// Kill is a function that runs the docker compose kill command
// it returns an error if the command fails, or nil if the command succeeds
// it also returns nil if the context is canceled
func (c *compose) Kill(ctx context.Context, setters ...SetComposeKillOption) error {
	opt := &ComposeKillOptions{
		Flags:  []string{"kill"},
		Writer: os.Stdout,
	}
	for _, setter := range setters {
		if err := setter(opt); err != nil {
			return NewComposeKillError(err)
		}
	}
	flags, err := opt.GenerateFlags()
	if err != nil {
		return NewComposeKillError(err)
	}
	cmd, err := c.command(ctx, opt.Writer, flags, opt.Profiles...)
	if err != nil {
		return NewComposeKillError(err)
	}
	return handleContextCancellation(ctx, cmd.Run())
}

func (c *compose) Up(ctx context.Context, setters ...SetComposeUpOption) error {
	opt := &ComposeUpOptions{
		//default flags
		Flags:  []string{"up"},
		Writer: os.Stdout,
	}
	for _, setter := range setters {
		if err := setter(opt); err != nil {
			return NewComposeUpError(err)
		}
	}
	flags, err := opt.GenerateFlags()
	if err != nil {
		return NewComposeUpError(err)
	}
	cmd, err := c.command(ctx, opt.Writer, flags, opt.Profiles...)
	if err != nil {
		return NewComposeUpError(err)
	}
	return handleContextCancellation(ctx, cmd.Run())
}

func (c *compose) Down(ctx context.Context, setters ...SetComposeDownOption) error {
	opt := &ComposeDownOptions{
		Flags:  []string{"down"},
		Writer: os.Stdout,
	}
	for _, setter := range setters {
		if err := setter(opt); err != nil {
			return NewComposeDownError(err)
		}
	}
	flags, err := opt.GenerateFlags()
	if err != nil {
		return NewComposeDownError(err)
	}

	cmd, err := c.command(ctx, opt.Writer, flags)
	if err != nil {
		return NewComposeDownError(err)
	}
	return handleContextCancellation(ctx, cmd.Run())
}

// Logs is a function that runs the docker compose logs command
// it returns an error if the command fails, or nil if the command succeeds
// it also returns nil if the context is canceled
func (c *compose) Logs(ctx context.Context, setters ...SetComposeLogsOption) error {
	opt := &ComposeLogsOptions{
		Flags:  []string{"logs"},
		Writer: os.Stdout,
	}
	for _, setter := range setters {
		if err := setter(opt); err != nil {
			return NewComposeLogsError(err)
		}
	}
	flags, err := opt.GenerateFlags()
	if err != nil {
		return NewComposeLogsError(err)
	}

	cmd, err := c.command(ctx, opt.Writer, flags, opt.Profiles...)
	if err != nil {
		return NewComposeLogsError(err)
	}

	return handleContextCancellation(ctx, cmd.Run())
}
func (c *compose) command(ctx context.Context, writer io.Writer, args []string, profiles ...string) (*exec.Cmd, error) {
	file, err := c.project.Marshal()
	if err != nil {
		return nil, NewComposeError(err)
	}
	base := []string{"compose"}

	// if profiles args are passed, we need to check if any services have profiles
	// if they do not, we need to return an error
	if profiles != nil {
		pmap := map[string]struct{}{}
		c.project.ForEachService(func(name string, service *types.ServiceConfig) error {
			if len(service.Profiles) > 0 {
				for _, profile := range service.Profiles {
					pmap[profile] = struct{}{}
				}
			}
			return nil
		})
		if len(pmap) == 0 {
			return nil, NewComposeError(fmt.Errorf("when using %s.WithProfiles, you must create a service with a profile via sc.WithProfiles", args[0]))
		}
	}
	// if no profiles args are passed, we need to check if any services have profiles
	// if they do, we need to return an error
	if len(profiles) == 0 {
		pmap := map[string]struct{}{}

		c.project.ForEachService(func(name string, service *types.ServiceConfig) error {
			if len(service.Profiles) > 0 {
				for _, profile := range service.Profiles {
					pmap[profile] = struct{}{}
				}
			}
			return nil
		})
		if len(pmap) > 0 {
			names := []string{}
			for profile := range pmap {
				names = append(names, profile)
			}
			return nil, NewComposeError(fmt.Errorf("when you have profiles, you must specify them via %s.WithProfiles\nprofiles found: %s", args[0], strings.Join(names, ", ")))
		}
	}
	// if profiles are set, add them to the base command
	if len(profiles) > 0 {
		for _, profile := range profiles {
			base = append(base, "--profile", profile)
		}
	}

	// for file passed via stdin, we need to add the -f flag
	base = append(base, "-f", "-")
	cmd := exec.CommandContext(ctx, "docker", base...)
	cmd.Args = append(cmd.Args, args...)
	cmd.Stdin = strings.NewReader(string(file))
	cmd.Stdout = writer
	cmd.Stderr = writer
	return cmd, nil
}

func handleContextCancellation(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) {
		return nil
	}
	if ctx.Err() == context.Canceled {
		return nil
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return nil
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil
	}
	return err
}
