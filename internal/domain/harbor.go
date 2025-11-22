package domain

import (
	"context"
	"time"
)

// HarborRegistry represents a Harbor container registry
// @Description Harbor registry configuration and connection details
type HarborRegistry struct {
	ID          string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null;uniqueIndex" validate:"required,min=3,max=255" example:"prod-harbor"`
	Description string     `json:"description" gorm:"type:text" example:"Production Harbor registry"`
	URL         string     `json:"url" gorm:"type:varchar(500);not null" validate:"required,url" example:"https://harbor.example.com"`
	Username    string     `json:"username" gorm:"type:varchar(255);not null" validate:"required" example:"admin"`
	Password    string     `json:"password,omitempty" gorm:"type:varchar(500)" validate:"required" example:"Harbor12345"` // Should be encrypted
	Version     string     `json:"version" gorm:"type:varchar(50)" example:"v2.9.0"`
	IsDefault   bool       `json:"is_default" gorm:"type:boolean;default:false" example:"false"`
	IsActive    bool       `json:"is_active" gorm:"type:boolean;default:true;index" example:"true"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime" example:"2024-01-01T00:00:00Z"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string" example:"2024-01-01T00:00:00Z"`
}

// TableName specifies the table name for HarborRegistry model
func (HarborRegistry) TableName() string {
	return "harbor_registries"
}

