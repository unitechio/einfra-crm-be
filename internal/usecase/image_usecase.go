package usecase

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/unitechio/einfra-be/internal/domain"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

const (
	thumbnailWidth = 200
)

type ImageUsecase interface {
	UploadImage(ctx context.Context, userID string, fileHeader *multipart.FileHeader) (*domain.Image, error)
}

type imageUsecase struct {
	imageRepo   domain.ImageRepository
	storagePath string
}

func NewImageUsecase(imageRepo domain.ImageRepository, storagePath string) ImageUsecase {
	return &imageUsecase{
		imageRepo:   imageRepo,
		storagePath: storagePath,
	}
}

func (uc *imageUsecase) UploadImage(ctx context.Context, userID string, fileHeader *multipart.FileHeader) (*domain.Image, error) {
	image, err := uc.saveOriginalFile(ctx, userID, fileHeader)
	if err != nil {
		return nil, err
	}

	go uc.processImageVariants(image)

	return image, nil
}

func (uc *imageUsecase) saveOriginalFile(ctx context.Context, userID string, fileHeader *multipart.FileHeader) (*domain.Image, error) {
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

func (uc *imageUsecase) processImageVariants(image *domain.Image) {
	ctx := context.Background()

	srcImage, err := imaging.Open(image.Filepath)
	if err != nil {
		log.Printf("ERROR: Failed to open image for processing: %v", err)
		return
	}

	baseFilename := strings.TrimSuffix(image.Filepath, filepath.Ext(image.Filepath))
	image.ThumbnailPath = baseFilename + "_thumb.jpg"
	image.WebpPath = baseFilename + ".webp"

	thumb := imaging.Resize(srcImage, thumbnailWidth, 0, imaging.Lanczos)
	if err = imaging.Save(thumb, image.ThumbnailPath); err != nil {
		log.Printf("ERROR: Failed to create thumbnail: %v", err)
	} else {
		log.Printf("INFO: Created thumbnail for image ID %s at %s", image.ID, image.ThumbnailPath)
	}

	if err = imaging.Save(srcImage, image.WebpPath, imaging.WebPWithQuality(80)); err != nil {
		log.Printf("ERROR: Failed to create WebP image: %v", err)
	}
	log.Printf("INFO: Created WebP for image ID %s at %s", image.ID, image.WebpPath)

	if err := uc.imageRepo.UpdatePaths(ctx, image); err != nil {
		log.Printf("ERROR: Failed to update image paths in DB for image ID %s: %v", image.ID, err)
	}

	log.Printf("INFO: Finished processing variants for image ID %s", image.ID)
}
