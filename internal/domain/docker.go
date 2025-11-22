package domain

import (
	"context"
	"time"
)

// ContainerStatus represents the status of a Docker container
type ContainerStatus string

const (
	// ContainerStatusCreated indicates container is created but not started
	ContainerStatusCreated ContainerStatus = "created"
	// ContainerStatusRunning indicates container is running
	ContainerStatusRunning ContainerStatus = "running"
	// ContainerStatusPaused indicates container is paused
	ContainerStatusPaused ContainerStatus = "paused"
	// ContainerStatusRestarting indicates container is restarting
	ContainerStatusRestarting ContainerStatus = "restarting"
	// ContainerStatusExited indicates container has exited
	ContainerStatusExited ContainerStatus = "exited"
	// ContainerStatusDead indicates container is dead
	ContainerStatusDead ContainerStatus = "dead"
)

// DockerHost represents a Docker host/daemon in the infrastructure
// @Description Docker host configuration and connection details
type DockerHost struct {
	ID          string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null;uniqueIndex" validate:"required,min=3,max=255" example:"docker-host-01"`
	Description string     `json:"description" gorm:"type:text" example:"Production Docker host"`
	Endpoint    string     `json:"endpoint" gorm:"type:varchar(500);not null" validate:"required" example:"tcp://192.168.1.100:2376"`
	TLSEnabled  bool       `json:"tls_enabled" gorm:"type:boolean;default:true" example:"true"`
	CertPath    string     `json:"cert_path,omitempty" gorm:"type:varchar(500)" example:"/etc/docker/certs"`
	Version     string     `json:"version" gorm:"type:varchar(50)" example:"24.0.7"`
	ServerID    *string    `json:"server_id,omitempty" gorm:"type:uuid;index" example:"550e8400-e29b-41d4-a716-446655440000"`
	IsActive    bool       `json:"is_active" gorm:"type:boolean;default:true;index" example:"true"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime" example:"2024-01-01T00:00:00Z"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string" example:"2024-01-01T00:00:00Z"`
}

// TableName specifies the table name for DockerHost model
func (DockerHost) TableName() string {
	return "docker_hosts"
}

