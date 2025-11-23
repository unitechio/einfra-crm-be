package domain

import (
	"context"
	"time"
)

// ServerStatus represents the operational status of a server
type ServerStatus string

const (
	// ServerStatusOnline indicates the server is running and reachable
	ServerStatusOnline ServerStatus = "online"
	// ServerStatusOffline indicates the server is not reachable
	ServerStatusOffline ServerStatus = "offline"
	// ServerStatusMaintenance indicates the server is under maintenance
	ServerStatusMaintenance ServerStatus = "maintenance"
	// ServerStatusError indicates the server has errors
	ServerStatusError ServerStatus = "error"
)

// ServerOS represents the operating system type
type ServerOS string

const (
	// Linux distributions
	ServerOSUbuntu     ServerOS = "ubuntu"
	ServerOSDebian     ServerOS = "debian"
	ServerOSCentOS     ServerOS = "centos"
	ServerOSRockyLinux ServerOS = "rocky"
	ServerOSAlmaLinux  ServerOS = "alma"
	ServerOSFedora     ServerOS = "fedora"
	ServerOSRHEL       ServerOS = "rhel"

	// Windows
	ServerOSWindowsServer2016 ServerOS = "windows-server-2016"
	ServerOSWindowsServer2019 ServerOS = "windows-server-2019"
	ServerOSWindowsServer2022 ServerOS = "windows-server-2022"

	// macOS
	ServerOSMacOS ServerOS = "macos"

	// Generic fallbacks
	ServerOSLinux   ServerOS = "linux"
	ServerOSWindows ServerOS = "windows"
)

// Server represents a physical or virtual server in the infrastructure
// @Description Server entity with hardware specifications and status
type Server struct {
	ID          string       `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string       `json:"name" gorm:"type:varchar(255);not null;uniqueIndex" validate:"required,min=3,max=255" example:"web-server-01"`
	Description string       `json:"description" gorm:"type:text" example:"Production web server"`
	IPAddress   string       `json:"ip_address" gorm:"type:varchar(45);not null;uniqueIndex" validate:"required,ip" example:"192.168.1.100"`
	Hostname    string       `json:"hostname" gorm:"type:varchar(255)" example:"web01.example.com"`
	OS          ServerOS     `json:"os" gorm:"type:varchar(50);not null" validate:"required" example:"ubuntu"`
	OSVersion   string       `json:"os_version" gorm:"type:varchar(100)" example:"Ubuntu 22.04 LTS"`
	CPUCores    int          `json:"cpu_cores" gorm:"type:int;not null" validate:"required,min=1" example:"8"`
	CPUModel    string       `json:"cpu_model" gorm:"type:varchar(255)" example:"Intel Xeon E5-2680 v4"`
	MemoryGB    float64      `json:"memory_gb" gorm:"type:decimal(10,2);not null" validate:"required,min=0.1" example:"32.00"`
	DiskGB      float64      `json:"disk_gb" gorm:"type:decimal(10,2);not null" validate:"required,min=1" example:"500.00"`
	Status      ServerStatus `json:"status" gorm:"type:varchar(50);not null;index" validate:"required,oneof=online offline maintenance error" example:"online"`
	Tags        []string     `json:"tags" gorm:"type:jsonb" example:"production,web,nginx"`
	Location    string       `json:"location" gorm:"type:varchar(255)" example:"DC-US-EAST-1"`
	Provider    string       `json:"provider" gorm:"type:varchar(100)" example:"AWS"`

	// SSH Connection Configuration
	SSHPort     int    `json:"ssh_port" gorm:"type:int;default:22" example:"22"`
	SSHUser     string `json:"ssh_user" gorm:"type:varchar(100)" example:"root"`
	SSHPassword string `json:"ssh_password,omitempty" gorm:"type:varchar(255)" swaggerignore:"true"` // Encrypted
	SSHKeyPath  string `json:"ssh_key_path,omitempty" gorm:"type:varchar(500)" example:"/path/to/key.pem"`

	// SSH Tunnel Configuration (for private servers)
	TunnelEnabled bool   `json:"tunnel_enabled" gorm:"type:boolean;default:false" example:"false"`
	TunnelHost    string `json:"tunnel_host,omitempty" gorm:"type:varchar(255)" example:"bastion.example.com"`
	TunnelPort    int    `json:"tunnel_port,omitempty" gorm:"type:int;default:22" example:"22"`
	TunnelUser    string `json:"tunnel_user,omitempty" gorm:"type:varchar(100)" example:"tunnel-user"`
	TunnelKeyPath string `json:"tunnel_key_path,omitempty" gorm:"type:varchar(500)" example:"/path/to/tunnel-key.pem"`

	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime" example:"2024-01-01T00:00:00Z"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string" example:"2024-01-01T00:00:00Z"`
}

// TableName specifies the table name for Server model
func (Server) TableName() string {
	return "servers"
}

// ServerMetrics represents real-time metrics of a server
// @Description Real-time server performance metrics
type ServerMetrics struct {
	ServerID       string    `json:"server_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	CPUUsage       float64   `json:"cpu_usage" example:"45.5"`           // Percentage
	MemoryUsage    float64   `json:"memory_usage" example:"70.2"`        // Percentage
	DiskUsage      float64   `json:"disk_usage" example:"60.8"`          // Percentage
	NetworkInMbps  float64   `json:"network_in_mbps" example:"125.5"`    // Mbps
	NetworkOutMbps float64   `json:"network_out_mbps" example:"85.3"`    // Mbps
	Uptime         int64     `json:"uptime" example:"864000"`            // Seconds
	LoadAverage    []float64 `json:"load_average" example:"1.5,1.2,1.0"` // 1, 5, 15 minutes
	Timestamp      time.Time `json:"timestamp" example:"2024-01-01T00:00:00Z"`
}

// ServerRepository defines the interface for server data persistence
type ServerRepository interface {
	Create(ctx context.Context, server *Server) error
	GetByID(ctx context.Context, id string) (*Server, error)
	GetByIPAddress(ctx context.Context, ip string) (*Server, error)
	List(ctx context.Context, filter ServerFilter) ([]*Server, int64, error)
	Update(ctx context.Context, server *Server) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status ServerStatus) error
}

type ServerFilter struct {
	Status   ServerStatus `json:"status,omitempty"`
	OS       ServerOS     `json:"os,omitempty"`
	Location string       `json:"location,omitempty"`
	Provider string       `json:"provider,omitempty"`
	Tags     []string     `json:"tags,omitempty"`
	Page     int          `json:"page" validate:"min=1"`
	PageSize int          `json:"page_size" validate:"min=1,max=100"`
}

type ServerUsecase interface {
	CreateServer(ctx context.Context, server *Server) error
	GetServer(ctx context.Context, id string) (*Server, error)
	ListServers(ctx context.Context, filter ServerFilter) ([]*Server, int64, error)
	UpdateServer(ctx context.Context, server *Server) error
	DeleteServer(ctx context.Context, id string) error
	GetServerMetrics(ctx context.Context, serverID string) (*ServerMetrics, error)
	HealthCheck(ctx context.Context, serverID string) (bool, error)
}
