package monitoring

import (
	"context"
	"time"

	"mymodule/internal/messaging"
)

// --- Publisher Middleware ---

type publisherMetricsMiddleware struct {
	next    messaging.Publisher
	metrics *Metrics
}

// NewPublisherMetricsMiddleware creates a new publisher middleware that records metrics.
func NewPublisherMetricsMiddleware(next messaging.Publisher, metrics *Metrics) messaging.Publisher {
	return &publisherMetricsMiddleware{
		next:    next,
		metrics: metrics,
	}
}

func (m *publisherMetricsMiddleware) Publish(ctx context.Context, msg messaging.Message) error {
	status := "success"
	err := m.next.Publish(ctx, msg)
	if err != nil {
		status = "error"
	}

	m.metrics.MessagesPublishedTotal.WithLabelValues(msg.Topic, status).Inc()

	if msg.Topic == "dead-letter-queue" { // Assuming this is your DLQ topic name
		m.metrics.MessagesToDLQTotal.WithLabelValues(msg.Topic).Inc()
	}

	return err
}

func (m *publisherMetricsMiddleware) Close() {
	m.next.Close()
}

// --- Subscriber Middleware ---

type subscriberMetricsMiddleware struct {
	next    messaging.Subscriber
	metrics *Metrics
}

// NewSubscriberMetricsMiddleware creates a new subscriber middleware that records metrics.
func NewSubscriberMetricsMiddleware(next messaging.Subscriber, metrics *Metrics) messaging.Subscriber {
	return &subscriberMetricsMiddleware{
		next:    next,
		metrics: metrics,
	}
}

func (m *subscriberMetricsMiddleware) Subscribe(ctx context.Context, topic string, handler messaging.MessageHandler) error {
	// Wrap the original handler to record metrics
	metricsHandler := func(ctx context.Context, msg messaging.Message) error {
		startTime := time.Now()

		err := handler(ctx, msg)

		status := "success"
		if err != nil {
			status = "error"
		}

		duration := time.Since(startTime).Seconds()
		m.metrics.MessageProcessingTime.WithLabelValues(topic).Observe(duration)
		m.metrics.MessagesConsumedTotal.WithLabelValues(topic, status).Inc()

		return err
	}

	// Pass the wrapped handler to the actual subscriber
	return m.next.Subscribe(ctx, topic, metricsHandler)
}

func (m *subscriberMetricsMiddleware) Close() {
	m.next.Close()
}