// Container represents a Docker container
// @Description Docker container with configuration and status
type Container struct {
	ID       string            `json:"id" example:"a1b2c3d4e5f6"`
	Name     string            `json:"name" example:"nginx-web"`
	Image    string            `json:"image" example:"nginx:latest"`
	ImageID  string            `json:"image_id" example:"sha256:abcd1234"`
	Command  string            `json:"command" example:"nginx -g 'daemon off;'"`
	Created  time.Time         `json:"created" example:"2024-01-01T00:00:00Z"`
	Status   ContainerStatus   `json:"status" example:"running"`
	State    string            `json:"state" example:"Up 2 hours"`
	Ports    []PortMapping     `json:"ports"`
	Labels   map[string]string `json:"labels" example:"app:web,env:prod"`
	Networks []string          `json:"networks" example:"bridge,custom-network"`
	Mounts   []MountPoint      `json:"mounts"`
	HostID   string            `json:"host_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// PortMapping represents a port mapping configuration
// @Description Container port mapping
type PortMapping struct {
	HostIP        string `json:"host_ip" example:"0.0.0.0"`
	HostPort      string `json:"host_port" example:"8080"`
	ContainerPort string `json:"container_port" example:"80"`
	Protocol      string `json:"protocol" example:"tcp"`
}

// MountPoint represents a volume mount
// @Description Container volume mount point
type MountPoint struct {
	Type        string `json:"type" example:"bind"`
	Source      string `json:"source" example:"/host/path"`
	Destination string `json:"destination" example:"/container/path"`
	Mode        string `json:"mode" example:"rw"`
}

// DockerImage represents a Docker image
// @Description Docker image information
type DockerImage struct {
	ID          string            `json:"id" example:"sha256:abcd1234"`
	RepoTags    []string          `json:"repo_tags" example:"nginx:latest,nginx:1.25"`
	RepoDigests []string          `json:"repo_digests"`
	Size        int64             `json:"size" example:"142000000"`
	Created     time.Time         `json:"created" example:"2024-01-01T00:00:00Z"`
	Labels      map[string]string `json:"labels"`
	HostID      string            `json:"host_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// Network represents a Docker network
// @Description Docker network configuration
type Network struct {
	ID      string            `json:"id" example:"a1b2c3d4e5f6"`
	Name    string            `json:"name" example:"custom-network"`
	Driver  string            `json:"driver" example:"bridge"`
	Scope   string            `json:"scope" example:"local"`
	Labels  map[string]string `json:"labels"`
	Options map[string]string `json:"options"`
	HostID  string            `json:"host_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// Volume represents a Docker volume
// @Description Docker volume information
type Volume struct {
	Name       string            `json:"name" example:"my-volume"`
	Driver     string            `json:"driver" example:"local"`
	Mountpoint string            `json:"mountpoint" example:"/var/lib/docker/volumes/my-volume/_data"`
	Labels     map[string]string `json:"labels"`
	Options    map[string]string `json:"options"`
	Scope      string            `json:"scope" example:"local"`
	HostID     string            `json:"host_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// ContainerStats represents container resource usage statistics
// @Description Real-time container resource statistics
type ContainerStats struct {
	ContainerID    string    `json:"container_id" example:"a1b2c3d4e5f6"`
	CPUPercent     float64   `json:"cpu_percent" example:"25.5"`
	MemoryUsage    int64     `json:"memory_usage" example:"536870912"`  // bytes
	MemoryLimit    int64     `json:"memory_limit" example:"2147483648"` // bytes
	MemoryPercent  float64   `json:"memory_percent" example:"25.0"`
	NetworkRxBytes int64     `json:"network_rx_bytes" example:"1048576"`
	NetworkTxBytes int64     `json:"network_tx_bytes" example:"524288"`
	BlockRead      int64     `json:"block_read" example:"2097152"`
	BlockWrite     int64     `json:"block_write" example:"1048576"`
	PIDs           int       `json:"pids" example:"5"`
	Timestamp      time.Time `json:"timestamp" example:"2024-01-01T00:00:00Z"`
}

// DockerHostRepository defines the interface for Docker host data persistence
type DockerHostRepository interface {
	// Create creates a new Docker host record
	Create(ctx context.Context, host *DockerHost) error

	// GetByID retrieves a Docker host by ID
	GetByID(ctx context.Context, id string) (*DockerHost, error)

	// List retrieves all Docker hosts with filtering
	List(ctx context.Context, filter DockerHostFilter) ([]*DockerHost, int64, error)

	// Update updates a Docker host
	Update(ctx context.Context, host *DockerHost) error

	// Delete soft deletes a Docker host
	Delete(ctx context.Context, id string) error
}

// DockerHostFilter represents filtering options for Docker host queries
type DockerHostFilter struct {
	ServerID *string `json:"server_id,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
	Page     int     `json:"page" validate:"min=1"`
	PageSize int     `json:"page_size" validate:"min=1,max=100"`
}

// DockerUsecase defines the business logic for Docker management
type DockerUsecase interface {
	// Docker Host Management
	CreateDockerHost(ctx context.Context, host *DockerHost) error
	GetDockerHost(ctx context.Context, id string) (*DockerHost, error)
	ListDockerHosts(ctx context.Context, filter DockerHostFilter) ([]*DockerHost, int64, error)
	UpdateDockerHost(ctx context.Context, host *DockerHost) error
	DeleteDockerHost(ctx context.Context, id string) error

	// Container Management
	ListContainers(ctx context.Context, hostID string, all bool) ([]*Container, error)
	GetContainer(ctx context.Context, hostID, containerID string) (*Container, error)
	CreateContainer(ctx context.Context, hostID string, config interface{}) (*Container, error)
	StartContainer(ctx context.Context, hostID, containerID string) error
	StopContainer(ctx context.Context, hostID, containerID string, timeout int) error
	RestartContainer(ctx context.Context, hostID, containerID string, timeout int) error
	RemoveContainer(ctx context.Context, hostID, containerID string, force bool) error
	GetContainerLogs(ctx context.Context, hostID, containerID string, tail int) (string, error)
	GetContainerStats(ctx context.Context, hostID, containerID string) (*ContainerStats, error)

	// Image Management
	ListImages(ctx context.Context, hostID string) ([]*DockerImage, error)
	PullImage(ctx context.Context, hostID, imageName string) error
	RemoveImage(ctx context.Context, hostID, imageID string, force bool) error

	// Network Management
	ListNetworks(ctx context.Context, hostID string) ([]*Network, error)
	CreateNetwork(ctx context.Context, hostID, name, driver string) (*Network, error)
	RemoveNetwork(ctx context.Context, hostID, networkID string) error

	// Volume Management
	ListVolumes(ctx context.Context, hostID string) ([]*Volume, error)
	CreateVolume(ctx context.Context, hostID, name, driver string) (*Volume, error)
	RemoveVolume(ctx context.Context, hostID, volumeName string, force bool) error
}
