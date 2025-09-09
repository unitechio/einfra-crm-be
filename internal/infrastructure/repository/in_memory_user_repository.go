
package repository

import (
	"context"
	"fmt"
	"mymodule/internal/domain"
	"sync"

	"github.com/google/uuid"
)

// InMemoryUserRepository is a simple in-memory implementation of UserRepository.
type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

// NewInMemoryUserRepository creates a new InMemoryUserRepository.
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

// Create saves a new user.
func (r *InMemoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.Email]; exists {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	user.ID = uuid.New().String()
	r.users[user.Email] = user
	return nil
}

// GetByEmail retrieves a user by email.
func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[email]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// GetByID retrieves a user by ID.
func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// GetAll retrieves all users.
func (r *InMemoryUserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    allUsers := make([]*domain.User, 0, len(r.users))
    for _, user := range r.users {
        allUsers = append(allUsers, user)
    }
    return allUsers, nil
}


// CreateBatch saves multiple users.
func (r *InMemoryUserRepository) CreateBatch(ctx context.Context, users []*domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, user := range users {
		if _, exists := r.users[user.Email]; exists {
			// In a real app, you might want to handle this differently (e.g., return an error for the specific user)
			// For simplicity, we'll just skip existing users.
			continue
		}
		user.ID = uuid.New().String()
		r.users[user.Email] = user
	}
	return nil
}
