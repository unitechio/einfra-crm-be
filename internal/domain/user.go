package domain

import (
	"context"
	"time"
)

// Role represents a user role in the system
// @Description User role with permissions
type Role struct {
	ID          string       `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string       `json:"name" gorm:"type:varchar(100);not null;uniqueIndex" validate:"required,min=3,max=100" example:"admin"`
	DisplayName string       `json:"display_name" gorm:"type:varchar(255);not null" validate:"required" example:"Administrator"`
	Description string       `json:"description" gorm:"type:text" example:"Full system access"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
	IsSystem    bool         `json:"is_system" gorm:"type:boolean;default:false" example:"false"` // System roles cannot be deleted
	IsActive    bool         `json:"is_active" gorm:"type:boolean;default:true;index" example:"true"`
	CreatedAt   time.Time    `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time    `json:"updated_at" gorm:"autoUpdateTime" example:"2024-01-01T00:00:00Z"`
	DeletedAt   *time.Time   `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string" example:"2024-01-01T00:00:00Z"`
}

// TableName specifies the table name for Role model
func (Role) TableName() string {
	return "roles"
}

// Permission represents a system permission
// @Description System permission for access control
type Permission struct {
	ID          string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string     `json:"name" gorm:"type:varchar(100);not null;uniqueIndex" validate:"required" example:"server.create"`
	Resource    string     `json:"resource" gorm:"type:varchar(100);not null;index" validate:"required" example:"server"`
	Action      string     `json:"action" gorm:"type:varchar(50);not null;index" validate:"required" example:"create"`
	Description string     `json:"description" gorm:"type:text" example:"Create new servers"`
	IsSystem    bool       `json:"is_system" gorm:"type:boolean;default:false" example:"false"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime" example:"2024-01-01T00:00:00Z"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string" example:"2024-01-01T00:00:00Z"`
}

// TableName specifies the table name for Permission model
func (Permission) TableName() string {
	return "permissions"
}

// User represents a user in the system
// @Description User entity with authentication and profile information
type User struct {
	ID                string       `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username          string       `json:"username" gorm:"type:varchar(100);not null;uniqueIndex" validate:"required,min=3,max=100" example:"john.doe"`
	Email             string       `json:"email" gorm:"type:varchar(255);not null;uniqueIndex" validate:"required,email" example:"john.doe@example.com"`
	Password          string       `json:"-" gorm:"type:varchar(255);not null"` // Password hash, never exposed in JSON
	FirstName         string       `json:"first_name" gorm:"type:varchar(100)" validate:"max=100" example:"John"`
	LastName          string       `json:"last_name" gorm:"type:varchar(100)" validate:"max=100" example:"Doe"`
	Phone             string       `json:"phone" gorm:"type:varchar(20)" validate:"max=20" example:"+1234567890"`
	Avatar            string       `json:"avatar" gorm:"type:varchar(500)" example:"https://example.com/avatar.jpg"`
	RoleID            string       `json:"role_id" gorm:"type:uuid;index" example:"550e8400-e29b-41d4-a716-446655440000"`
	Role              *Role        `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	IsActive          bool         `json:"is_active" gorm:"type:boolean;default:true;index" example:"true"`
	IsEmailVerified   bool         `json:"is_email_verified" gorm:"type:boolean;default:false" example:"false"`
	EmailVerifiedAt   *time.Time   `json:"email_verified_at,omitempty" swaggertype:"string" example:"2024-01-01T00:00:00Z"`
	LastLoginAt       *time.Time   `json:"last_login_at,omitempty" swaggertype:"string" example:"2024-01-01T00:00:00Z"`
	LastLoginIP       string       `json:"last_login_ip,omitempty" gorm:"type:varchar(45)" example:"192.168.1.1"`
	FailedLoginCount  int          `json:"failed_login_count" gorm:"type:int;default:0" example:"0"`
	LockedUntil       *time.Time   `json:"locked_until,omitempty" swaggertype:"string" example:"2024-01-01T00:00:00Z"`
	PasswordChangedAt *time.Time   `json:"password_changed_at,omitempty" swaggertype:"string" example:"2024-01-01T00:00:00Z"`
	MFAEnabled        bool         `json:"mfa_enabled" gorm:"type:boolean;default:false" example:"false"`
	MFASecret         string       `json:"-" gorm:"type:varchar(255)"` // TOTP secret, never exposed
	Settings          UserSettings `json:"settings" gorm:"type:jsonb"`
	CreatedAt         time.Time    `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
	UpdatedAt         time.Time    `json:"updated_at" gorm:"autoUpdateTime" example:"2024-01-01T00:00:00Z"`
	DeletedAt         *time.Time   `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string" example:"2024-01-01T00:00:00Z"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// IsLocked checks if the user account is locked
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// HasPermission checks if the user has a specific permission
func (u *User) HasPermission(permission string) bool {
	if u.Role == nil {
		return false
	}
	for _, p := range u.Role.Permissions {
		if p.Name == permission {
			return true
		}
	}
	return false
}

// HasRole checks if the user has a specific role
func (u *User) HasRole(roleName string) bool {
	if u.Role == nil {
		return false
	}
	return u.Role.Name == roleName
}

// RefreshToken represents a refresh token for JWT authentication
// @Description Refresh token for renewing access tokens
type RefreshToken struct {
	ID        string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID    string     `json:"user_id" gorm:"type:uuid;not null;index" example:"550e8400-e29b-41d4-a716-446655440000"`
	Token     string     `json:"token" gorm:"type:varchar(500);not null;uniqueIndex" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null;index" example:"2024-01-08T00:00:00Z"`
	IsRevoked bool       `json:"is_revoked" gorm:"type:boolean;default:false;index" example:"false"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" swaggertype:"string" example:"2024-01-01T00:00:00Z"`
}

// TableName specifies the table name for RefreshToken model
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// IsExpired checks if the refresh token is expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsValid checks if the refresh token is valid (not expired and not revoked)
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsExpired() && !rt.IsRevoked
}

