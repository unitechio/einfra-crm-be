package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
)

// Events returns a channel of Docker events
func (c *Client) Events(ctx context.Context) (<-chan events.Message, <-chan error) {
	options := types.EventsOptions{}
	return c.cli.Events(ctx, options)
}
