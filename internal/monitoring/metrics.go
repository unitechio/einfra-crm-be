package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus metrics for the application.
type Metrics struct {
	// HTTP Metrics
	HTTPRequestsTotal   *prometheus.CounterVec
	HTTPRequestDuration *prometheus.HistogramVec

	// Messaging Metrics
	MessagesPublishedTotal *prometheus.CounterVec
	MessagesConsumedTotal  *prometheus.CounterVec
	MessageProcessingTime  *prometheus.HistogramVec
	MessagesToDLQTotal     *prometheus.CounterVec
}

// NewMetrics creates and registers the Prometheus metrics.
func NewMetrics(reg prometheus.Registerer) *Metrics {
	return &Metrics{
		// HTTP Metrics
		HTTPRequestsTotal: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests.",
			},
			[]string{"method", "path", "status"},
		),
		HTTPRequestDuration: promauto.With(reg).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),

		// Messaging Metrics
		MessagesPublishedTotal: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "messaging_published_total",
				Help: "Total number of messages published.",
			},
			[]string{"topic", "status"}, // status can be "success" or "error"
		),
		MessagesConsumedTotal: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "messaging_consumed_total",
				Help: "Total number of messages consumed.",
			},
			[]string{"topic", "status"}, // status can be "success", "error", or "dlq"
		),
		MessageProcessingTime: promauto.With(reg).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "messaging_processing_duration_seconds",
				Help:    "Histogram of message processing time.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"topic"},
		),
		MessagesToDLQTotal: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "messaging_dlq_total",
				Help: "Total number of messages sent to the Dead Letter Queue.",
			},
			[]string{"topic"},
		),
	}
}
