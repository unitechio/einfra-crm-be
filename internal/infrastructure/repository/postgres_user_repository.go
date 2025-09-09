
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"mymodule/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresUserRepository is an implementation of UserRepository for PostgreSQL.
type PostgresUserRepository struct {
	db *pgxpool.Pool
}

// NewPostgresUserRepository creates a new PostgresUserRepository.
func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// Create inserts a new user into the database.
func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	// Set default role if not provided
	if user.Role == "" {
		user.Role = domain.UserRole
	}
	query := `INSERT INTO users (username, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRow(ctx, query, user.Username, user.Email, user.Password, user.Role).Scan(&user.ID)
}

// GetByEmail retrieves a user by their email address.
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, role, settings FROM users WHERE email = $1`
	user := &domain.User{}
	var settingsJSON []byte
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &settingsJSON)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(settingsJSON, &user.Settings); err != nil {
		return nil, err
	}
	return user, nil
}

// GetByID retrieves a user by their ID.
func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, role, settings FROM users WHERE id = $1`
	user := &domain.User{}
	var settingsJSON []byte
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &settingsJSON)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(settingsJSON, &user.Settings); err != nil {
		return nil, err
	}
	return user, nil
}

// GetAll retrieves all users from the database.
func (r *PostgresUserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	query := `SELECT id, username, email, role FROM users ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, &user)
	}

	return users, nil
}

// CreateBatch efficiently inserts multiple users into the database using a batch operation.
func (r *PostgresUserRepository) CreateBatch(ctx context.Context, users []*domain.User) error {
	batch := &pgx.Batch{}
	for _, user := range users {
		if user.Role == "" {
			user.Role = domain.UserRole
		}
		query := `INSERT INTO users (username, email, password_hash, role) VALUES ($1, $2, $3, $4) ON CONFLICT (email) DO NOTHING`
		batch.Queue(query, user.Username, user.Email, user.Password, user.Role)
	}

	br := r.db.SendBatch(ctx, batch)
	defer br.Close()

	// We need to check the result of each command in the batch.
	for i := 0; i < len(users); i++ {
		_, err := br.Exec()
		if err != nil {
			// Log or handle the error for the specific user insert
			// For now, we'll just return the first error we encounter
			return fmt.Errorf("error in batch insert on user %s: %w", users[i].Email, err)
		}
	}

	return nil
}

// UpdateSettings updates the user's settings in the database.
func (r *PostgresUserRepository) UpdateSettings(ctx context.Context, userID string, settings domain.UserSettings) error {
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	query := `UPDATE users SET settings = $1 WHERE id = $2`
	_, err = r.db.Exec(ctx, query, settingsJSON, userID)
	if err != nil {
		return fmt.Errorf("failed to update user settings: %w", err)
	}

	return nil
}
