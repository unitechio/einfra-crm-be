package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"mymodule/internal/domain"
	"time"

	"github.com/google/uuid"
)

type auditRepository struct {
	kafkaProducer *KafkaProducer
	// In-memory store for demonstration purposes
	store map[string]domain.AuditEntry
}

func NewAuditRepository(kafkaProducer *KafkaProducer) domain.AuditRepository {
	return &auditRepository{
		kafkaProducer: kafkaProducer,
		store:         make(map[string]domain.AuditEntry),
	}
}

func (r *auditRepository) Add(ctx context.Context, entry domain.AuditEntry) (domain.AuditEntry, error) {
	entry.ID = uuid.New().String()
	entry.Timestamp = time.Now()

	payload, err := json.Marshal(entry)
	if err != nil {
		return domain.AuditEntry{}, err
	}

	if err := r.kafkaProducer.ProduceMessage("audit-log", payload); err != nil {
		return domain.AuditEntry{}, err
	}

	r.store[entry.ID] = entry
	return entry, nil
}

func (r *auditRepository) GetAll(ctx context.Context) ([]domain.AuditEntry, error) {
	audits := make([]domain.AuditEntry, 0, len(r.store))
	for _, audit := range r.store {
		audits = append(audits, audit)
	}
	return audits, nil
}

func (r *auditRepository) GetByID(ctx context.Context, id string) (domain.AuditEntry, error) {
	audit, ok := r.store[id]
	if !ok {
		return domain.AuditEntry{}, fmt.Errorf("audit entry with id %s not found", id)
	}
	return audit, nil
}

func (r *auditRepository) Update(ctx context.Context, id string, entry domain.AuditEntry) (domain.AuditEntry, error) {
	_, ok := r.store[id]
	if !ok {
		return domain.AuditEntry{}, fmt.Errorf("audit entry with id %s not found", id)
	}
	entry.ID = id
	r.store[id] = entry
	return entry, nil
}

func (r *auditRepository) Delete(ctx context.Context, id string) error {
	_, ok := r.store[id]
	if !ok {
		return fmt.Errorf("audit entry with id %s not found", id)
	}
	delete(r.store, id)
	return nil
}
