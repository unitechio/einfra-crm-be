
package usecase

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"mymodule/internal/domain"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

const (
	thumbnailWidth = 200
)

// ImageUsecase handles the business logic for image operations.
type ImageUsecase struct {
	imageRepo   domain.ImageRepository
	storagePath string // The base path where images are stored
}

// NewImageUsecase creates a new ImageUsecase.
func NewImageUsecase(imageRepo domain.ImageRepository, storagePath string) *ImageUsecase {
	return &ImageUsecase{
		imageRepo:   imageRepo,
		storagePath: storagePath,
	}
}

// UploadImage handles saving, processing, and creating a database record for an image.
func (uc *ImageUsecase) UploadImage(ctx context.Context, userID string, fileHeader *multipart.FileHeader) (*domain.Image, error) {
	// 1. Save the original file
	image, err := uc.saveOriginalFile(ctx, userID, fileHeader)
	if err != nil {
		return nil, err
	}

	// 2. Process the image to create variants (thumbnail, webp)
	// In a production app, this could be offloaded to a background worker.
	go uc.processImageVariants(image)

	// 3. Return the initial image object. The paths for variants will be updated in the background.
	return image, nil
}

// saveOriginalFile saves the uploaded file and creates the initial DB record.
func (uc *ImageUsecase) saveOriginalFile(ctx context.Context, userID string, fileHeader *multipart.FileHeader) (*domain.Image, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	if err := os.MkdirAll(uc.storagePath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	ext := filepath.Ext(fileHeader.Filename)
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	dstPath := filepath.Join(uc.storagePath, newFilename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	image := &domain.Image{
		UserID:   userID,
		Filename: fileHeader.Filename,
		Filepath: dstPath,
	}

	if err := uc.imageRepo.Create(ctx, image); err != nil {
		os.Remove(dstPath)
		return nil, fmt.Errorf("failed to create image record in db: %w", err)
	}

	return image, nil
}

// processImageVariants creates thumbnail and WebP versions of the image.
// It runs in a goroutine and updates the DB record asynchronously.
func (uc *ImageUsecase) processImageVariants(image *domain.Image) {
	// We use a background context because the original request might have already finished.
	ctx := context.Background()

	// Open the saved original image
	srcImage, err := imaging.Open(image.Filepath)
	if err != nil {
		log.Printf("ERROR: Failed to open image for processing: %v", err)
		return
	}

	// Generate new file paths
	baseFilename := strings.TrimSuffix(image.Filepath, filepath.Ext(image.Filepath))
	image.ThumbnailPath = baseFilename + "_thumb.jpg"
	image.WebpPath = baseFilename + ".webp"

	// Create Thumbnail
	thumb := imaging.Resize(srcImage, thumbnailWidth, 0, imaging.Lanczos)
	if err = imaging.Save(thumb, image.ThumbnailPath); err != nil {
		log.Printf("ERROR: Failed to create thumbnail: %v", err)
		// Continue to try creating WebP version
	} else {
		log.Printf("INFO: Created thumbnail for image ID %s at %s", image.ID, image.ThumbnailPath)
	}

	// Create WebP
	if err = imaging.Save(srcImage, image.WebpPath, imaging.WebPWithQuality(80)); err != nil {
		log.Printf("ERROR: Failed to create WebP image: %v", err)
	}
	log.Printf("INFO: Created WebP for image ID %s at %s", image.ID, image.WebpPath)

	// Update the database with the new paths
	if err := uc.imageRepo.UpdatePaths(ctx, image); err != nil {
		log.Printf("ERROR: Failed to update image paths in DB for image ID %s: %v", image.ID, err)
	}

	log.Printf("INFO: Finished processing variants for image ID %s", image.ID)
}