// HarborProject represents a Harbor project
// @Description Harbor project with metadata and quotas
type HarborProject struct {
	ID           int64             `json:"id" example:"1"`
	Name         string            `json:"name" example:"library"`
	RegistryID   string            `json:"registry_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Public       bool              `json:"public" example:"false"`
	OwnerName    string            `json:"owner_name" example:"admin"`
	RepoCount    int64             `json:"repo_count" example:"10"`
	ChartCount   int64             `json:"chart_count" example:"5"`
	Metadata     map[string]string `json:"metadata"`
	CVEAllowlist []string          `json:"cve_allowlist"`
	StorageLimit int64             `json:"storage_limit" example:"10737418240"` // bytes
	CreatedAt    time.Time         `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt    time.Time         `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// HarborRepository represents a repository within a Harbor project
// @Description Harbor repository information
type HarborRepository struct {
	ID            int64     `json:"id" example:"1"`
	Name          string    `json:"name" example:"library/nginx"`
	ProjectID     int64     `json:"project_id" example:"1"`
	RegistryID    string    `json:"registry_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Description   string    `json:"description" example:"Nginx web server images"`
	ArtifactCount int64     `json:"artifact_count" example:"15"`
	PullCount     int64     `json:"pull_count" example:"1000"`
	CreatedAt     time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt     time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// HarborArtifact represents a container image artifact
// @Description Harbor artifact (container image) with tags and scan results
type HarborArtifact struct {
	ID                int64               `json:"id" example:"1"`
	Digest            string              `json:"digest" example:"sha256:abcd1234..."`
	RepositoryName    string              `json:"repository_name" example:"library/nginx"`
	ProjectID         int64               `json:"project_id" example:"1"`
	RegistryID        string              `json:"registry_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Tags              []HarborTag         `json:"tags"`
	Type              string              `json:"type" example:"IMAGE"`
	Size              int64               `json:"size" example:"142000000"`
	PushTime          time.Time           `json:"push_time" example:"2024-01-01T00:00:00Z"`
	PullTime          *time.Time          `json:"pull_time,omitempty" example:"2024-01-01T00:00:00Z"`
	ScanOverview      *HarborScanOverview `json:"scan_overview,omitempty"`
	Labels            []HarborLabel       `json:"labels"`
	Annotations       map[string]string   `json:"annotations"`
	ManifestMediaType string              `json:"manifest_media_type" example:"application/vnd.docker.distribution.manifest.v2+json"`
}

// HarborTag represents an image tag
// @Description Harbor artifact tag
type HarborTag struct {
	ID        int64      `json:"id" example:"1"`
	Name      string     `json:"name" example:"latest"`
	PushTime  time.Time  `json:"push_time" example:"2024-01-01T00:00:00Z"`
	PullTime  *time.Time `json:"pull_time,omitempty" example:"2024-01-01T00:00:00Z"`
	Immutable bool       `json:"immutable" example:"false"`
	Signed    bool       `json:"signed" example:"false"`
}

// HarborScanOverview represents vulnerability scan results
// @Description Harbor vulnerability scan overview
type HarborScanOverview struct {
	ScanStatus string                     `json:"scan_status" example:"Success"` // Pending, Running, Success, Error
	Severity   string                     `json:"severity" example:"High"`       // None, Unknown, Low, Medium, High, Critical
	Duration   int64                      `json:"duration" example:"45"`         // seconds
	StartTime  time.Time                  `json:"start_time" example:"2024-01-01T00:00:00Z"`
	EndTime    time.Time                  `json:"end_time" example:"2024-01-01T00:00:00Z"`
	Summary    HarborVulnerabilitySummary `json:"summary"`
}

// HarborVulnerabilitySummary represents a summary of vulnerabilities
// @Description Vulnerability count by severity
type HarborVulnerabilitySummary struct {
	Total    int `json:"total" example:"25"`
	Critical int `json:"critical" example:"2"`
	High     int `json:"high" example:"5"`
	Medium   int `json:"medium" example:"10"`
	Low      int `json:"low" example:"8"`
	Unknown  int `json:"unknown" example:"0"`
}

// HarborVulnerability represents a detailed vulnerability
// @Description Detailed vulnerability information
type HarborVulnerability struct {
	ID          string   `json:"id" example:"CVE-2024-1234"`
	Package     string   `json:"package" example:"openssl"`
	Version     string   `json:"version" example:"1.1.1k"`
	FixVersion  string   `json:"fix_version" example:"1.1.1w"`
	Severity    string   `json:"severity" example:"High"`
	Description string   `json:"description" example:"Buffer overflow vulnerability"`
	Links       []string `json:"links" example:"https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2024-1234"`
	CVSSScore   float64  `json:"cvss_score" example:"7.5"`
}

// HarborLabel represents a label for artifacts
// @Description Harbor label
type HarborLabel struct {
	ID          int64     `json:"id" example:"1"`
	Name        string    `json:"name" example:"production"`
	Description string    `json:"description" example:"Production environment"`
	Color       string    `json:"color" example:"#0099CC"`
	Scope       string    `json:"scope" example:"global"` // global or project
	ProjectID   int64     `json:"project_id,omitempty" example:"1"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// HarborQuota represents storage quota information
// @Description Harbor project storage quota
type HarborQuota struct {
	ID        int64                  `json:"id" example:"1"`
	Ref       map[string]interface{} `json:"ref"`
	Hard      map[string]int64       `json:"hard" example:"storage:10737418240"`
	Used      map[string]int64       `json:"used" example:"storage:5368709120"`
	CreatedAt time.Time              `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time              `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// HarborRegistryRepository defines the interface for Harbor registry data persistence
type HarborRegistryRepository interface {
	// Create creates a new Harbor registry record
	Create(ctx context.Context, registry *HarborRegistry) error

	// GetByID retrieves a registry by ID
	GetByID(ctx context.Context, id string) (*HarborRegistry, error)

	// List retrieves all registries with filtering
	List(ctx context.Context, filter HarborRegistryFilter) ([]*HarborRegistry, int64, error)

	// Update updates a registry
	Update(ctx context.Context, registry *HarborRegistry) error

	// Delete soft deletes a registry
	Delete(ctx context.Context, id string) error

	// GetDefault retrieves the default registry
	GetDefault(ctx context.Context) (*HarborRegistry, error)
}

// HarborRegistryFilter represents filtering options for registry queries
type HarborRegistryFilter struct {
	IsActive  *bool `json:"is_active,omitempty"`
	IsDefault *bool `json:"is_default,omitempty"`
	Page      int   `json:"page" validate:"min=1"`
	PageSize  int   `json:"page_size" validate:"min=1,max=100"`
}

// HarborUsecase defines the business logic for Harbor management
type HarborUsecase interface {
	// Registry Management
	CreateRegistry(ctx context.Context, registry *HarborRegistry) error
	GetRegistry(ctx context.Context, id string) (*HarborRegistry, error)
	ListRegistries(ctx context.Context, filter HarborRegistryFilter) ([]*HarborRegistry, int64, error)
	UpdateRegistry(ctx context.Context, registry *HarborRegistry) error
	DeleteRegistry(ctx context.Context, id string) error
	TestRegistryConnection(ctx context.Context, id string) (bool, error)

	// Project Management
	ListProjects(ctx context.Context, registryID string, public *bool) ([]*HarborProject, error)
	GetProject(ctx context.Context, registryID string, projectID int64) (*HarborProject, error)
	CreateProject(ctx context.Context, registryID, name string, public bool, storageLimit int64) (*HarborProject, error)
	UpdateProject(ctx context.Context, registryID string, projectID int64, updates map[string]interface{}) error
	DeleteProject(ctx context.Context, registryID string, projectID int64) error

	// Repository Management
	ListRepositories(ctx context.Context, registryID string, projectName string) ([]*HarborRepository, error)
	GetRepository(ctx context.Context, registryID, repositoryName string) (*HarborRepository, error)
	DeleteRepository(ctx context.Context, registryID, repositoryName string) error

	// Artifact Management
	ListArtifacts(ctx context.Context, registryID, repositoryName string) ([]*HarborArtifact, error)
	GetArtifact(ctx context.Context, registryID, repositoryName, reference string) (*HarborArtifact, error)
	DeleteArtifact(ctx context.Context, registryID, repositoryName, reference string) error
	CopyArtifact(ctx context.Context, registryID, srcRepo, dstRepo, reference string) error

	// Vulnerability Scanning
	ScanArtifact(ctx context.Context, registryID, repositoryName, reference string) error
	GetScanReport(ctx context.Context, registryID, repositoryName, reference string) (*HarborScanOverview, error)
	GetVulnerabilities(ctx context.Context, registryID, repositoryName, reference string) ([]*HarborVulnerability, error)

	// Label Management
	ListLabels(ctx context.Context, registryID string, scope string, projectID *int64) ([]*HarborLabel, error)
	CreateLabel(ctx context.Context, registryID string, label *HarborLabel) error
	DeleteLabel(ctx context.Context, registryID string, labelID int64) error
	AddLabelToArtifact(ctx context.Context, registryID, repositoryName, reference string, labelID int64) error
	RemoveLabelFromArtifact(ctx context.Context, registryID, repositoryName, reference string, labelID int64) error

	// Quota Management
	GetProjectQuota(ctx context.Context, registryID string, projectID int64) (*HarborQuota, error)
	UpdateProjectQuota(ctx context.Context, registryID string, projectID int64, storageLimit int64) error
}
