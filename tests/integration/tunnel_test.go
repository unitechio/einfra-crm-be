package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unitechio/einfra-be/pkg/ssh"
)

// TestTunnelCreation tests SSH tunnel creation and lifecycle
func TestTunnelCreation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("Create and start tunnel", func(t *testing.T) {
		// This test requires actual SSH infrastructure
		// Skip if not available
		t.Skip("Requires SSH infrastructure")

		tunnelCfg := ssh.TunnelConfig{
			SSHConfig: ssh.Config{
				Host:    "bastion.example.com",
				Port:    22,
				User:    "test-user",
				KeyPath: "/path/to/key.pem",
				Timeout: 10 * time.Second,
			},
			LocalAddr:  "localhost:13306",
			RemoteAddr: "10.0.1.100:3306",
		}

		tunnel, err := ssh.NewTunnel(tunnelCfg)
		require.NoError(t, err)

		err = tunnel.Start()
		require.NoError(t, err)

		// Verify tunnel is active
		assert.True(t, tunnel.IsActive())

		// Get stats
		stats := tunnel.GetStats()
		assert.True(t, stats["active"].(bool))
		assert.Equal(t, "localhost:13306", stats["local_addr"])

		// Stop tunnel
		err = tunnel.Stop()
		assert.NoError(t, err)
		assert.False(t, tunnel.IsActive())
	})
}

// TestTunnelManager tests tunnel manager operations
func TestTunnelManager(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	manager := ssh.NewTunnelManager()

	t.Run("Create multiple tunnels", func(t *testing.T) {
		t.Skip("Requires SSH infrastructure")

		// Create first tunnel
		cfg1 := ssh.TunnelConfig{
			SSHConfig: ssh.Config{
				Host:    "bastion1.example.com",
				Port:    22,
				User:    "user1",
				KeyPath: "/path/to/key1.pem",
			},
			LocalAddr:  "localhost:13307",
			RemoteAddr: "10.0.1.101:3306",
		}

		err := manager.CreateTunnel("tunnel-1", cfg1)
		require.NoError(t, err)

		// Create second tunnel
		cfg2 := ssh.TunnelConfig{
			SSHConfig: ssh.Config{
				Host:    "bastion2.example.com",
				Port:    22,
				User:    "user2",
				KeyPath: "/path/to/key2.pem",
			},
			LocalAddr:  "localhost:13308",
			RemoteAddr: "10.0.1.102:3306",
		}

		err = manager.CreateTunnel("tunnel-2", cfg2)
		require.NoError(t, err)

		// List tunnels
		tunnels := manager.ListTunnels()
		assert.Equal(t, 2, len(tunnels))

		// Stop all tunnels
		err = manager.StopAll()
		assert.NoError(t, err)

		tunnels = manager.ListTunnels()
		assert.Equal(t, 0, len(tunnels))
	})

	t.Run("Prevent duplicate tunnel IDs", func(t *testing.T) {
		t.Skip("Requires SSH infrastructure")

		cfg := ssh.TunnelConfig{
			SSHConfig: ssh.Config{
				Host: "bastion.example.com",
				Port: 22,
				User: "user",
			},
			LocalAddr:  "localhost:13309",
			RemoteAddr: "10.0.1.103:3306",
		}

		err := manager.CreateTunnel("duplicate-test", cfg)
		require.NoError(t, err)

		// Try to create with same ID
		err = manager.CreateTunnel("duplicate-test", cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")

		// Cleanup
		manager.StopTunnel("duplicate-test")
	})
}

// TestSSHCommandExecution tests command execution through SSH
func TestSSHCommandExecution(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("Execute command on remote server", func(t *testing.T) {
		t.Skip("Requires SSH server")

		sshCfg := ssh.Config{
			Host:    "test-server.example.com",
			Port:    22,
			User:    "test-user",
			KeyPath: "/path/to/key.pem",
			Timeout: 10 * time.Second,
		}

		client, err := ssh.NewClient(sshCfg)
		require.NoError(t, err)

		err = client.Connect()
		require.NoError(t, err)
		defer client.Close()

		ctx := context.Background()

		// Execute simple command
		result, err := client.ExecuteCommand(ctx, "echo 'Hello, World!'")
		require.NoError(t, err)
		assert.Equal(t, 0, result.ExitCode)
		assert.Contains(t, result.Stdout, "Hello, World!")
	})
}

// TestTunnelConnectionFlow tests end-to-end tunnel connection flow
func TestTunnelConnectionFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("Connect to private server through tunnel", func(t *testing.T) {
		t.Skip("Requires full SSH infrastructure")

		manager := ssh.NewTunnelManager()

		// Step 1: Create tunnel to bastion
		tunnelCfg := ssh.TunnelConfig{
			SSHConfig: ssh.Config{
				Host:    "bastion.example.com",
				Port:    22,
				User:    "tunnel-user",
				KeyPath: "/path/to/bastion-key.pem",
			},
			LocalAddr:  "localhost:12222",
			RemoteAddr: "10.0.1.100:22", // Private server
		}

		err := manager.CreateTunnel("test-tunnel", tunnelCfg)
		require.NoError(t, err)

		// Step 2: Connect to private server through tunnel
		privateSshCfg := ssh.Config{
			Host:    "localhost",
			Port:    12222, // Tunnel local port
			User:    "ubuntu",
			KeyPath: "/path/to/server-key.pem",
		}

		client, err := ssh.NewClient(privateSshCfg)
		require.NoError(t, err)

		err = client.Connect()
		require.NoError(t, err)

		// Step 3: Execute command on private server
		ctx := context.Background()
		result, err := client.ExecuteCommand(ctx, "hostname")
		require.NoError(t, err)
		assert.Equal(t, 0, result.ExitCode)

		// Cleanup
		client.Close()
		manager.StopTunnel("test-tunnel")
	})
}

