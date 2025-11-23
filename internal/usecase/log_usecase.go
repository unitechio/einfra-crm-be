package usecase

import (
	"context"
	"fmt"

	"github.com/unitechio/einfra-be/pkg/docker"
	"github.com/unitechio/einfra-be/pkg/logstream"
)

// LogUsecase handles log operations
type LogUsecase interface {
	StreamContainerLogs(ctx context.Context, containerID string, tail string) (<-chan logstream.LogMessage, <-chan error, error)
}

type logUsecase struct {
	dockerClient *docker.Client
}

// NewLogUsecase creates a new log usecase
func NewLogUsecase(dockerClient *docker.Client) LogUsecase {
	return &logUsecase{
		dockerClient: dockerClient,
	}
}

// StreamContainerLogs streams logs from a container
func (u *logUsecase) StreamContainerLogs(ctx context.Context, containerID string, tail string) (<-chan logstream.LogMessage, <-chan error, error) {
	if containerID == "" {
		return nil, nil, fmt.Errorf("container ID is required")
	}

	if tail == "" {
		tail = "100"
	}

	return u.dockerClient.ContainerLogsStream(ctx, containerID, tail, true)
}
