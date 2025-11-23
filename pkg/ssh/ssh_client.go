package ssh

import (
	"context"
	"fmt"
	"io"
	"time"

	"golang.org/x/crypto/ssh"
)

// Client represents an SSH client connection
type Client struct {
	config *ssh.ClientConfig
	host   string
	port   int
	client *ssh.Client
}

// Config represents SSH connection configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	KeyPath  string
	Timeout  time.Duration
}

// CommandResult represents the result of a command execution
type CommandResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Duration time.Duration
}

// NewClient creates a new SSH client
func NewClient(cfg Config) (*Client, error) {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}

	authMethods := []ssh.AuthMethod{}

	// Add password authentication if provided
	if cfg.Password != "" {
		authMethods = append(authMethods, ssh.Password(cfg.Password))
	}

	// Add key-based authentication if provided
	if cfg.KeyPath != "" {
		key, err := loadPrivateKey(cfg.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load private key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(key))
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("no authentication method provided")
	}

	config := &ssh.ClientConfig{
		User:            cfg.User,
		Auth:            authMethods,
		Timeout:         cfg.Timeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: Implement proper host key verification
	}

	return &Client{
		config: config,
		host:   cfg.Host,
		port:   cfg.Port,
	}, nil
}

// Connect establishes an SSH connection
func (c *Client) Connect() error {
	addr := fmt.Sprintf("%s:%d", c.host, c.port)
	client, err := ssh.Dial("tcp", addr, c.config)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	c.client = client
	return nil
}

// Close closes the SSH connection
func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// ExecuteCommand executes a command on the remote server
func (c *Client) ExecuteCommand(ctx context.Context, command string) (*CommandResult, error) {
	if c.client == nil {
		if err := c.Connect(); err != nil {
			return nil, err
		}
	}

	session, err := c.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Capture stdout and stderr
	stdoutPipe, err := session.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderrPipe, err := session.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	startTime := time.Now()

	// Start the command
	if err := session.Start(command); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	// Read output
	stdoutBytes, err := io.ReadAll(stdoutPipe)
	if err != nil {
		return nil, fmt.Errorf("failed to read stdout: %w", err)
	}

	stderrBytes, err := io.ReadAll(stderrPipe)
	if err != nil {
		return nil, fmt.Errorf("failed to read stderr: %w", err)
	}

	// Wait for command to complete
	err = session.Wait()
	duration := time.Since(startTime)

	result := &CommandResult{
		Stdout:   string(stdoutBytes),
		Stderr:   string(stderrBytes),
		Duration: duration,
	}

	// Get exit code
	if err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			result.ExitCode = exitErr.ExitStatus()
		} else {
			return result, fmt.Errorf("command execution failed: %w", err)
		}
	} else {
		result.ExitCode = 0
	}

	return result, nil
}

// ExecuteCommands executes multiple commands sequentially
func (c *Client) ExecuteCommands(ctx context.Context, commands []string) ([]*CommandResult, error) {
	results := make([]*CommandResult, 0, len(commands))

	for _, cmd := range commands {
		result, err := c.ExecuteCommand(ctx, cmd)
		if err != nil {
			return results, err
		}
		results = append(results, result)

		// Stop if command failed
		if result.ExitCode != 0 {
			return results, fmt.Errorf("command failed with exit code %d: %s", result.ExitCode, result.Stderr)
		}
	}

	return results, nil
}

// FileExists checks if a file exists on the remote server
func (c *Client) FileExists(ctx context.Context, path string) (bool, error) {
	command := fmt.Sprintf("test -f %s && echo 'exists' || echo 'not found'", path)
	result, err := c.ExecuteCommand(ctx, command)
	if err != nil {
		return false, err
	}
	return result.Stdout == "exists\n", nil
}

// loadPrivateKey loads a private key from file
func loadPrivateKey(path string) (ssh.Signer, error) {
	// This is a placeholder - in production, implement proper key loading
	// with support for encrypted keys
	return nil, fmt.Errorf("key loading not implemented yet")
}