// TestTunnelPerformance tests tunnel performance metrics
func TestTunnelPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("Measure tunnel throughput", func(t *testing.T) {
		t.Skip("Requires SSH infrastructure")

		// Create tunnel
		manager := ssh.NewTunnelManager()
		tunnelCfg := ssh.TunnelConfig{
			SSHConfig: ssh.Config{
				Host: "bastion.example.com",
				Port: 22,
				User: "perf-test",
			},
			LocalAddr:  "localhost:15432",
			RemoteAddr: "10.0.1.100:5432",
		}

		err := manager.CreateTunnel("perf-tunnel", tunnelCfg)
		require.NoError(t, err)

		// Simulate multiple connections
		for i := 0; i < 10; i++ {
			go func(id int) {
				// Simulate connection through tunnel
				time.Sleep(100 * time.Millisecond)
			}(i)
		}

		time.Sleep(2 * time.Second)

		// Get tunnel stats
		tunnel, err := manager.GetTunnel("perf-tunnel")
		require.NoError(t, err)

		stats := tunnel.GetStats()
		fmt.Printf("Tunnel stats: %+v\n", stats)

		// Cleanup
		manager.StopTunnel("perf-tunnel")
	})
}

// BenchmarkTunnelCreation benchmarks tunnel creation
func BenchmarkTunnelCreation(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	b.Skip("Requires SSH infrastructure")

	manager := ssh.NewTunnelManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tunnelID := fmt.Sprintf("bench-tunnel-%d", i)
		cfg := ssh.TunnelConfig{
			SSHConfig: ssh.Config{
				Host: "bastion.example.com",
				Port: 22,
				User: "bench-user",
			},
			LocalAddr:  fmt.Sprintf("localhost:%d", 20000+i),
			RemoteAddr: "10.0.1.100:22",
		}

		manager.CreateTunnel(tunnelID, cfg)
	}
	b.StopTimer()

	// Cleanup
	manager.StopAll()
}
