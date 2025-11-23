package domain

import "time"

// TunnelConfig represents SSH tunnel configuration for infrastructure connections
type TunnelConfig struct {
	Enabled    bool   `json:"enabled" gorm:"type:boolean;default:false"`
	SSHHost    string `json:"ssh_host" gorm:"type:varchar(255)"` // Jump/bastion host
	SSHPort    int    `json:"ssh_port" gorm:"type:int;default:22"`
	SSHUser    string `json:"ssh_user" gorm:"type:varchar(100)"`
	SSHKeyPath string `json:"ssh_key_path" gorm:"type:varchar(500)"`
	LocalPort  int    `json:"local_port" gorm:"type:int"` // Local port for tunnel
}

// KubeConfig represents Kubernetes configuration
type KubeConfig struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	ClusterID   string    `json:"cluster_id" gorm:"type:uuid;index"`            // Reference to K8s cluster
	ConfigType  string    `json:"config_type" gorm:"type:varchar(50);not null"` // "file", "inline", "credentials"
	ConfigData  string    `json:"config_data" gorm:"type:text"`                 // Base64 encoded kubeconfig file or inline YAML
	ContextName string    `json:"context_name" gorm:"type:varchar(255)"`        // Default context to use
	Description string    `json:"description" gorm:"type:text"`
	IsDefault   bool      `json:"is_default" gorm:"type:boolean;default:false"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for KubeConfig model
func (KubeConfig) TableName() string {
	return "kube_configs"
}
