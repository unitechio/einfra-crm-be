package logstream

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// LogMessage represents a single log line
type LogMessage struct {
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"` // stdout or stderr
	Message   string    `json:"message"`
}

// LogStreamer handles streaming logs from Docker containers
type LogStreamer struct {
	dockerClient *client.Client
}

// NewLogStreamer creates a new log streamer
func NewLogStreamer(cli *client.Client) *LogStreamer {
	return &LogStreamer{
		dockerClient: cli,
	}
}

// StreamLogs streams logs from a container to a channel
func (ls *LogStreamer) StreamLogs(ctx context.Context, containerID string, tail string, follow bool) (<-chan LogMessage, <-chan error) {
	logChan := make(chan LogMessage, 100)
	errChan := make(chan error, 1)

	go func() {
		defer close(logChan)
		defer close(errChan)

		options := container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     follow,
			Timestamps: true,
			Tail:       tail,
		}

		reader, err := ls.dockerClient.ContainerLogs(ctx, containerID, options)
		if err != nil {
			errChan <- fmt.Errorf("failed to get container logs: %w", err)
			return
		}
		defer reader.Close()

		// Docker logs format:
		// [8 bytes header] [payload]
		// Header: [1 byte stream type] [3 bytes ignored] [4 bytes payload size]
		// Stream type: 0: stdin, 1: stdout, 2: stderr

		header := make([]byte, 8)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Read header
				_, err := io.ReadFull(reader, header)
				if err != nil {
					if err == io.EOF {
						return
					}
					errChan <- err
					return
				}

				streamType := header[0]
				payloadSize := binary.BigEndian.Uint32(header[4:])

				// Read payload
				payload := make([]byte, payloadSize)
				_, err = io.ReadFull(reader, payload)
				if err != nil {
					errChan <- err
					return
				}

				// Parse timestamp and message
				// Docker adds timestamp at the beginning if Timestamps: true
				// Format: 2006-01-02T15:04:05.999999999Z message
				line := string(payload)
				var ts time.Time
				var msg string

				// Try to parse timestamp (RFC3339Nano)
				// Usually the first space separates timestamp and message
				if len(line) > 30 { // Basic check to avoid index out of range
					spaceIdx := -1
					for i, r := range line {
						if r == ' ' {
							spaceIdx = i
							break
						}
					}

					if spaceIdx > 0 {
						parsedTs, err := time.Parse(time.RFC3339Nano, line[:spaceIdx])
						if err == nil {
							ts = parsedTs
							msg = line[spaceIdx+1:]
						} else {
							// Fallback if parsing fails
							ts = time.Now()
							msg = line
						}
					} else {
						ts = time.Now()
						msg = line
					}
				} else {
					ts = time.Now()
					msg = line
				}

				source := "stdout"
				if streamType == 2 {
					source = "stderr"
				}

				logChan <- LogMessage{
					Timestamp: ts,
					Source:    source,
					Message:   msg,
				}
			}
		}
	}()

	return logChan, errChan
}

// MultiLogStreamer streams logs from multiple containers (e.g., for a stack)
type MultiLogStreamer struct {
	dockerClient *client.Client
}

func NewMultiLogStreamer(cli *client.Client) *MultiLogStreamer {
	return &MultiLogStreamer{dockerClient: cli}
}

func (mls *MultiLogStreamer) StreamStackLogs(ctx context.Context, containerIDs []string) (<-chan LogMessage, <-chan error) {
	mergedLogChan := make(chan LogMessage, 100*len(containerIDs))
	mergedErrChan := make(chan error, len(containerIDs))
	var wg sync.WaitGroup

	for _, cid := range containerIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			ls := NewLogStreamer(mls.dockerClient)
			logs, errs := ls.StreamLogs(ctx, id, "100", true)

			for {
				select {
				case msg, ok := <-logs:
					if !ok {
						return
					}
					// Prepend container ID or Name to message?
					// For now, let's keep it simple or maybe add a field to LogMessage if needed
					// msg.Message = fmt.Sprintf("[%s] %s", id[:12], msg.Message)
					mergedLogChan <- msg
				case err, ok := <-errs:
					if !ok {
						return
					}
					mergedErrChan <- err
					return
				case <-ctx.Done():
					return
				}
			}
		}(cid)
	}

	go func() {
		wg.Wait()
		close(mergedLogChan)
		close(mergedErrChan)
	}()

	return mergedLogChan, mergedErrChan
}
