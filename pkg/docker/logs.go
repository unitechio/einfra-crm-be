package docker

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/unitechio/einfra-be/pkg/logstream"
)

// ContainerLogsStream streams logs from a container
func (c *Client) ContainerLogsStream(ctx context.Context, containerID string, tail string, follow bool) (<-chan logstream.LogMessage, <-chan error, error) {
	logChan := make(chan logstream.LogMessage, 100)
	errChan := make(chan error, 1)

	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     follow,
		Timestamps: true,
		Tail:       tail,
	}

	reader, err := c.cli.ContainerLogs(ctx, containerID, options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get container logs: %w", err)
	}

	go func() {
		defer close(logChan)
		defer close(errChan)
		defer reader.Close()

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
				line := string(payload)
				var ts time.Time
				var msg string

				if len(line) > 30 {
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

				logChan <- logstream.LogMessage{
					Timestamp: ts,
					Source:    source,
					Message:   msg,
				}
			}
		}
	}()

	return logChan, errChan, nil
}
