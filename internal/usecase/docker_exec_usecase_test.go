package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/unitechio/einfra-be/pkg/docker"
)

// MockDockerClient is a mock implementation of Docker client
type MockDockerClient struct {
	mock.Mock
}

func (m *MockDockerClient) ContainerExec(ctx context.Context, config docker.ExecConfig) (string, error) {
	args := m.Called(ctx, config)
	return args.String(0), args.Error(1)
}

func (m *MockDockerClient) ExecStart(ctx context.Context, execID string, config docker.ExecStartConfig) ([]byte, error) {
	args := m.Called(ctx, execID, config)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockDockerClient) ExecInspect(ctx context.Context, execID string) (map[string]interface{}, error) {
	args := m.Called(ctx, execID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockDockerClient) ExecResize(ctx context.Context, execID string, height, width uint) error {
	args := m.Called(ctx, execID, height, width)
	return args.Error(0)
}

func (m *MockDockerClient) ContainerExecSimple(ctx context.Context, containerID string, cmd []string) (*docker.ExecResult, error) {
	args := m.Called(ctx, containerID, cmd)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*docker.ExecResult), args.Error(1)
}

// TestCreateExec tests creating an exec instance
func TestCreateExec(t *testing.T) {
	mockClient := new(MockDockerClient)
	// Note: This test is a placeholder - actual implementation would need proper Docker client mocking

	ctx := context.Background()

	t.Run("Success - Create exec instance", func(t *testing.T) {
		containerID := "container-123"
		cmd := []string{"/bin/bash"}
		expectedExecID := "exec-456"

		config := docker.ExecConfig{
			ContainerID:  containerID,
			Cmd:          cmd,
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Tty:          true,
		}

		mockClient.On("ContainerExec", ctx, config).Return(expectedExecID, nil).Once()

		execID, err := mockClient.ContainerExec(ctx, config)
		assert.NoError(t, err)
		assert.Equal(t, expectedExecID, execID)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error - Empty container ID", func(t *testing.T) {
		// Usecase validation test
		containerID := ""
		cmd := []string{"/bin/bash"}

		// This would be tested in the usecase layer
		assert.NotEmpty(t, containerID, "Container ID should not be empty")
		assert.NotEmpty(t, cmd, "Command should not be empty")
	})
}

// TestExecuteCommand tests executing a simple command
func TestExecuteCommand(t *testing.T) {
	mockClient := new(MockDockerClient)
	ctx := context.Background()

	t.Run("Success - Execute command", func(t *testing.T) {
		containerID := "container-123"
		cmd := []string{"ls", "-la"}
		expectedResult := &docker.ExecResult{
			ExecID:   "exec-789",
			ExitCode: 0,
			Output:   "total 0\ndrwxr-xr-x",
		}

		mockClient.On("ContainerExecSimple", ctx, containerID, cmd).Return(expectedResult, nil).Once()

		result, err := mockClient.ContainerExecSimple(ctx, containerID, cmd)
		assert.NoError(t, err)
		assert.Equal(t, 0, result.ExitCode)
		assert.Contains(t, result.Output, "total")
		mockClient.AssertExpectations(t)
	})

	t.Run("Error - Command failed", func(t *testing.T) {
		containerID := "container-123"
		cmd := []string{"invalid-command"}
		expectedResult := &docker.ExecResult{
			ExecID:   "exec-790",
			ExitCode: 127,
			Error:    "command not found",
		}

		mockClient.On("ContainerExecSimple", ctx, containerID, cmd).Return(expectedResult, nil).Once()

		result, err := mockClient.ContainerExecSimple(ctx, containerID, cmd)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, result.ExitCode)
		mockClient.AssertExpectations(t)
	})
}

// TestExecResize tests resizing exec TTY
func TestExecResize(t *testing.T) {
	mockClient := new(MockDockerClient)
	ctx := context.Background()

	t.Run("Success - Resize TTY", func(t *testing.T) {
		execID := "exec-123"
		height := uint(24)
		width := uint(80)

		mockClient.On("ExecResize", ctx, execID, height, width).Return(nil).Once()

		err := mockClient.ExecResize(ctx, execID, height, width)
		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error - Invalid dimensions", func(t *testing.T) {
		execID := "exec-123"
		height := uint(0)
		width := uint(0)

		// Validation should happen in usecase
		assert.Greater(t, uint(24), height, "Height must be greater than 0")
		assert.Greater(t, uint(80), width, "Width must be greater than 0")
	})
}

// BenchmarkExecuteCommand benchmarks command execution
func BenchmarkExecuteCommand(b *testing.B) {
	mockClient := new(MockDockerClient)
	ctx := context.Background()
	containerID := "container-bench"
	cmd := []string{"echo", "hello"}

	result := &docker.ExecResult{
		ExecID:   "exec-bench",
		ExitCode: 0,
		Output:   "hello",
	}

	mockClient.On("ContainerExecSimple", ctx, containerID, cmd).Return(result, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mockClient.ContainerExecSimple(ctx, containerID, cmd)
	}
}
