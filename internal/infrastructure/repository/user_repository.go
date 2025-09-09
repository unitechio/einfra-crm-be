
package repository

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"mymodule/internal/domain"
)

// inMemoryUserRepository is a mock implementation of UserRepository.
// In a real application, this would interact with a database.
type inMemoryUserRepository struct {
	users map[string]*domain.User
	mutex sync.RWMutex
}

// NewInMemoryUserRepository creates a new inMemoryUserRepository.
func NewInMemoryUserRepository() domain.UserRepository {
	// Create a mock user.
	users := make(map[string]*domain.User)
	users["1"] = &domain.User{
		ID:           "1",
		Username:     "testuser",
		Role:         "user",
		PasswordHash: "$2a$10$...", // In a real app, this would be a hash
		AuthProvider: domain.AuthProviderLocal,
	}

	return &inMemoryUserRepository{users: users}
}

// FindByUsername finds a user by their username, only for local accounts.
func (r *inMemoryUserRepository) FindByUsername(username string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Username == username && user.AuthProvider == domain.AuthProviderLocal {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// FindByID finds a user by their ID.
func (r *inMemoryUserRepository) FindByID(id string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// FindByProviderID finds a user by their authentication provider and provider-specific ID.
func (r *inMemoryUserRepository) FindByProviderID(provider domain.AuthProvider, providerID string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.AuthProvider == provider && user.AuthProviderID == providerID {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// Create creates a new user in the repository.
func (r *inMemoryUserRepository) Create(user *domain.User) (*domain.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Simple check to prevent username collision, especially for local accounts.
	// For OAuth, we should primarily rely on FindByProviderID.
	if user.AuthProvider == domain.AuthProviderLocal {
		for _, existingUser := range r.users {
			if existingUser.Username == user.Username && existingUser.AuthProvider == domain.AuthProviderLocal {
				return nil, errors.New("username already exists")
			}
		}
	}

	// Generate a new ID for the user.
	newID := strconv.FormatInt(time.Now().UnixNano(), 10)
	user.ID = newID

	r.users[newID] = user
	return user, nil
}
