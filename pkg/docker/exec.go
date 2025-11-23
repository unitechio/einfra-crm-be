package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

// ExecConfig represents configuration for container exec
type ExecConfig struct {
	ContainerID  string   `json:"container_id"`
	Cmd          []string `json:"cmd"`
	AttachStdin  bool     `json:"attach_stdin"`
	AttachStdout bool     `json:"attach_stdout"`
	AttachStderr bool     `json:"attach_stderr"`
	Tty          bool     `json:"tty"`
	Env          []string `json:"env,omitempty"`
	WorkingDir   string   `json:"working_dir,omitempty"`
	User         string   `json:"user,omitempty"`
	Privileged   bool     `json:"privileged,omitempty"`
}

// ExecStartConfig represents configuration for starting exec
type ExecStartConfig struct {
	Detach bool `json:"detach"`
	Tty    bool `json:"tty"`
}

// ExecResult represents the result of exec operation
type ExecResult struct {
	ExecID   string
	ExitCode int
	Output   string
	Error    string
}

// ContainerExec creates an exec instance in a container
func (c *Client) ContainerExec(ctx context.Context, config ExecConfig) (string, error) {
	execConfig := container.ExecOptions{
		AttachStdin:  config.AttachStdin,
		AttachStdout: config.AttachStdout,
		AttachStderr: config.AttachStderr,
		Tty:          config.Tty,
		Cmd:          config.Cmd,
		Env:          config.Env,
		WorkingDir:   config.WorkingDir,
		User:         config.User,
		Privileged:   config.Privileged,
	}

	resp, err := c.cli.ContainerExecCreate(ctx, config.ContainerID, execConfig)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

// ExecStart starts an exec instance
func (c *Client) ExecStart(ctx context.Context, execID string, config ExecStartConfig) (io.ReadCloser, error) {
	startConfig := types.ExecStartCheck{
		Detach: config.Detach,
		Tty:    config.Tty,
	}

	resp, err := c.cli.ContainerExecAttach(ctx, execID, startConfig)
	if err != nil {
		return nil, err
	}

	return resp.Reader, nil
}

// ExecInspect returns information about an exec instance
func (c *Client) ExecInspect(ctx context.Context, execID string) (container.ExecInspect, error) {
	return c.cli.ContainerExecInspect(ctx, execID)
}

// ExecResize resizes the TTY of an exec instance
func (c *Client) ExecResize(ctx context.Context, execID string, height, width uint) error {
	options := container.ResizeOptions{
		Height: height,
		Width:  width,
	}

	return c.cli.ContainerExecResize(ctx, execID, options)
}

// ContainerExecSimple executes a simple command and returns output
func (c *Client) ContainerExecSimple(ctx context.Context, containerID string, cmd []string) (*ExecResult, error) {
	// Create exec instance
	execConfig := ExecConfig{
		ContainerID:  containerID,
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	execID, err := c.ContainerExec(ctx, execConfig)
	if err != nil {
		return nil, err
	}

	// Start exec
	startConfig := ExecStartConfig{
		Detach: false,
	}

	reader, err := c.ExecStart(ctx, execID, startConfig)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// Read output
	output, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Get exit code
	inspect, err := c.ExecInspect(ctx, execID)
	if err != nil {
		return nil, err
	}

	result := &ExecResult{
		ExecID:   execID,
		ExitCode: inspect.ExitCode,
		Output:   string(output),
	}

	return result, nil
}
