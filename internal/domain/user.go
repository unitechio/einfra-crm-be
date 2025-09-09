
package domain

import (
	"context"

	"github.com/xuri/excelize/v2"
)

// Role defines the user roles in the system.
type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

// User represents a user in the system.
type User struct {
	ID       string       `json:"id"`
	Username string       `json:"username"`
	Email    string       `json:"email"`
	Password string       `json:"-"` // Password should not be exposed
	Role     Role         `json:"role"`
	Settings UserSettings `json:"settings"`
}

// UserRepository defines the interface for user data storage.
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	CreateBatch(ctx context.Context, users []*User) error
	UpdateSettings(ctx context.Context, userID string, settings UserSettings) error
}

// UserUsecase defines the business logic for user-related operations.
type UserUsecase interface {
	ImportUsersFromExcel(ctx context.Context, filePath string) error
	ExportUsersToExcel(ctx context.Context) (*excelize.File, string, error)
	UpdateUserSettings(ctx context.Context, userID string, settings UserSettings) error
}
