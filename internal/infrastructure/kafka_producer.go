package infrastructure

import (
	"context"
	"encoding/json"
	"log"
	"mymodule/internal/domain"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// KafkaProducer is a Kafka producer that sends notification messages.
type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

// NewKafkaProducer creates a new KafkaProducer.
func NewKafkaProducer(brokers, topic string) (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		return nil, err
	}

	// Go-routine to handle delivery reports
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	return &KafkaProducer{producer: p, topic: topic}, nil
}

// SendNotification sends a notification message to Kafka.
func (kp *KafkaProducer) SendNotification(ctx context.Context, notification *domain.Notification) error {
	value, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kp.topic, Partition: kafka.PartitionAny},
		Value:          value,
	}

	// Produce the message
	deliveryChan := make(chan kafka.Event)
	err = kp.producer.Produce(message, deliveryChan)
	if err != nil {
		return err
	}

	// Wait for the delivery report
	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	} 

	close(deliveryChan)

	return nil
}

// Close closes the Kafka producer.
func (kp *KafkaProducer) Close() {
	// Wait for message deliveries before shutting down
	kp.producer.Flush(15 * 1000)
	kp.producer.Close()
}
