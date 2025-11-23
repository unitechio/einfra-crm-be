package domain

import "time"

// DockerStack represents a Docker Compose stack
type DockerStack struct {
	ID          string            `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string            `json:"name" gorm:"type:varchar(255);not null;uniqueIndex"`
	ComposeFile string            `json:"compose_file" gorm:"type:text;not null"` // YAML content
	EnvVars     map[string]string `json:"env_vars" gorm:"type:jsonb"`
	Status      StackStatus       `json:"status" gorm:"type:varchar(50);not null"`
	DockerHost  string            `json:"docker_host" gorm:"type:varchar(255)"`
	ProjectName string            `json:"project_name" gorm:"type:varchar(255)"` // Docker Compose project name
	CreatedBy   string            `json:"created_by" gorm:"type:uuid"`
	CreatedAt   time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time        `json:"deleted_at,omitempty" gorm:"index"`
}

// StackStatus represents the status of a stack
type StackStatus string

const (
	StackStatusDeploying StackStatus = "deploying"
	StackStatusRunning   StackStatus = "running"
	StackStatusStopped   StackStatus = "stopped"
	StackStatusFailed    StackStatus = "failed"
	StackStatusUpdating  StackStatus = "updating"
)

// StackService represents a service within a stack
type StackService struct {
	ID          string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StackID     string     `json:"stack_id" gorm:"type:uuid;not null;index"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null"`
	Image       string     `json:"image" gorm:"type:varchar(500)"`
	Replicas    int        `json:"replicas" gorm:"type:int;default:1"`
	Status      string     `json:"status" gorm:"type:varchar(50)"`
	Ports       []string   `json:"ports" gorm:"type:jsonb"`
	Environment []string   `json:"environment" gorm:"type:jsonb"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName specifies the table name for DockerStack
func (DockerStack) TableName() string {
	return "docker_stacks"
}

// TableName specifies the table name for StackService
func (StackService) TableName() string {
	return "stack_services"
}

// StackDeployRequest represents a request to deploy a stack
type StackDeployRequest struct {
	Name        string            `json:"name" binding:"required" example:"my-app"`
	ComposeFile string            `json:"compose_file" binding:"required"` // YAML content
	EnvVars     map[string]string `json:"env_vars,omitempty"`
	DockerHost  string            `json:"docker_host,omitempty" example:"unix:///var/run/docker.sock"`
}

// StackUpdateRequest represents a request to update a stack
type StackUpdateRequest struct {
	ComposeFile string            `json:"compose_file" binding:"required"`
	EnvVars     map[string]string `json:"env_vars,omitempty"`
}

// StackInfo represents detailed stack information
type StackInfo struct {
	Stack    DockerStack    `json:"stack"`
	Services []StackService `json:"services"`
}