// UserFilter represents filtering options for user queries
type UserFilter struct {
	RoleID   string `json:"role_id,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
	Search   string `json:"search,omitempty"` // Search in username, email, first_name, last_name
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

// RoleFilter represents filtering options for role queries
type RoleFilter struct {
	IsActive *bool `json:"is_active,omitempty"`
	Page     int   `json:"page" validate:"min=1"`
	PageSize int   `json:"page_size" validate:"min=1,max=100"`
}

// PermissionFilter represents filtering options for permission queries
type PermissionFilter struct {
	Resource string `json:"resource,omitempty"`
	Action   string `json:"action,omitempty"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

// RoleUsecase defines the business logic for role-related operations
type RoleUsecase interface {
	// CreateRole creates a new role
	CreateRole(ctx context.Context, role *Role) error

	// GetRole retrieves a role by ID
	GetRole(ctx context.Context, id string) (*Role, error)

	// ListRoles retrieves all roles
	ListRoles(ctx context.Context, filter RoleFilter) ([]*Role, int64, error)

	// UpdateRole updates role information
	UpdateRole(ctx context.Context, role *Role) error

	// DeleteRole soft deletes a role
	DeleteRole(ctx context.Context, id string) error

	// AssignPermissions assigns permissions to a role
	AssignPermissions(ctx context.Context, roleID string, permissionIDs []string) error

	// RemovePermissions removes permissions from a role
	RemovePermissions(ctx context.Context, roleID string, permissionIDs []string) error
}

// PermissionUsecase defines the business logic for permission-related operations
type PermissionUsecase interface {
	// CreatePermission creates a new permission
	CreatePermission(ctx context.Context, permission *Permission) error

	// GetPermission retrieves a permission by ID
	GetPermission(ctx context.Context, id string) (*Permission, error)

	// ListPermissions retrieves all permissions
	ListPermissions(ctx context.Context, filter PermissionFilter) ([]*Permission, int64, error)

	// UpdatePermission updates permission information
	UpdatePermission(ctx context.Context, permission *Permission) error

	// DeletePermission soft deletes a permission
	DeletePermission(ctx context.Context, id string) error
}
