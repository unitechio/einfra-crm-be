
package util

import (
	"github.com/google/uuid"
)

// NewUUID generates a new UUID and returns it as a string.
func NewUUID() string {
	return uuid.New().String()
}
