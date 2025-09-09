package domain

import "context"

// HealthStatus represents the status of a single component.
type HealthStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// HealthRepository defines the interface for health checks.
type HealthRepository interface {
	Check(ctx context.Context) []HealthStatus
}
