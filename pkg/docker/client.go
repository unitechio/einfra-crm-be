package docker

import (
	"github.com/docker/docker/client"
)

type Client struct {
	cli *client.Client
}

func NewClient(host string) (*Client, error) {
	cli, err := client.NewClientWithOpts(
		client.WithHost(host),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		cli: cli,
	}, nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}
