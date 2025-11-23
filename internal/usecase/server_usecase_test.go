package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/unitechio/einfra-be/internal/domain"
	"github.com/unitechio/einfra-be/internal/usecase"
	"github.com/unitechio/einfra-be/pkg/ssh"
)

// MockServerRepository is a mock implementation of ServerRepository
type MockServerRepository struct {
	mock.Mock
}

func (m *MockServerRepository) Create(ctx context.Context, server *domain.Server) error {
	args := m.Called(ctx, server)
	return args.Error(0)
}

func (m *MockServerRepository) GetByID(ctx context.Context, id string) (*domain.Server, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Server), args.Error(1)
}

func (m *MockServerRepository) GetByIPAddress(ctx context.Context, ip string) (*domain.Server, error) {
	args := m.Called(ctx, ip)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Server), args.Error(1)
}

func (m *MockServerRepository) List(ctx context.Context, filter domain.ServerFilter) ([]*domain.Server, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*domain.Server), args.Get(1).(int64), args.Error(2)
}

func (m *MockServerRepository) Update(ctx context.Context, server *domain.Server) error {
	args := m.Called(ctx, server)
	return args.Error(0)
}

func (m *MockServerRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockServerRepository) UpdateStatus(ctx context.Context, id string, status domain.ServerStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

// TestCreateServer tests server creation
func TestCreateServer(t *testing.T) {
	mockRepo := new(MockServerRepository)
	tunnelManager := ssh.NewTunnelManager()
	uc := usecase.NewServerUsecase(mockRepo, tunnelManager)

	ctx := context.Background()

	t.Run("Success - Create server without tunnel", func(t *testing.T) {
		server := &domain.Server{
			Name:      "test-server",
			IPAddress: "192.168.1.100",
			CPUCores:  4,
			MemoryGB:  8.0,
			DiskGB:    100.0,
			SSHPort:   22,
			SSHUser:   "root",
		}

		mockRepo.On("GetByIPAddress", ctx, server.IPAddress).Return(nil, nil).Once()
		mockRepo.On("Create", ctx, server).Return(nil).Once()

		err := uc.CreateServer(ctx, server)
		assert.NoError(t, err)
		assert.Equal(t, domain.ServerStatusOffline, server.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Create server with tunnel", func(t *testing.T) {
		server := &domain.Server{
			Name:          "private-server",
			IPAddress:     "10.0.1.100",
			CPUCores:      4,
			MemoryGB:      8.0,
			DiskGB:        100.0,
			SSHPort:       22,
			SSHUser:       "ubuntu",
			TunnelEnabled: true,
			TunnelHost:    "bastion.example.com",
			TunnelPort:    22,
			TunnelUser:    "tunnel-user",
		}

		mockRepo.On("GetByIPAddress", ctx, server.IPAddress).Return(nil, nil).Once()
		mockRepo.On("Create", ctx, server).Return(nil).Once()

		err := uc.CreateServer(ctx, server)
		assert.NoError(t, err)
		assert.True(t, server.TunnelEnabled)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Duplicate IP address", func(t *testing.T) {
		server := &domain.Server{
			Name:      "duplicate-server",
			IPAddress: "192.168.1.100",
			CPUCores:  4,
			MemoryGB:  8.0,
			DiskGB:    100.0,
		}

		existingServer := &domain.Server{ID: "existing-id"}
		mockRepo.On("GetByIPAddress", ctx, server.IPAddress).Return(existingServer, nil).Once()

		err := uc.CreateServer(ctx, server)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Missing required fields", func(t *testing.T) {
		server := &domain.Server{
			Name: "incomplete-server",
		}

		err := uc.CreateServer(ctx, server)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "IP address is required")
	})
}

// TestListServers tests server listing
func TestListServers(t *testing.T) {
	mockRepo := new(MockServerRepository)
	tunnelManager := ssh.NewTunnelManager()
	uc := usecase.NewServerUsecase(mockRepo, tunnelManager)

	ctx := context.Background()

	t.Run("Success - List servers with pagination", func(t *testing.T) {
		filter := domain.ServerFilter{
			Page:     1,
			PageSize: 10,
		}

		servers := []*domain.Server{
			{ID: "1", Name: "server-1"},
			{ID: "2", Name: "server-2"},
		}

		mockRepo.On("List", ctx, filter).Return(servers, int64(2), nil).Once()

		result, total, err := uc.ListServers(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, int64(2), total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Default pagination", func(t *testing.T) {
		filter := domain.ServerFilter{}

		expectedFilter := domain.ServerFilter{
			Page:     1,
			PageSize: 20,
		}

		mockRepo.On("List", ctx, expectedFilter).Return([]*domain.Server{}, int64(0), nil).Once()

		_, _, err := uc.ListServers(ctx, filter)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// TestUpdateServer tests server updates
func TestUpdateServer(t *testing.T) {
	mockRepo := new(MockServerRepository)
	tunnelManager := ssh.NewTunnelManager()
	uc := usecase.NewServerUsecase(mockRepo, tunnelManager)

	ctx := context.Background()

	t.Run("Success - Update server", func(t *testing.T) {
		server := &domain.Server{
			ID:       "server-1",
			Name:     "updated-server",
			CPUCores: 8,
		}

		existingServer := &domain.Server{ID: "server-1"}
		mockRepo.On("GetByID", ctx, server.ID).Return(existingServer, nil).Once()
		mockRepo.On("Update", ctx, server).Return(nil).Once()

		err := uc.UpdateServer(ctx, server)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Server not found", func(t *testing.T) {
		server := &domain.Server{
			ID:   "non-existent",
			Name: "test",
		}

		mockRepo.On("GetByID", ctx, server.ID).Return(nil, nil).Once()

		err := uc.UpdateServer(ctx, server)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		mockRepo.AssertExpectations(t)
	})
}

// TestDeleteServer tests server deletion
func TestDeleteServer(t *testing.T) {
	mockRepo := new(MockServerRepository)
	tunnelManager := ssh.NewTunnelManager()
	uc := usecase.NewServerUsecase(mockRepo, tunnelManager)

	ctx := context.Background()

	t.Run("Success - Delete server", func(t *testing.T) {
		serverID := "server-1"

		existingServer := &domain.Server{ID: serverID}
		mockRepo.On("GetByID", ctx, serverID).Return(existingServer, nil).Once()
		mockRepo.On("Delete", ctx, serverID).Return(nil).Once()

		err := uc.DeleteServer(ctx, serverID)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Server not found", func(t *testing.T) {
		serverID := "non-existent"

		mockRepo.On("GetByID", ctx, serverID).Return(nil, nil).Once()

		err := uc.DeleteServer(ctx, serverID)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// TestHealthCheck tests server health check
func TestHealthCheck(t *testing.T) {
	mockRepo := new(MockServerRepository)
	tunnelManager := ssh.NewTunnelManager()
	uc := usecase.NewServerUsecase(mockRepo, tunnelManager)

	ctx := context.Background()

	t.Run("Success - Online server", func(t *testing.T) {
		serverID := "server-1"
		server := &domain.Server{
			ID:     serverID,
			Status: domain.ServerStatusOnline,
		}

		mockRepo.On("GetByID", ctx, serverID).Return(server, nil).Once()
		mockRepo.On("UpdateStatus", ctx, serverID, domain.ServerStatusOnline).Return(nil).Once()

		isHealthy, err := uc.HealthCheck(ctx, serverID)
		assert.NoError(t, err)
		assert.True(t, isHealthy)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Offline server", func(t *testing.T) {
		serverID := "server-2"
		server := &domain.Server{
			ID:     serverID,
			Status: domain.ServerStatusOffline,
		}

		mockRepo.On("GetByID", ctx, serverID).Return(server, nil).Once()
		mockRepo.On("UpdateStatus", ctx, serverID, domain.ServerStatusOffline).Return(nil).Once()

		isHealthy, err := uc.HealthCheck(ctx, serverID)
		assert.NoError(t, err)
		assert.False(t, isHealthy)
		mockRepo.AssertExpectations(t)
	})
}

// Benchmark tests
func BenchmarkCreateServer(b *testing.B) {
	mockRepo := new(MockServerRepository)
	tunnelManager := ssh.NewTunnelManager()
	uc := usecase.NewServerUsecase(mockRepo, tunnelManager)

	ctx := context.Background()
	server := &domain.Server{
		Name:      "bench-server",
		IPAddress: "192.168.1.200",
		CPUCores:  4,
		MemoryGB:  8.0,
		DiskGB:    100.0,
	}

	mockRepo.On("GetByIPAddress", ctx, mock.Anything).Return(nil, nil)
	mockRepo.On("Create", ctx, mock.Anything).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uc.CreateServer(ctx, server)
	}
}
