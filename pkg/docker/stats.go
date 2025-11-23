package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
)

// ContainerStats represents container resource usage statistics
type ContainerStats struct {
	Read      time.Time `json:"read"`
	Preread   time.Time `json:"preread"`
	PidsStats struct {
		Current uint64 `json:"current"`
	} `json:"pids_stats"`
	BlkioStats struct {
		IoServiceBytesRecursive []interface{} `json:"io_service_bytes_recursive"`
	} `json:"blkio_stats"`
	NumProcs     uint32   `json:"num_procs"`
	StorageStats struct{} `json:"storage_stats"`
	CPUStats     struct {
		CPUUsage struct {
			TotalUsage        uint64   `json:"total_usage"`
			PercpuUsage       []uint64 `json:"percpu_usage"`
			UsageInKernelmode uint64   `json:"usage_in_kernelmode"`
			UsageInUsermode   uint64   `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		SystemCPUUsage uint64 `json:"system_cpu_usage"`
		OnlineCPUs     uint32 `json:"online_cpus"`
		ThrottlingData struct {
			Periods          uint64 `json:"periods"`
			ThrottledPeriods uint64 `json:"throttled_periods"`
			ThrottledTime    uint64 `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"cpu_stats"`
	PrecpuStats struct {
		CPUUsage struct {
			TotalUsage        uint64   `json:"total_usage"`
			PercpuUsage       []uint64 `json:"percpu_usage"`
			UsageInKernelmode uint64   `json:"usage_in_kernelmode"`
			UsageInUsermode   uint64   `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		SystemCPUUsage uint64 `json:"system_cpu_usage"`
		OnlineCPUs     uint32 `json:"online_cpus"`
		ThrottlingData struct {
			Periods          uint64 `json:"periods"`
			ThrottledPeriods uint64 `json:"throttled_periods"`
			ThrottledTime    uint64 `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"precpu_stats"`
	MemoryStats struct {
		Usage    uint64 `json:"usage"`
		MaxUsage uint64 `json:"max_usage"`
		Stats    struct {
			ActiveAnon              uint64 `json:"active_anon"`
			ActiveFile              uint64 `json:"active_file"`
			Cache                   uint64 `json:"cache"`
			Dirty                   uint64 `json:"dirty"`
			HierarchicalMemoryLimit uint64 `json:"hierarchical_memory_limit"`
			HierarchicalMemswLimit  uint64 `json:"hierarchical_memsw_limit"`
			InactiveAnon            uint64 `json:"inactive_anon"`
			InactiveFile            uint64 `json:"inactive_file"`
			MappedFile              uint64 `json:"mapped_file"`
			Pgfault                 uint64 `json:"pgfault"`
			Pgmajfault              uint64 `json:"pgmajfault"`
			Pgpgin                  uint64 `json:"pgpgin"`
			Pgpgout                 uint64 `json:"pgpgout"`
			Rss                     uint64 `json:"rss"`
			RssHuge                 uint64 `json:"rss_huge"`
			TotalActiveAnon         uint64 `json:"total_active_anon"`
			TotalActiveFile         uint64 `json:"total_active_file"`
			TotalCache              uint64 `json:"total_cache"`
			TotalDirty              uint64 `json:"total_dirty"`
			TotalInactiveAnon       uint64 `json:"total_inactive_anon"`
			TotalInactiveFile       uint64 `json:"total_inactive_file"`
			TotalMappedFile         uint64 `json:"total_mapped_file"`
			TotalPgfault            uint64 `json:"total_pgfault"`
			TotalPgmajfault         uint64 `json:"total_pgmajfault"`
			TotalPgpgin             uint64 `json:"total_pgpgin"`
			TotalPgpgout            uint64 `json:"total_pgpgout"`
			TotalRss                uint64 `json:"total_rss"`
			TotalRssHuge            uint64 `json:"total_rss_huge"`
			TotalUnevictable        uint64 `json:"total_unevictable"`
			TotalWriteback          uint64 `json:"total_writeback"`
			Unevictable             uint64 `json:"unevictable"`
			Writeback               uint64 `json:"writeback"`
		} `json:"stats"`
		Limit uint64 `json:"limit"`
	} `json:"memory_stats"`
	Name     string `json:"name"`
	ID       string `json:"id"`
	Networks map[string]struct {
		RxBytes   uint64 `json:"rx_bytes"`
		RxPackets uint64 `json:"rx_packets"`
		RxErrors  uint64 `json:"rx_errors"`
		RxDropped uint64 `json:"rx_dropped"`
		TxBytes   uint64 `json:"tx_bytes"`
		TxPackets uint64 `json:"tx_packets"`
		TxErrors  uint64 `json:"tx_errors"`
		TxDropped uint64 `json:"tx_dropped"`
	} `json:"networks"`
}

// ContainerCommitConfig represents configuration for committing a container
type ContainerCommitConfig struct {
	ContainerID string
	Repository  string
	Tag         string
	Comment     string
	Author      string
	Pause       bool
	Changes     []string
	Config      *container.Config
}

// ContainerStatsStream streams container stats
func (c *Client) ContainerStatsStream(ctx context.Context, containerID string) (<-chan *ContainerStats, <-chan error, error) {
	statsChan := make(chan *ContainerStats)
	errChan := make(chan error)

	stats, err := c.cli.ContainerStats(ctx, containerID, true)
	if err != nil {
		return nil, nil, err
	}

	go func() {
		defer close(statsChan)
		defer close(errChan)
		defer stats.Body.Close()

		decoder := json.NewDecoder(stats.Body)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				var s ContainerStats
				if err := decoder.Decode(&s); err != nil {
					if err == io.EOF {
						return
					}
					errChan <- err
					return
				}
				statsChan <- &s
			}
		}
	}()

	return statsChan, errChan, nil
}

// ContainerStatsOnce gets container stats once
func (c *Client) ContainerStatsOnce(ctx context.Context, containerID string) (*ContainerStats, error) {
	stats, err := c.cli.ContainerStats(ctx, containerID, false)
	if err != nil {
		return nil, err
	}
	defer stats.Body.Close()

	var s ContainerStats
	if err := json.NewDecoder(stats.Body).Decode(&s); err != nil {
		return nil, fmt.Errorf("failed to decode stats: %w", err)
	}

	return &s, nil
}

// ContainerPause pauses a container
func (c *Client) ContainerPause(ctx context.Context, containerID string) error {
	return c.cli.ContainerPause(ctx, containerID)
}

// ContainerUnpause unpauses a container
func (c *Client) ContainerUnpause(ctx context.Context, containerID string) error {
	return c.cli.ContainerUnpause(ctx, containerID)
}

// ContainerCommit commits a container to a new image
func (c *Client) ContainerCommit(ctx context.Context, config ContainerCommitConfig) (string, error) {
	options := container.CommitOptions{
		Reference: fmt.Sprintf("%s:%s", config.Repository, config.Tag),
		Comment:   config.Comment,
		Author:    config.Author,
		Changes:   config.Changes,
		Pause:     config.Pause,
		Config:    config.Config,
	}

	resp, err := c.cli.ContainerCommit(ctx, config.ContainerID, options)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}
