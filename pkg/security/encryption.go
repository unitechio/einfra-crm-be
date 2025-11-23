package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var (
	// ErrInvalidKey is returned when the encryption key is invalid
	ErrInvalidKey = errors.New("invalid encryption key")
	// ErrInvalidCiphertext is returned when the ciphertext is invalid or tampered
	ErrInvalidCiphertext = errors.New("invalid or tampered ciphertext")
	// ErrEmptyPlaintext is returned when trying to encrypt empty data
	ErrEmptyPlaintext = errors.New("plaintext cannot be empty")
)

// EncryptionService provides encryption and decryption functionality
type EncryptionService interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
	EncryptBytes(plaintext []byte) ([]byte, error)
	DecryptBytes(ciphertext []byte) ([]byte, error)
}

// AESEncryption implements EncryptionService using AES-256-GCM
type AESEncryption struct {
	key []byte
}

// NewAESEncryption creates a new AES encryption service
// The key should be 32 bytes for AES-256
func NewAESEncryption(masterKey string) (*AESEncryption, error) {
	if masterKey == "" {
		return nil, ErrInvalidKey
	}

	// Derive a 32-byte key from the master key using SHA-256
	hash := sha256.Sum256([]byte(masterKey))

	return &AESEncryption{
		key: hash[:],
	}, nil
}

// Encrypt encrypts plaintext using AES-256-GCM and returns base64-encoded ciphertext
func (e *AESEncryption) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", ErrEmptyPlaintext
	}

	ciphertext, err := e.EncryptBytes([]byte(plaintext))
	if err != nil {
		return "", err
	}

	// Encode to base64 for safe database storage
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts base64-encoded ciphertext using AES-256-GCM
func (e *AESEncryption) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	// Decode from base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	plaintext, err := e.DecryptBytes(data)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// EncryptBytes encrypts plaintext bytes using AES-256-GCM
// Format: [nonce][ciphertext][tag]
func (e *AESEncryption) EncryptBytes(plaintext []byte) ([]byte, error) {
	if len(plaintext) == 0 {
		return nil, ErrEmptyPlaintext
	}

	// Create AES cipher
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt and authenticate
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

// DecryptBytes decrypts ciphertext bytes using AES-256-GCM
func (e *AESEncryption) DecryptBytes(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) == 0 {
		return nil, nil
	}

	// Create AES cipher
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Check minimum ciphertext length
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrInvalidCiphertext
	}

	// Extract nonce and ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt and verify
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrInvalidCiphertext
	}

	return plaintext, nil
}

// GenerateKey generates a secure random encryption key
func GenerateKey() (string, error) {
	key := make([]byte, 32) // 256 bits
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", fmt.Errorf("failed to generate key: %w", err)
	}
	return base64.StdEncoding.EncodeToString(key), nil
}
