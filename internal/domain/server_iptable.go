package domain

import (
	"context"
	"time"
)

// IPTableChain represents an iptables chain
type IPTableChain string

const (
	// IPTableChainInput represents INPUT chain
	IPTableChainInput IPTableChain = "INPUT"
	// IPTableChainOutput represents OUTPUT chain
	IPTableChainOutput IPTableChain = "OUTPUT"
	// IPTableChainForward represents FORWARD chain
	IPTableChainForward IPTableChain = "FORWARD"
	// IPTableChainPrerouting represents PREROUTING chain
	IPTableChainPrerouting IPTableChain = "PREROUTING"
	// IPTableChainPostrouting represents POSTROUTING chain
	IPTableChainPostrouting IPTableChain = "POSTROUTING"
)

// IPTableAction represents an iptables action/target
type IPTableAction string

const (
	// IPTableActionAccept accepts the packet
	IPTableActionAccept IPTableAction = "ACCEPT"
	// IPTableActionDrop drops the packet
	IPTableActionDrop IPTableAction = "DROP"
	// IPTableActionReject rejects the packet
	IPTableActionReject IPTableAction = "REJECT"
	// IPTableActionLog logs the packet
	IPTableActionLog IPTableAction = "LOG"
	// IPTableActionMasquerade masquerades the packet (NAT)
	IPTableActionMasquerade IPTableAction = "MASQUERADE"
)

// IPTableProtocol represents network protocol
type IPTableProtocol string

const (
	// IPTableProtocolTCP represents TCP protocol
	IPTableProtocolTCP IPTableProtocol = "tcp"
	// IPTableProtocolUDP represents UDP protocol
	IPTableProtocolUDP IPTableProtocol = "udp"
	// IPTableProtocolICMP represents ICMP protocol
	IPTableProtocolICMP IPTableProtocol = "icmp"
	// IPTableProtocolAll represents all protocols
	IPTableProtocolAll IPTableProtocol = "all"
)

// ServerIPTable represents an iptables rule on a server
// @Description Server iptables/firewall rule configuration
type ServerIPTable struct {
	ID          string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	ServerID    string `json:"server_id" gorm:"type:uuid;not null;index" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string `json:"name" gorm:"type:varchar(255);not null" validate:"required" example:"allow-http"`
	Description string `json:"description" gorm:"type:text" example:"Allow HTTP traffic"`
	Enabled     bool   `json:"enabled" gorm:"type:boolean;default:true" example:"true"`

	// Rule configuration
	Chain      IPTableChain    `json:"chain" gorm:"type:varchar(50);not null" validate:"required" example:"INPUT"`
	Action     IPTableAction   `json:"action" gorm:"type:varchar(50);not null" validate:"required" example:"ACCEPT"`
	Protocol   IPTableProtocol `json:"protocol" gorm:"type:varchar(10)" example:"tcp"`
	SourceIP   string          `json:"source_ip,omitempty" gorm:"type:varchar(45)" example:"0.0.0.0/0"`
	SourcePort string          `json:"source_port,omitempty" gorm:"type:varchar(20)" example:"1024:65535"`
	DestIP     string          `json:"dest_ip,omitempty" gorm:"type:varchar(45)" example:"192.168.1.100"`
	DestPort   string          `json:"dest_port,omitempty" gorm:"type:varchar(20)" example:"80"`
	Interface  string          `json:"interface,omitempty" gorm:"type:varchar(50)" example:"eth0"`
	State      string          `json:"state,omitempty" gorm:"type:varchar(100)" example:"NEW,ESTABLISHED"`

	// Rule metadata
	Position int    `json:"position" gorm:"type:int" example:"1"` // Rule position in chain
	RawRule  string `json:"raw_rule,omitempty" gorm:"type:text"`  // Full iptables command
	Comment  string `json:"comment,omitempty" gorm:"type:varchar(255)" example:"Allow web traffic"`

	// Tracking
	PacketCount int64     `json:"packet_count" gorm:"type:bigint;default:0" example:"1000"`
	ByteCount   int64     `json:"byte_count" gorm:"type:bigint;default:0" example:"1048576"`
	LastApplied time.Time `json:"last_applied" gorm:"type:timestamp" example:"2024-01-01T00:00:00Z"`

	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime" example:"2024-01-01T00:00:00Z"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string"`
}

