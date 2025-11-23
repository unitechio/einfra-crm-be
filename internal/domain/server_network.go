package domain

import (
	"context"
	"time"
)

// NetworkInterface represents a network interface on a server
// @Description Network interface with configuration and statistics
type NetworkInterface struct {
	ID         string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	ServerID   string `json:"server_id" gorm:"type:uuid;not null;index" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name       string `json:"name" gorm:"type:varchar(100);not null" validate:"required" example:"eth0"`
	Type       string `json:"type" gorm:"type:varchar(50)" example:"ethernet"` // ethernet, wifi, loopback, etc.
	IPAddress  string `json:"ip_address" gorm:"type:varchar(45)" example:"192.168.1.100"`
	MACAddress string `json:"mac_address" gorm:"type:varchar(17)" example:"00:1B:44:11:3A:B7"`
	Netmask    string `json:"netmask" gorm:"type:varchar(45)" example:"255.255.255.0"`
	Gateway    string `json:"gateway" gorm:"type:varchar(45)" example:"192.168.1.1"`
	MTU        int    `json:"mtu" gorm:"type:int" example:"1500"`
	Speed      int    `json:"speed" gorm:"type:int" example:"1000"` // Mbps
	IsUp       bool   `json:"is_up" gorm:"type:boolean" example:"true"`

	// Statistics
	BytesReceived   int64 `json:"bytes_received" gorm:"type:bigint" example:"1073741824"`
	BytesSent       int64 `json:"bytes_sent" gorm:"type:bigint" example:"536870912"`
	PacketsReceived int64 `json:"packets_received" gorm:"type:bigint" example:"1000000"`
	PacketsSent     int64 `json:"packets_sent" gorm:"type:bigint" example:"500000"`
	ErrorsReceived  int64 `json:"errors_received" gorm:"type:bigint" example:"10"`
	ErrorsSent      int64 `json:"errors_sent" gorm:"type:bigint" example:"5"`
	DroppedReceived int64 `json:"dropped_received" gorm:"type:bigint" example:"2"`
	DroppedSent     int64 `json:"dropped_sent" gorm:"type:bigint" example:"1"`

	LastUpdatedAt time.Time  `json:"last_updated_at" gorm:"type:timestamp" example:"2024-01-01T00:00:00Z"`
	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"autoUpdateTime" example:"2024-01-01T00:00:00Z"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string"`
}

// TableName specifies the table name for NetworkInterface model
func (NetworkInterface) TableName() string {
	return "network_interfaces"
}

// NetworkStats represents network statistics for a server
// @Description Network performance statistics
type NetworkStats struct {
	ServerID         string    `json:"server_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	InterfaceName    string    `json:"interface_name" example:"eth0"`
	BandwidthInMbps  float64   `json:"bandwidth_in_mbps" example:"125.5"`
	BandwidthOutMbps float64   `json:"bandwidth_out_mbps" example:"85.3"`
	Latency          float64   `json:"latency" example:"15.5"`    // ms
	PacketLoss       float64   `json:"packet_loss" example:"0.1"` // percentage
	Timestamp        time.Time `json:"timestamp" example:"2024-01-01T00:00:00Z"`
}

// NetworkConnectivityCheck represents a network connectivity test
// @Description Network connectivity test result
type NetworkConnectivityCheck struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	ServerID     string    `json:"server_id" gorm:"type:uuid;not null;index" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	TargetHost   string    `json:"target_host" gorm:"type:varchar(255);not null" validate:"required" example:"8.8.8.8"`
	TargetPort   int       `json:"target_port,omitempty" gorm:"type:int" example:"80"`
	Protocol     string    `json:"protocol" gorm:"type:varchar(10)" example:"tcp"` // tcp, udp, icmp
	Success      bool      `json:"success" gorm:"type:boolean" example:"true"`
	Latency      float64   `json:"latency" gorm:"type:decimal(10,2)" example:"15.50"` // ms
	ErrorMessage string    `json:"error_message,omitempty" gorm:"type:text"`
	TestedAt     time.Time `json:"tested_at" gorm:"type:timestamp;not null" example:"2024-01-01T00:00:00Z"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
}

// TableName specifies the table name for NetworkConnectivityCheck model
func (NetworkConnectivityCheck) TableName() string {
	return "network_connectivity_checks"
}

// PortCheckRequest represents a request to check port connectivity
type PortCheckRequest struct {
	Host     string `json:"host" validate:"required" example:"example.com"`
	Port     int    `json:"port" validate:"required,min=1,max=65535" example:"443"`
	Protocol string `json:"protocol" validate:"required,oneof=tcp udp" example:"tcp"`
	Timeout  int    `json:"timeout" example:"5"` // seconds
}

// ServerNetworkRepository defines the interface for network data persistence
type ServerNetworkRepository interface {
	// CreateInterface creates a new network interface record
	CreateInterface(ctx context.Context, iface *NetworkInterface) error

	// GetInterfaceByID retrieves a network interface by its ID
	GetInterfaceByID(ctx context.Context, id string) (*NetworkInterface, error)

	// GetInterfacesByServerID retrieves all network interfaces for a server
	GetInterfacesByServerID(ctx context.Context, serverID string) ([]*NetworkInterface, error)

	// UpdateInterface updates an existing network interface
	UpdateInterface(ctx context.Context, iface *NetworkInterface) error

	// DeleteInterface soft deletes a network interface
	DeleteInterface(ctx context.Context, id string) error

	// CreateConnectivityCheck creates a new connectivity check record
	CreateConnectivityCheck(ctx context.Context, check *NetworkConnectivityCheck) error

	// GetConnectivityHistory retrieves connectivity check history
	GetConnectivityHistory(ctx context.Context, serverID string, limit int) ([]*NetworkConnectivityCheck, error)
}

// ServerNetworkUsecase defines the business logic for network management
type ServerNetworkUsecase interface {
	// GetNetworkInterfaces retrieves all network interfaces for a server
	GetNetworkInterfaces(ctx context.Context, serverID string) ([]*NetworkInterface, error)

	// RefreshNetworkInterfaces refreshes network interface data from the server
	RefreshNetworkInterfaces(ctx context.Context, serverID string) error

	// GetNetworkStats retrieves current network statistics
	GetNetworkStats(ctx context.Context, serverID string) ([]*NetworkStats, error)

	// CheckConnectivity checks network connectivity to a target
	CheckConnectivity(ctx context.Context, serverID, targetHost string, targetPort int, protocol string) (*NetworkConnectivityCheck, error)

	// TestPort tests connectivity to a specific port
	TestPort(ctx context.Context, serverID string, request PortCheckRequest) (bool, error)

	// GetConnectivityHistory retrieves connectivity check history
	GetConnectivityHistory(ctx context.Context, serverID string, limit int) ([]*NetworkConnectivityCheck, error)

	// MonitorBandwidth monitors bandwidth usage
	MonitorBandwidth(ctx context.Context, serverID string, duration int) ([]*NetworkStats, error)
}
