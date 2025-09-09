
package infrastructure

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"mymodule/internal/domain"
)

// healthRepository is the implementation of HealthRepository.
type healthRepository struct {
	db          *gorm.DB
	redisClient *redis.Client
	minioClient *minio.Client
}

// NewHealthRepository creates a new HealthRepository.
func NewHealthRepository(db *gorm.DB, redisClient *redis.Client, minioClient *minio.Client) domain.HealthRepository {
	return &healthRepository{
		db:          db,
		redisClient: redisClient,
		minioClient: minioClient,
	}
}

// Check performs health checks on all dependencies concurrently.
func (r *healthRepository) Check(ctx context.Context) []domain.HealthStatus {
	var wg sync.WaitGroup
	statusChan := make(chan domain.HealthStatus, 3) // Buffer for 3 checks

	// Check MySQL.
	wg.Add(1)
	go func() {
		defer wg.Done()
		sqlDB, err := r.db.DB()
		if err != nil {
			statusChan <- domain.HealthStatus{Name: "mysql", Status: "down", Error: err.Error()}
			return
		}
		if err := sqlDB.PingContext(ctx); err != nil {
			statusChan <- domain.HealthStatus{Name: "mysql", Status: "down", Error: err.Error()}
		} else {
			statusChan <- domain.HealthStatus{Name: "mysql", Status: "up"}
		}
	}()

	// Check Redis.
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := r.redisClient.Ping(ctx).Err(); err != nil {
			statusChan <- domain.HealthStatus{Name: "redis", Status: "down", Error: err.Error()}
		} else {
			statusChan <- domain.HealthStatus{Name: "redis", Status: "up"}
		}
	}()

	// Check Minio.
	wg.Add(1)
	go func() {
		defer wg.Done()
		// A more lightweight check for Minio is to check for a bucket's existence.
		// For this example, we'll assume a bucket named "health-check" should exist.
		// In a real application, you might use a more robust check.
		if _, err := r.minioClient.BucketExists(ctx, "health-check"); err != nil {
			statusChan <- domain.HealthStatus{Name: "minio", Status: "down", Error: err.Error()}
		} else {
			statusChan <- domain.HealthStatus{Name: "minio", Status: "up"}
		}
	}()

	wg.Wait()
	close(statusChan)

	statuses := make([]domain.HealthStatus, 0, 3)
	for status := range statusChan {
		statuses = append(statuses, status)
	}

	return statuses
}
