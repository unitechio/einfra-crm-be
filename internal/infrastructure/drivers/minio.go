package drivers

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"mymodule/internal/config"
)

// NewMinioClient creates a new Minio client.
func NewMinioClient(cfg config.MinioConfig) (*minio.Client, error) {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	// Create the bucket if it doesn't exist.
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
		log.Printf("Successfully created bucket: %s\n", cfg.BucketName)
	} else {
		log.Printf("Bucket already exists: %s\n", cfg.BucketName)
	}

	log.Println("Successfully connected to Minio.")

	return minioClient, nil
}
