package domain

import (
	"context"
)

// Image represents an uploaded image file with its variants.
type Image struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	Filename      string `json:"filename"`       // Original filename
	Filepath      string `json:"filepath"`       // Path to the original uploaded file
	WebpPath      string `json:"webp_path"`      // Path to the WebP version
	ThumbnailPath string `json:"thumbnail_path"` // Path to the thumbnail version
}

// ImageRepository defines the interface for image data storage.
type ImageRepository interface {
	Create(ctx context.Context, image *Image) error
	UpdatePaths(ctx context.Context, image *Image) error // To save the new paths
	GetByID(ctx context.Context, id string) (*Image, error)
	GetByUserID(ctx context.Context, userID string) ([]*Image, error)
}
