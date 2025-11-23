package domain

import (
	"context"
	"time"
)

// PermissionScope defines the scope of a permission
type PermissionScope string

const (
	// PermissionScopeGlobal applies to all resources
	PermissionScopeGlobal PermissionScope = "global"
	// PermissionScopeEnvironment applies to specific environment
	PermissionScopeEnvironment PermissionScope = "environment"
	// PermissionScopeResource applies to specific resource
	PermissionScopeResource PermissionScope = "resource"
)

// ResourceType represents the type of resource for permissions
type ResourceType string

const (
	// ResourceTypeServer represents server resources
	ResourceTypeServer ResourceType = "server"
	// ResourceTypeK8sCluster represents Kubernetes cluster resources
	ResourceTypeK8sCluster ResourceType = "k8s_cluster"
	// ResourceTypeK8sNamespace represents Kubernetes namespace resources
	ResourceTypeK8sNamespace ResourceType = "k8s_namespace"
	// ResourceTypeDockerContainer represents Docker container resources
	ResourceTypeDockerContainer ResourceType = "docker_container"
	// ResourceTypeHarborProject represents Harbor project resources
	ResourceTypeHarborProject ResourceType = "harbor_project"
)

// UserEnvironmentRole represents a user's role in a specific environment
// @Description User role assignment for a specific environment
type UserEnvironmentRole struct {
	ID            string       `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID        string       `json:"user_id" gorm:"type:uuid;not null;index" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	User          *User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	EnvironmentID *string      `json:"environment_id,omitempty" gorm:"type:uuid;index" example:"550e8400-e29b-41d4-a716-446655440000"` // NULL = all environments
	Environment   *Environment `json:"environment,omitempty" gorm:"foreignKey:EnvironmentID"`
	RoleID        string       `json:"role_id" gorm:"type:uuid;not null;index" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Role          *Role        `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	CreatedBy     string       `json:"created_by" gorm:"type:uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	CreatedAt     time.Time    `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
	UpdatedAt     time.Time    `json:"updated_at" gorm:"autoUpdateTime" example:"2024-01-01T00:00:00Z"`
	DeletedAt     *time.Time   `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string" example:"2024-01-01T00:00:00Z"`
}

// TableName specifies the table name for UserEnvironmentRole model
func (UserEnvironmentRole) TableName() string {
	return "user_environment_roles"
}

// ResourcePermission represents a user's permission on a specific resource
// @Description User permission for a specific resource
type ResourcePermission struct {
	ID            string       `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID        string       `json:"user_id" gorm:"type:uuid;not null;index" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	User          *User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ResourceType  ResourceType `json:"resource_type" gorm:"type:varchar(50);not null;index" validate:"required" example:"k8s_cluster"`
	ResourceID    string       `json:"resource_id" gorm:"type:varchar(255);not null;index" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Actions       []string     `json:"actions" gorm:"type:jsonb;not null" validate:"required" example:"read,update,delete"`            // ["read", "update", "delete"]
	EnvironmentID *string      `json:"environment_id,omitempty" gorm:"type:uuid;index" example:"550e8400-e29b-41d4-a716-446655440000"` // NULL = all environments
	Environment   *Environment `json:"environment,omitempty" gorm:"foreignKey:EnvironmentID"`
	ExpiresAt     *time.Time   `json:"expires_at,omitempty" swaggertype:"string" example:"2024-12-31T23:59:59Z"` // NULL = never expires
	GrantedBy     string       `json:"granted_by" gorm:"type:uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Reason        string       `json:"reason" gorm:"type:text" example:"Temporary access for deployment"`
	CreatedAt     time.Time    `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
	UpdatedAt     time.Time    `json:"updated_at" gorm:"autoUpdateTime" example:"2024-01-01T00:00:00Z"`
	DeletedAt     *time.Time   `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string" example:"2024-01-01T00:00:00Z"`
}

// TableName specifies the table name for ResourcePermission model
func (ResourcePermission) TableName() string {
	return "resource_permissions"
}

// IsExpired checks if the resource permission has expired
func (rp *ResourcePermission) IsExpired() bool {
	if rp.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*rp.ExpiresAt)
}

// HasAction checks if the permission includes a specific action
func (rp *ResourcePermission) HasAction(action string) bool {
	for _, a := range rp.Actions {
		if a == action {
			return true
		}
	}
	return false
}

// GrantPermissionRequest represents a request to grant permission to a user
// @Description Request to grant resource permission to a user
type GrantPermissionRequest struct {
	UserID        string     `json:"user_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	ResourceType  string     `json:"resource_type" validate:"required,oneof=server k8s_cluster k8s_namespace docker_container harbor_project" example:"k8s_cluster"`
	ResourceID    string     `json:"resource_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Actions       []string   `json:"actions" validate:"required,min=1" example:"read,update"`
	EnvironmentID *string    `json:"environment_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
	Reason        string     `json:"reason" validate:"max=500" example:"Temporary access for deployment"`
}

// RevokePermissionRequest represents a request to revoke permission from a user
// @Description Request to revoke resource permission from a user
type RevokePermissionRequest struct {
	PermissionID string `json:"permission_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Reason       string `json:"reason" validate:"max=500" example:"Access no longer needed"`
}

// UserPermissions represents all permissions for a user
// @Description Complete permission set for a user
type UserPermissions struct {
	UserID               string                 `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	GlobalRoles          []*Role                `json:"global_roles"`
	EnvironmentRoles     []*UserEnvironmentRole `json:"environment_roles"`
	ResourcePermissions  []*ResourcePermission  `json:"resource_permissions"`
	EffectivePermissions []string               `json:"effective_permissions"` // Computed list of all permission names
}

// AuthorizationUsecase defines the business logic for authorization
type AuthorizationUsecase interface {
	// Permission checking
	CheckPermission(ctx context.Context, userID, permission string) (bool, error)
	CheckEnvironmentPermission(ctx context.Context, userID, permission, environmentID string) (bool, error)
	CheckResourcePermission(ctx context.Context, userID string, resourceType ResourceType, resourceID, action string) (bool, error)
	CheckNamespacePermission(ctx context.Context, userID, clusterID, namespace, action string) (bool, error)

	// Grant/Revoke permissions
	GrantResourcePermission(ctx context.Context, req *GrantPermissionRequest, grantedBy string) error
	RevokeResourcePermission(ctx context.Context, req *RevokePermissionRequest, revokedBy string) error

	// Assign environment role
	AssignEnvironmentRole(ctx context.Context, userID, roleID string, environmentID *string, assignedBy string) error
	RemoveEnvironmentRole(ctx context.Context, id string) error

	// List permissions
	ListUserPermissions(ctx context.Context, userID string) (*UserPermissions, error)
	ListResourcePermissions(ctx context.Context, resourceType ResourceType, resourceID string) ([]*ResourcePermission, error)

	// Cleanup
	CleanupExpiredPermissions(ctx context.Context) (int64, error)
}
