package domain

import (
	"time"
)

// AlertRule represents a rule for resource alerting
type AlertRule struct {
	ID          string    `json:"id"`
	ContainerID string    `json:"container_id"`
	Metric      string    `json:"metric"` // cpu, memory
	Threshold   float64   `json:"threshold"`
	Duration    string    `json:"duration"` // e.g., "5m"
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AlertHistory represents a history of triggered alerts
type AlertHistory struct {
	ID          string    `json:"id"`
	RuleID      string    `json:"rule_id"`
	ContainerID string    `json:"container_id"`
	Metric      string    `json:"metric"`
	Value       float64   `json:"value"`
	Threshold   float64   `json:"threshold"`
	TriggeredAt time.Time `json:"triggered_at"`
	ResolvedAt  time.Time `json:"resolved_at,omitempty"`
}

// DefaultAlertRules returns default rules applied to all containers
func DefaultAlertRules() []AlertRule {
	return []AlertRule{
		{
			ID:        "default-cpu-high",
			Metric:    "cpu",
			Threshold: 80.0, // 80%
			Duration:  "1m",
			Enabled:   true,
		},
		{
			ID:        "default-memory-high",
			Metric:    "memory",
			Threshold: 90.0, // 90%
			Duration:  "1m",
			Enabled:   true,
		},
	}
}
