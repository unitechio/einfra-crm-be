package domain

import "context"

type HealthStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type HealthRepository interface {
	Check(ctx context.Context) []HealthStatus
}
