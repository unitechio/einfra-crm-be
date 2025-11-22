package repository

import (
	"context"
	"fmt"

	"github.com/unitechio/einfra-be/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ImageRepository struct {
	db *pgxpool.Pool
}

func NewImageRepository(db *pgxpool.Pool) *ImageRepository {
	return &ImageRepository{db: db}
}

func (r *ImageRepository) Create(ctx context.Context, image *domain.Image) error {
	query := `INSERT INTO images (user_id, filename, filepath) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRow(ctx, query, image.UserID, image.Filename, image.Filepath).Scan(&image.ID)
}

func (r *ImageRepository) UpdatePaths(ctx context.Context, image *domain.Image) error {
	query := `UPDATE images SET webp_path = $1, thumbnail_path = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, image.WebpPath, image.ThumbnailPath, image.ID)
	return err
}

func (r *ImageRepository) GetByID(ctx context.Context, id string) (*domain.Image, error) {
	query := `SELECT id, user_id, filename, filepath, webp_path, thumbnail_path FROM images WHERE id = $1`
	image := &domain.Image{}
	err := r.db.QueryRow(ctx, query, id).Scan(&image.ID, &image.UserID, &image.Filename, &image.Filepath, &image.WebpPath, &image.ThumbnailPath)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("image not found")
	}
	return image, err
}

func (r *ImageRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Image, error) {
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
