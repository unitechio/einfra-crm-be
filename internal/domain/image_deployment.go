package domain

import (
	"context"
	"time"
)

// ImageDeployment represents a deployment of a Harbor image to a Kubernetes cluster
type ImageDeployment struct {
	ID              string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ClusterID       string    `json:"cluster_id" gorm:"type:uuid;not null"`
	Namespace       string    `json:"namespace" gorm:"type:varchar(255);not null"`
	DeploymentName  string    `json:"deployment_name" gorm:"type:varchar(255);not null"`
	ContainerName   string    `json:"container_name" gorm:"type:varchar(255);not null"`
	ImageRepository string    `json:"image_repository" gorm:"type:varchar(500);not null"` // e.g., harbor.example.com/project/image
	ImageTag        string    `json:"image_tag" gorm:"type:varchar(255);not null"`        // e.g., v1.0.0
	DeployedAt      time.Time `json:"deployed_at" gorm:"autoCreateTime"`
	DeployedBy      string    `json:"deployed_by" gorm:"type:varchar(255)"` // User who triggered deployment
	Status          string    `json:"status" gorm:"type:varchar(50)"`       // active, rolled_back, replaced
}

// ImageDeploymentHistory tracks the history of image deployments
type ImageDeploymentHistory struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	DeploymentID string    `json:"deployment_id" gorm:"type:uuid;not null"`
	PreviousTag  string    `json:"previous_tag" gorm:"type:varchar(255)"`
	NewTag       string    `json:"new_tag" gorm:"type:varchar(255);not null"`
	ChangedAt    time.Time `json:"changed_at" gorm:"autoCreateTime"`
	ChangedBy    string    `json:"changed_by" gorm:"type:varchar(255)"`
	Reason       string    `json:"reason" gorm:"type:text"`
}

// ImageDeploymentRepository defines operations for managing image deployments
type ImageDeploymentRepository interface {
	Create(ctx context.Context, deployment *ImageDeployment) error
	GetByID(ctx context.Context, id string) (*ImageDeployment, error)
	List(ctx context.Context, clusterID, namespace string) ([]*ImageDeployment, error)
	GetActiveDeployment(ctx context.Context, clusterID, namespace, deploymentName, containerName string) (*ImageDeployment, error)
	UpdateStatus(ctx context.Context, id, status string) error
}

// ImageDeploymentUsecase defines business logic for image deployments
type ImageDeploymentUsecase interface {
	TrackDeployment(ctx context.Context, clusterID, namespace, deploymentName, containerName, imageRepo, imageTag, user string) error
	GetDeploymentHistory(ctx context.Context, clusterID, namespace, deploymentName string) ([]*ImageDeployment, error)
	GetCurrentDeployments(ctx context.Context, clusterID string) ([]*ImageDeployment, error)
	SyncDeploymentsFromK8s(ctx context.Context, clusterID string) error
}
