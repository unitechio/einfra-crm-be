package domain

import (
	"time"
)

// AuthProvider represents different authentication providers
type AuthProvider string

const (
	// AuthProviderLocal represents local username/password authentication
	AuthProviderLocal AuthProvider = "local"
	// AuthProviderGoogle represents Google OAuth authentication
	AuthProviderGoogle AuthProvider = "google"
	// AuthProviderAzure represents Azure AD authentication
	AuthProviderAzure AuthProvider = "azure"
	// AuthProviderGitHub represents GitHub OAuth authentication
	AuthProviderGitHub AuthProvider = "github"
	// AuthProviderLDAP represents LDAP authentication
	AuthProviderLDAP AuthProvider = "ldap"
)

// TokenType represents the type of token
type TokenType string

const (
	// TokenTypeAccess represents an access token
	TokenTypeAccess TokenType = "access"
	// TokenTypeRefresh represents a refresh token
	TokenTypeRefresh TokenType = "refresh"
	// TokenTypePasswordReset represents a password reset token
	TokenTypePasswordReset TokenType = "password_reset"
	// TokenTypeEmailVerification represents an email verification token
	TokenTypeEmailVerification TokenType = "email_verification"
)

// AuthCredentials represents login credentials
// @Description Authentication credentials for login
type AuthCredentials struct {
	Username string `json:"username" validate:"required" example:"john.doe"`
	Password string `json:"password" validate:"required" example:"SecurePass123!"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=3,max=50"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   int    `json:"role_id" `
}

// AuthResponse represents the authentication response
// @Description Authentication response with tokens and user info
type AuthResponse struct {
	AccessToken  string    `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string    `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType    string    `json:"token_type" example:"Bearer"`
	ExpiresIn    int       `json:"expires_in" example:"3600"` // seconds
	User         *User     `json:"user"`
	IssuedAt     time.Time `json:"issued_at" example:"2024-01-01T00:00:00Z"`
}

// TokenClaims represents JWT token claims
// @Description JWT token claims
type TokenClaims struct {
	UserID      string    `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username    string    `json:"username" example:"john.doe"`
	Email       string    `json:"email" example:"john.doe@example.com"`
	RoleID      string    `json:"role_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	RoleName    string    `json:"role_name" example:"admin"`
	Permissions []string  `json:"permissions" example:"server.create,server.read"`
	TokenType   TokenType `json:"token_type" example:"access"`
	IssuedAt    time.Time `json:"iat" example:"2024-01-01T00:00:00Z"`
	ExpiresAt   time.Time `json:"exp" example:"2024-01-01T01:00:00Z"`
}

// PasswordResetRequest represents a password reset request
// @Description Password reset request
type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email" example:"john.doe@example.com"`
}

// PasswordResetConfirm represents password reset confirmation
// @Description Password reset confirmation with new password
type PasswordResetConfirm struct {
	Token       string `json:"token" validate:"required" example:"reset-token-123"`
	NewPassword string `json:"new_password" validate:"required,min=8" example:"NewSecurePass123!"`
}

// EmailVerificationRequest represents email verification request
// @Description Email verification request
type EmailVerificationRequest struct {
	Token string `json:"token" validate:"required" example:"verify-token-123"`
}

// ChangePasswordRequest represents a password change request
// @Description Change password request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required" example:"OldPass123!"`
	NewPassword string `json:"new_password" validate:"required,min=8" example:"NewPass123!"`
}

// OAuthState represents OAuth state for CSRF protection
type OAuthState struct {
	State     string       `json:"state"`
	Provider  AuthProvider `json:"provider"`
	CreatedAt time.Time    `json:"created_at"`
	ExpiresAt time.Time    `json:"expires_at"`
}

// Session represents a user session
// @Description User session information
type Session struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID       string    `json:"user_id" gorm:"type:uuid;not null;index" example:"550e8400-e29b-41d4-a716-446655440000"`
	Token        string    `json:"token" gorm:"type:varchar(500);not null;uniqueIndex" example:"session-token-123"`
	IPAddress    string    `json:"ip_address" gorm:"type:varchar(45)" example:"192.168.1.1"`
	UserAgent    string    `json:"user_agent" gorm:"type:text" example:"Mozilla/5.0..."`
	DeviceType   string    `json:"device_type" gorm:"type:varchar(50)" example:"desktop"`
	Location     string    `json:"location" gorm:"type:varchar(255)" example:"New York, US"`
	IsActive     bool      `json:"is_active" gorm:"type:boolean;default:true;index" example:"true"`
	LastActivity time.Time `json:"last_activity" gorm:"index" example:"2024-01-01T00:00:00Z"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"not null;index" example:"2024-01-08T00:00:00Z"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
}

// TableName specifies the table name for Session model
func (Session) TableName() string {
	return "sessions"
}

// IsExpired checks if the session is expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsValid checks if the session is valid (active and not expired)
func (s *Session) IsValid() bool {
	return s.IsActive && !s.IsExpired()
}

// LoginAttempt represents a login attempt for security tracking
// @Description Login attempt tracking for security
type LoginAttempt struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username   string    `json:"username" gorm:"type:varchar(100);index" example:"john.doe"`
	IPAddress  string    `json:"ip_address" gorm:"type:varchar(45);index" example:"192.168.1.1"`
	UserAgent  string    `json:"user_agent" gorm:"type:text" example:"Mozilla/5.0..."`
	Success    bool      `json:"success" gorm:"type:boolean;index" example:"true"`
	FailReason string    `json:"fail_reason,omitempty" gorm:"type:varchar(255)" example:"Invalid password"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime;index" example:"2024-01-01T00:00:00Z"`
}

// TableName specifies the table name for LoginAttempt model
func (LoginAttempt) TableName() string {
	return "login_attempts"
}
