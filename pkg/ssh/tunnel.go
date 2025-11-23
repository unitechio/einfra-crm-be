package ssh

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// Tunnel represents an SSH tunnel connection
type Tunnel struct {
	client       *Client
	localAddr    string
	remoteAddr   string
	listener     net.Listener
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	mu           sync.RWMutex
	isActive     bool
	connections  int
	lastActivity time.Time
}

// TunnelConfig represents SSH tunnel configuration
type TunnelConfig struct {
	SSHConfig  Config // SSH connection config
	LocalAddr  string // Local address to bind (e.g., "localhost:3307")
	RemoteAddr string // Remote address to forward to (e.g., "localhost:3306")
}

// NewTunnel creates a new SSH tunnel
func NewTunnel(cfg TunnelConfig) (*Tunnel, error) {
	client, err := NewClient(cfg.SSHConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH client: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Tunnel{
		client:       client,
		localAddr:    cfg.LocalAddr,
		remoteAddr:   cfg.RemoteAddr,
		ctx:          ctx,
		cancel:       cancel,
		lastActivity: time.Now(),
	}, nil
}

// Start starts the SSH tunnel
func (t *Tunnel) Start() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.isActive {
		return fmt.Errorf("tunnel is already active")
	}

	// Connect SSH client
	if err := t.client.Connect(); err != nil {
		return fmt.Errorf("failed to connect SSH client: %w", err)
	}

	// Create local listener
	listener, err := net.Listen("tcp", t.localAddr)
	if err != nil {
		t.client.Close()
		return fmt.Errorf("failed to create local listener: %w", err)
	}

	t.listener = listener
	t.isActive = true

	// Start accepting connections
	t.wg.Add(1)
	go t.acceptConnections()

	return nil
}

// Stop stops the SSH tunnel
func (t *Tunnel) Stop() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.isActive {
		return nil
	}

	t.isActive = false
	t.cancel()

	// Close listener
	if t.listener != nil {
		t.listener.Close()
	}

	// Wait for all connections to finish
	t.wg.Wait()

	// Close SSH client
	return t.client.Close()
}

// IsActive returns whether the tunnel is currently active
func (t *Tunnel) IsActive() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.isActive
}

// GetStats returns tunnel statistics
func (t *Tunnel) GetStats() map[string]interface{} {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return map[string]interface{}{
		"active":        t.isActive,
		"connections":   t.connections,
		"last_activity": t.lastActivity,
		"local_addr":    t.localAddr,
		"remote_addr":   t.remoteAddr,
	}
}

// acceptConnections accepts incoming connections and forwards them through the tunnel
func (t *Tunnel) acceptConnections() {
	defer t.wg.Done()

	for {
		select {
		case <-t.ctx.Done():
			return
		default:
			// Set accept deadline to allow checking context
			t.listener.(*net.TCPListener).SetDeadline(time.Now().Add(1 * time.Second))

			conn, err := t.listener.Accept()
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				if !t.isActive {
					return
				}
				continue
			}

			t.wg.Add(1)
			go t.handleConnection(conn)
		}
	}
}

// handleConnection handles a single connection through the tunnel
func (t *Tunnel) handleConnection(localConn net.Conn) {
	defer t.wg.Done()
	defer localConn.Close()

	// Increment connection count
	t.mu.Lock()
	t.connections++
	t.lastActivity = time.Now()
	t.mu.Unlock()

	// Dial remote address through SSH
	remoteConn, err := t.client.client.Dial("tcp", t.remoteAddr)
	if err != nil {
		return
	}
	defer remoteConn.Close()

	// Bidirectional copy
	done := make(chan struct{}, 2)

	go func() {
		io.Copy(remoteConn, localConn)
		done <- struct{}{}
	}()

	go func() {
		io.Copy(localConn, remoteConn)
		done <- struct{}{}
	}()

	// Wait for one direction to finish
	<-done

	// Update last activity
	t.mu.Lock()
	t.lastActivity = time.Now()
	t.mu.Unlock()
}

// TunnelManager manages multiple SSH tunnels
type TunnelManager struct {
	tunnels map[string]*Tunnel
	mu      sync.RWMutex
}

// NewTunnelManager creates a new tunnel manager
func NewTunnelManager() *TunnelManager {
	return &TunnelManager{
		tunnels: make(map[string]*Tunnel),
	}
}

// CreateTunnel creates and starts a new tunnel
func (tm *TunnelManager) CreateTunnel(id string, cfg TunnelConfig) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if _, exists := tm.tunnels[id]; exists {
		return fmt.Errorf("tunnel with id %s already exists", id)
	}

	tunnel, err := NewTunnel(cfg)
	if err != nil {
		return err
	}

	if err := tunnel.Start(); err != nil {
		return err
	}

	tm.tunnels[id] = tunnel
	return nil
}

// GetTunnel retrieves a tunnel by ID
func (tm *TunnelManager) GetTunnel(id string) (*Tunnel, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tunnel, exists := tm.tunnels[id]
	if !exists {
		return nil, fmt.Errorf("tunnel with id %s not found", id)
	}

	return tunnel, nil
}

// StopTunnel stops and removes a tunnel
func (tm *TunnelManager) StopTunnel(id string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tunnel, exists := tm.tunnels[id]
	if !exists {
		return fmt.Errorf("tunnel with id %s not found", id)
	}

	if err := tunnel.Stop(); err != nil {
		return err
	}

	delete(tm.tunnels, id)
	return nil
}

// ListTunnels returns all active tunnels
func (tm *TunnelManager) ListTunnels() map[string]map[string]interface{} {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	result := make(map[string]map[string]interface{})
	for id, tunnel := range tm.tunnels {
		result[id] = tunnel.GetStats()
	}

	return result
}

// StopAll stops all tunnels
func (tm *TunnelManager) StopAll() error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	var errors []error
	for id, tunnel := range tm.tunnels {
		if err := tunnel.Stop(); err != nil {
			errors = append(errors, fmt.Errorf("failed to stop tunnel %s: %w", id, err))
		}
	}

	tm.tunnels = make(map[string]*Tunnel)

	if len(errors) > 0 {
		return fmt.Errorf("errors stopping tunnels: %v", errors)
	}

	return nil
}
