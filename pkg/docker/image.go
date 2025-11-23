package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
)

// BuildConfig represents image build configuration
type BuildConfig struct {
	Dockerfile string
	Context    io.Reader
	Tags       []string
	BuildArgs  map[string]string
	Labels     map[string]string
	NoCache    bool
}

// AuthConfig represents registry authentication configuration
type AuthConfig struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Auth          string `json:"auth,omitempty"`
	Email         string `json:"email,omitempty"`
	ServerAddress string `json:"serveraddress,omitempty"`
	IdentityToken string `json:"identitytoken,omitempty"`
	RegistryToken string `json:"registrytoken,omitempty"`
}

// ImageBuild builds an image from a Dockerfile
func (c *Client) ImageBuild(ctx context.Context, config BuildConfig) (io.ReadCloser, error) {
	// Convert build args to map[string]*string
	buildArgs := make(map[string]*string)
	for k, v := range config.BuildArgs {
		val := v
		buildArgs[k] = &val
	}

	options := types.ImageBuildOptions{
		Dockerfile: config.Dockerfile,
		Tags:       config.Tags,
		BuildArgs:  buildArgs,
		Labels:     config.Labels,
		NoCache:    config.NoCache,
		Remove:     true,
	}

	resp, err := c.cli.ImageBuild(ctx, config.Context, options)
	if err != nil {
		return nil, fmt.Errorf("failed to build image: %w", err)
	}

	return resp.Body, nil
}

// ImagePush pushes an image to a registry
func (c *Client) ImagePush(ctx context.Context, imageName string, authConfig AuthConfig) (io.ReadCloser, error) {
	options := image.PushOptions{}

	if authConfig.Username != "" && authConfig.Password != "" {
		encodedAuth, err := encodeAuthToBase64(authConfig)
		if err != nil {
			return nil, err
		}
		options.RegistryAuth = encodedAuth
	}

	resp, err := c.cli.ImagePush(ctx, imageName, options)
	if err != nil {
		return nil, fmt.Errorf("failed to push image: %w", err)
	}

	return resp, nil
}

// ImagePull pulls an image from a registry
func (c *Client) ImagePull(ctx context.Context, imageName string, authConfig AuthConfig) (io.ReadCloser, error) {
	options := image.PullOptions{}

	if authConfig.Username != "" && authConfig.Password != "" {
		encodedAuth, err := encodeAuthToBase64(authConfig)
		if err != nil {
			return nil, err
		}
		options.RegistryAuth = encodedAuth
	}

	resp, err := c.cli.ImagePull(ctx, imageName, options)
	if err != nil {
		return nil, fmt.Errorf("failed to pull image: %w", err)
	}

	return resp, nil
}

// ImageRemove removes an image
func (c *Client) ImageRemove(ctx context.Context, imageID string, force bool) ([]string, error) {
	options := image.RemoveOptions{
		Force:         force,
		PruneChildren: true,
	}

	deleted, err := c.cli.ImageRemove(ctx, imageID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to remove image: %w", err)
	}

	result := make([]string, len(deleted))
	for i, d := range deleted {
		if d.Deleted != "" {
			result[i] = d.Deleted
		} else {
			result[i] = d.Untagged
		}
	}

	return result, nil
}

// ImageInspect inspects an image
func (c *Client) ImageInspect(ctx context.Context, imageID string) (map[string]interface{}, error) {
	inspect, _, err := c.cli.ImageInspectWithRaw(ctx, imageID)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect image: %w", err)
	}

	result := map[string]interface{}{
		"id":             inspect.ID,
		"repo_tags":      inspect.RepoTags,
		"repo_digests":   inspect.RepoDigests,
		"parent":         inspect.Parent,
		"comment":        inspect.Comment,
		"created":        inspect.Created,
		"container":      inspect.Container,
		"docker_version": inspect.DockerVersion,
		"author":         inspect.Author,
		"config":         inspect.Config,
		"architecture":   inspect.Architecture,
		"os":             inspect.Os,
		"size":           inspect.Size,
		"virtual_size":   inspect.VirtualSize,
	}

	return result, nil
}

func encodeAuthToBase64(authConfig AuthConfig) (string, error) {
	authJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", fmt.Errorf("failed to marshal auth config: %w", err)
	}
	return base64.URLEncoding.EncodeToString(authJSON), nil
}
