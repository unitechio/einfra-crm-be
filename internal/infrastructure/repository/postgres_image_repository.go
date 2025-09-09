
package repository

import (
	"context"
	"fmt"
	"mymodule/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresImageRepository is an implementation of ImageRepository for PostgreSQL.
type PostgresImageRepository struct {
	db *pgxpool.Pool
}

// NewPostgresImageRepository creates a new PostgresImageRepository.
func NewPostgresImageRepository(db *pgxpool.Pool) *PostgresImageRepository {
	return &PostgresImageRepository{db: db}
}

// Create inserts a new image record into the database.
func (r *PostgresImageRepository) Create(ctx context.Context, image *domain.Image) error {
	query := `INSERT INTO images (user_id, filename, filepath) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRow(ctx, query, image.UserID, image.Filename, image.Filepath).Scan(&image.ID)
}

// UpdatePaths updates the paths for the generated WebP and thumbnail images.
func (r *PostgresImageRepository) UpdatePaths(ctx context.Context, image *domain.Image) error {
	query := `UPDATE images SET webp_path = $1, thumbnail_path = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, image.WebpPath, image.ThumbnailPath, image.ID)
	return err
}

// GetByID retrieves an image by its ID.
func (r *PostgresImageRepository) GetByID(ctx context.Context, id string) (*domain.Image, error) {
	query := `SELECT id, user_id, filename, filepath, webp_path, thumbnail_path FROM images WHERE id = $1`
	image := &domain.Image{}
	err := r.db.QueryRow(ctx, query, id).Scan(&image.ID, &image.UserID, &image.Filename, &image.Filepath, &image.WebpPath, &image.ThumbnailPath)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("image not found")
	}
	return image, err
}

// GetByUserID retrieves all images uploaded by a specific user.
func (r *PostgresImageRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Image, error) {
	query := `SELECT id, user_id, filename, filepath, webp_path, thumbnail_path FROM images WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*domain.Image
	for rows.Next() {
		var img domain.Image
		if err := rows.Scan(&img.ID, &img.UserID, &img.Filename, &img.Filepath, &img.WebpPath, &img.ThumbnailPath); err != nil {
			return nil, err
		}
		images = append(images, &img)
	}

	return images, nil
}