// TableName specifies the table name for ServerIPTable model
func (ServerIPTable) TableName() string {
	return "server_iptables"
}

// IPTableBackup represents a backup of iptables configuration
// @Description Backup of server iptables configuration
type IPTableBackup struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	ServerID    string    `json:"server_id" gorm:"type:uuid;not null;index" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null" validate:"required" example:"backup-2024-01-01"`
	Description string    `json:"description" gorm:"type:text" example:"Pre-update backup"`
	Content     string    `json:"content" gorm:"type:text;not null"` // Full iptables-save output
	RuleCount   int       `json:"rule_count" gorm:"type:int" example:"25"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
}

// TableName specifies the table name for IPTableBackup model
func (IPTableBackup) TableName() string {
	return "iptable_backups"
}

// IPTableFilter represents filtering options for iptables queries
type IPTableFilter struct {
	ServerID string          `json:"server_id,omitempty"`
	Chain    IPTableChain    `json:"chain,omitempty"`
	Action   IPTableAction   `json:"action,omitempty"`
	Protocol IPTableProtocol `json:"protocol,omitempty"`
	Enabled  *bool           `json:"enabled,omitempty"`
	Page     int             `json:"page" validate:"min=1"`
	PageSize int             `json:"page_size" validate:"min=1,max=100"`
}

// ServerIPTableRepository defines the interface for iptables data persistence
type ServerIPTableRepository interface {
	// Create creates a new iptables rule record
	Create(ctx context.Context, rule *ServerIPTable) error

	// GetByID retrieves an iptables rule by its ID
	GetByID(ctx context.Context, id string) (*ServerIPTable, error)

	// List retrieves all iptables rules with pagination and filtering
	List(ctx context.Context, filter IPTableFilter) ([]*ServerIPTable, int64, error)

	// Update updates an existing iptables rule
	Update(ctx context.Context, rule *ServerIPTable) error

	// Delete soft deletes an iptables rule
	Delete(ctx context.Context, id string) error

	// GetByServerID retrieves all iptables rules for a server
	GetByServerID(ctx context.Context, serverID string) ([]*ServerIPTable, error)

	// CreateBackup creates a new iptables backup
	CreateBackup(ctx context.Context, backup *IPTableBackup) error

	// GetBackups retrieves iptables backups for a server
	GetBackups(ctx context.Context, serverID string, limit int) ([]*IPTableBackup, error)

	// GetBackupByID retrieves a specific backup
	GetBackupByID(ctx context.Context, id string) (*IPTableBackup, error)
}

// ServerIPTableUsecase defines the business logic for iptables management
type ServerIPTableUsecase interface {
	// ListRules retrieves all iptables rules for a server
	ListRules(ctx context.Context, serverID string) ([]*ServerIPTable, error)

	// GetRule retrieves an iptables rule by ID
	GetRule(ctx context.Context, id string) (*ServerIPTable, error)

	// AddRule adds a new iptables rule
	AddRule(ctx context.Context, rule *ServerIPTable) error

	// UpdateRule updates an existing iptables rule
	UpdateRule(ctx context.Context, rule *ServerIPTable) error

	// DeleteRule deletes an iptables rule
	DeleteRule(ctx context.Context, id string) error

	// ApplyRules applies all enabled rules to the server
	ApplyRules(ctx context.Context, serverID string) error

	// RefreshRules refreshes rules from the server
	RefreshRules(ctx context.Context, serverID string) error

	// BackupConfiguration creates a backup of current iptables configuration
	BackupConfiguration(ctx context.Context, serverID, name, description string) (*IPTableBackup, error)

	// RestoreConfiguration restores iptables from a backup
	RestoreConfiguration(ctx context.Context, backupID string) error

	// GetBackups retrieves backup history
	GetBackups(ctx context.Context, serverID string, limit int) ([]*IPTableBackup, error)

	// FlushRules removes all rules from a chain
	FlushRules(ctx context.Context, serverID string, chain IPTableChain) error
}
