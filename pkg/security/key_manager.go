package security

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
)

var (
	// ErrKeyNotFound is returned when encryption key is not configured
	ErrKeyNotFound = errors.New("encryption key not found in environment")
	// ErrInvalidKeyVersion is returned when key version is invalid
	ErrInvalidKeyVersion = errors.New("invalid encryption key version")
)

// KeyManager manages encryption keys with versioning support
type KeyManager struct {
	currentKey     string
	currentVersion int
	keys           map[int]string // version -> key
	mu             sync.RWMutex
}

// NewKeyManager creates a new key manager
func NewKeyManager() *KeyManager {
	return &KeyManager{
		keys: make(map[int]string),
	}
}

// LoadFromEnv loads encryption key from environment variables
func (km *KeyManager) LoadFromEnv() error {
	key := os.Getenv("ENCRYPTION_KEY")
	if key == "" {
		return ErrKeyNotFound
	}

	versionStr := os.Getenv("ENCRYPTION_KEY_VERSION")
	version := 1 // Default version
	if versionStr != "" {
		v, err := strconv.Atoi(versionStr)
		if err != nil {
			return fmt.Errorf("invalid key version: %w", err)
		}
		version = v
	}

	km.mu.Lock()
	defer km.mu.Unlock()

	km.currentKey = key
	km.currentVersion = version
	km.keys[version] = key

	return nil
}

// GetCurrentKey returns the current encryption key
func (km *KeyManager) GetCurrentKey() (string, int, error) {
	km.mu.RLock()
	defer km.mu.RUnlock()

	if km.currentKey == "" {
		return "", 0, ErrKeyNotFound
	}

	return km.currentKey, km.currentVersion, nil
}

// GetKey returns the encryption key for a specific version
func (km *KeyManager) GetKey(version int) (string, error) {
	km.mu.RLock()
	defer km.mu.RUnlock()

	key, exists := km.keys[version]
	if !exists {
		return "", ErrInvalidKeyVersion
	}

	return key, nil
}

// AddKey adds a new encryption key version
func (km *KeyManager) AddKey(version int, key string) error {
	if key == "" {
		return ErrInvalidKey
	}

	km.mu.Lock()
	defer km.mu.Unlock()

	km.keys[version] = key

	// Update current key if this is a newer version
	if version > km.currentVersion {
		km.currentKey = key
		km.currentVersion = version
	}

	return nil
}

// GetCurrentVersion returns the current key version
func (km *KeyManager) GetCurrentVersion() int {
	km.mu.RLock()
	defer km.mu.RUnlock()
	return km.currentVersion
}

// ValidateKey validates that a key meets security requirements
func ValidateKey(key string) error {
	if key == "" {
		return ErrInvalidKey
	}

	// Key should be at least 32 characters (256 bits when base64 decoded)
	if len(key) < 32 {
		return fmt.Errorf("encryption key too short (minimum 32 characters)")
	}

	return nil
}
