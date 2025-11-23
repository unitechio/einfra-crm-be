package domain

import (
	"context"
	"time"
)

// K8sResourceType defines the type of Kubernetes resource
type K8sResourceType string

const (
	K8sResourceConfigMap   K8sResourceType = "ConfigMap"
	K8sResourceSecret      K8sResourceType = "Secret"
	K8sResourceService     K8sResourceType = "Service"
	K8sResourceIngress     K8sResourceType = "Ingress"
	K8sResourceStatefulSet K8sResourceType = "StatefulSet"
	K8sResourceDaemonSet   K8sResourceType = "DaemonSet"
	K8sResourceJob         K8sResourceType = "Job"
	K8sResourceCronJob     K8sResourceType = "CronJob"
	K8sResourcePV          K8sResourceType = "PersistentVolume"
	K8sResourcePVC         K8sResourceType = "PersistentVolumeClaim"
	K8sResourceDeployment  K8sResourceType = "Deployment"
	K8sResourcePod         K8sResourceType = "Pod"
)

// K8sResource represents a generic Kubernetes resource
type K8sResource struct {
	ID        string          `json:"id"`
	ClusterID string          `json:"cluster_id"`
	Namespace string          `json:"namespace"`
	Name      string          `json:"name"`
	Kind      K8sResourceType `json:"kind"`
	Manifest  string          `json:"manifest"` // YAML content
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// K8sBackup represents a backup of Kubernetes resources
type K8sBackup struct {
	ID            string              `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ClusterID     string              `json:"cluster_id" gorm:"type:uuid;not null"`
	Name          string              `json:"name" gorm:"type:varchar(255);not null"`
	Description   string              `json:"description" gorm:"type:text"`
	Namespace     string              `json:"namespace" gorm:"type:varchar(255)"` // Optional: if backing up specific namespace
	ResourceCount int                 `json:"resource_count" gorm:"type:int"`
	SizeBytes     int64               `json:"size_bytes" gorm:"type:bigint"`
	Status        string              `json:"status" gorm:"type:varchar(50)"` // pending, completed, failed
	CreatedAt     time.Time           `json:"created_at" gorm:"autoCreateTime"`
	Resources     []K8sBackupResource `json:"resources" gorm:"foreignKey:BackupID"`
}

// K8sBackupResource represents a single resource within a backup
type K8sBackupResource struct {
	ID        string          `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	BackupID  string          `json:"backup_id" gorm:"type:uuid;not null"`
	Kind      K8sResourceType `json:"kind" gorm:"type:varchar(100);not null"`
	Namespace string          `json:"namespace" gorm:"type:varchar(255)"`
	Name      string          `json:"name" gorm:"type:varchar(255);not null"`
	Manifest  string          `json:"manifest" gorm:"type:text;not null"` // YAML content
	CreatedAt time.Time       `json:"created_at" gorm:"autoCreateTime"`
}

// K8sBackupRepository defines operations for managing backups
type K8sBackupRepository interface {
	Create(ctx context.Context, backup *K8sBackup) error
	GetByID(ctx context.Context, id string) (*K8sBackup, error)
	List(ctx context.Context, clusterID, namespace string) ([]*K8sBackup, error)
	Delete(ctx context.Context, id string) error
}

// K8sBackupUsecase defines business logic for backups
type K8sBackupUsecase interface {
	BackupNamespace(ctx context.Context, clusterID, namespace, name, description, user string) (*K8sBackup, error)
	RestoreBackup(ctx context.Context, backupID, user string) error
	ListBackups(ctx context.Context, clusterID, namespace string) ([]*K8sBackup, error)
	GetBackup(ctx context.Context, id string) (*K8sBackup, error)
	DeleteBackup(ctx context.Context, id string) error
}
