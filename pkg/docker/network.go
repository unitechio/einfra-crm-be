package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
)

// NetworkConnect connects a container to a network
func (c *Client) NetworkConnect(ctx context.Context, networkID, containerID string) error {
	return c.cli.NetworkConnect(ctx, networkID, containerID, nil)
}

// NetworkDisconnect disconnects a container from a network
func (c *Client) NetworkDisconnect(ctx context.Context, networkID, containerID string) error {
	return c.cli.NetworkDisconnect(ctx, networkID, containerID, false)
}

// NetworkCreate creates a new network
func (c *Client) NetworkCreate(ctx context.Context, name, driver string) (string, error) {
	options := types.NetworkCreate{
		Driver: driver,
	}

	resp, err := c.cli.NetworkCreate(ctx, name, options)
	if err != nil {
		return "", fmt.Errorf("failed to create network: %w", err)
	}

	return resp.ID, nil
}

// NetworkRemove removes a network
func (c *Client) NetworkRemove(ctx context.Context, networkID string) error {
	return c.cli.NetworkRemove(ctx, networkID)
}

// NetworkInspect inspects a network
func (c *Client) NetworkInspect(ctx context.Context, networkID string) (map[string]interface{}, error) {
	resource, err := c.cli.NetworkInspect(ctx, networkID, types.NetworkInspectOptions{})
	if err != nil {
		return nil, err
	}

	// Convert to map for generic response
	result := map[string]interface{}{
		"id":         resource.ID,
		"name":       resource.Name,
		"driver":     resource.Driver,
		"scope":      resource.Scope,
		"internal":   resource.Internal,
		"attachable": resource.Attachable,
		"ingress":    resource.Ingress,
		"ipam":       resource.IPAM,
		"containers": resource.Containers,
		"options":    resource.Options,
		"labels":     resource.Labels,
	}

	return result, nil
}
