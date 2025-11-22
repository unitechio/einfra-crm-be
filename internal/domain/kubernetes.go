package domain

import (
	"context"
	"time"
)

// K8sCluster represents a Kubernetes cluster
// @Description Kubernetes cluster configuration and connection details
type K8sCluster struct {
	ID          string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null;uniqueIndex" validate:"required,min=3,max=255" example:"prod-k8s-cluster"`
	Description string     `json:"description" gorm:"type:text" example:"Production Kubernetes cluster"`
	APIServer   string     `json:"api_server" gorm:"type:varchar(500);not null" validate:"required,url" example:"https://k8s.example.com:6443"`
	Version     string     `json:"version" gorm:"type:varchar(50)" example:"v1.28.0"`
	Provider    string     `json:"provider" gorm:"type:varchar(100)" example:"EKS"` // EKS, GKE, AKS, self-hosted
	Region      string     `json:"region" gorm:"type:varchar(100)" example:"us-east-1"`
	ConfigPath  string     `json:"config_path" gorm:"type:varchar(500)" example:"/etc/k8s/config"`
	IsActive    bool       `json:"is_active" gorm:"type:boolean;default:true;index" example:"true"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime" example:"2024-01-01T00:00:00Z"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string" example:"2024-01-01T00:00:00Z"`
}

// TableName specifies the table name for K8sCluster model
func (K8sCluster) TableName() string {
	return "k8s_clusters"
}

// K8sNamespace represents a Kubernetes namespace
// @Description Kubernetes namespace information
type K8sNamespace struct {
	Name        string            `json:"name" example:"production"`
	ClusterID   string            `json:"cluster_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Labels      map[string]string `json:"labels" example:"env:prod,team:platform"`
	Annotations map[string]string `json:"annotations"`
	Status      string            `json:"status" example:"Active"`
	CreatedAt   time.Time         `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

// K8sDeployment represents a Kubernetes deployment
// @Description Kubernetes deployment configuration and status
type K8sDeployment struct {
	Name              string            `json:"name" example:"nginx-deployment"`
	Namespace         string            `json:"namespace" example:"production"`
	ClusterID         string            `json:"cluster_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Labels            map[string]string `json:"labels" example:"app:nginx,version:1.0"`
	Replicas          int32             `json:"replicas" example:"3"`
	AvailableReplicas int32             `json:"available_replicas" example:"3"`
	ReadyReplicas     int32             `json:"ready_replicas" example:"3"`
	UpdatedReplicas   int32             `json:"updated_replicas" example:"3"`
	Image             string            `json:"image" example:"nginx:1.25"`
	Strategy          string            `json:"strategy" example:"RollingUpdate"`
	Conditions        []K8sCondition    `json:"conditions"`
	CreatedAt         time.Time         `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

// K8sService represents a Kubernetes service
// @Description Kubernetes service configuration
type K8sService struct {
	Name       string            `json:"name" example:"nginx-service"`
	Namespace  string            `json:"namespace" example:"production"`
	ClusterID  string            `json:"cluster_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Type       string            `json:"type" example:"LoadBalancer"` // ClusterIP, NodePort, LoadBalancer, ExternalName
	ClusterIP  string            `json:"cluster_ip" example:"10.96.0.1"`
	ExternalIP string            `json:"external_ip,omitempty" example:"203.0.113.1"`
	Ports      []K8sServicePort  `json:"ports"`
	Selector   map[string]string `json:"selector" example:"app:nginx"`
	Labels     map[string]string `json:"labels"`
	CreatedAt  time.Time         `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

// K8sServicePort represents a service port configuration
// @Description Kubernetes service port
type K8sServicePort struct {
	Name       string `json:"name" example:"http"`
	Protocol   string `json:"protocol" example:"TCP"`
	Port       int32  `json:"port" example:"80"`
	TargetPort string `json:"target_port" example:"8080"`
	NodePort   int32  `json:"node_port,omitempty" example:"30080"`
}

// K8sPod represents a Kubernetes pod
// @Description Kubernetes pod information and status
type K8sPod struct {
	Name         string            `json:"name" example:"nginx-deployment-7d6b8c5f9-abc12"`
	Namespace    string            `json:"namespace" example:"production"`
	ClusterID    string            `json:"cluster_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Labels       map[string]string `json:"labels" example:"app:nginx,pod-template-hash:7d6b8c5f9"`
	Phase        string            `json:"phase" example:"Running"` // Pending, Running, Succeeded, Failed, Unknown
	PodIP        string            `json:"pod_ip" example:"10.244.0.5"`
	HostIP       string            `json:"host_ip" example:"192.168.1.10"`
	NodeName     string            `json:"node_name" example:"node-1"`
	Containers   []K8sContainer    `json:"containers"`
	Conditions   []K8sCondition    `json:"conditions"`
	RestartCount int32             `json:"restart_count" example:"0"`
	CreatedAt    time.Time         `json:"created_at" example:"2024-01-01T00:00:00Z"`
	StartedAt    *time.Time        `json:"started_at,omitempty" example:"2024-01-01T00:00:00Z"`
}

// K8sContainer represents a container within a pod
// @Description Kubernetes container information
type K8sContainer struct {
	Name         string     `json:"name" example:"nginx"`
	Image        string     `json:"image" example:"nginx:1.25"`
	Ready        bool       `json:"ready" example:"true"`
	RestartCount int32      `json:"restart_count" example:"0"`
	State        string     `json:"state" example:"running"` // waiting, running, terminated
	StartedAt    *time.Time `json:"started_at,omitempty" example:"2024-01-01T00:00:00Z"`
}

// K8sCondition represents a condition of a Kubernetes resource
// @Description Kubernetes resource condition
type K8sCondition struct {
	Type               string    `json:"type" example:"Ready"`
	Status             string    `json:"status" example:"True"` // True, False, Unknown
	LastTransitionTime time.Time `json:"last_transition_time" example:"2024-01-01T00:00:00Z"`
	Reason             string    `json:"reason,omitempty" example:"MinimumReplicasAvailable"`
	Message            string    `json:"message,omitempty" example:"Deployment has minimum availability"`
}

// K8sNode represents a Kubernetes node
// @Description Kubernetes node information
type K8sNode struct {
	Name             string            `json:"name" example:"node-1"`
	ClusterID        string            `json:"cluster_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Labels           map[string]string `json:"labels"`
	Status           string            `json:"status" example:"Ready"`
	Roles            []string          `json:"roles" example:"master,control-plane"`
	KubeletVersion   string            `json:"kubelet_version" example:"v1.28.0"`
	OSImage          string            `json:"os_image" example:"Ubuntu 22.04.3 LTS"`
	KernelVersion    string            `json:"kernel_version" example:"5.15.0-91-generic"`
	ContainerRuntime string            `json:"container_runtime" example:"containerd://1.7.2"`
	CPUCapacity      string            `json:"cpu_capacity" example:"8"`
	MemoryCapacity   string            `json:"memory_capacity" example:"32Gi"`
	PodCapacity      string            `json:"pod_capacity" example:"110"`
	Conditions       []K8sCondition    `json:"conditions"`
	CreatedAt        time.Time         `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

// K8sClusterRepository defines the interface for Kubernetes cluster data persistence
type K8sClusterRepository interface {
	// Create creates a new Kubernetes cluster record
	Create(ctx context.Context, cluster *K8sCluster) error

	// GetByID retrieves a cluster by ID
	GetByID(ctx context.Context, id string) (*K8sCluster, error)

	// List retrieves all clusters with filtering
	List(ctx context.Context, filter K8sClusterFilter) ([]*K8sCluster, int64, error)

	// Update updates a cluster
	Update(ctx context.Context, cluster *K8sCluster) error

	// Delete soft deletes a cluster
	Delete(ctx context.Context, id string) error
}

// K8sClusterFilter represents filtering options for cluster queries
type K8sClusterFilter struct {
	Provider string `json:"provider,omitempty"`
	Region   string `json:"region,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

// KubernetesUsecase defines the business logic for Kubernetes management
type KubernetesUsecase interface {
	// Cluster Management
	CreateCluster(ctx context.Context, cluster *K8sCluster) error
	GetCluster(ctx context.Context, id string) (*K8sCluster, error)
	ListClusters(ctx context.Context, filter K8sClusterFilter) ([]*K8sCluster, int64, error)
	UpdateCluster(ctx context.Context, cluster *K8sCluster) error
	DeleteCluster(ctx context.Context, id string) error
	GetClusterInfo(ctx context.Context, clusterID string) (map[string]interface{}, error)

	// Namespace Management
	ListNamespaces(ctx context.Context, clusterID string) ([]*K8sNamespace, error)
	CreateNamespace(ctx context.Context, clusterID, name string, labels map[string]string) error
	DeleteNamespace(ctx context.Context, clusterID, name string) error

	// Deployment Management
	ListDeployments(ctx context.Context, clusterID, namespace string) ([]*K8sDeployment, error)
	GetDeployment(ctx context.Context, clusterID, namespace, name string) (*K8sDeployment, error)
	CreateDeployment(ctx context.Context, clusterID string, deployment interface{}) error
	UpdateDeployment(ctx context.Context, clusterID string, deployment interface{}) error
	DeleteDeployment(ctx context.Context, clusterID, namespace, name string) error
	ScaleDeployment(ctx context.Context, clusterID, namespace, name string, replicas int32) error

	// Service Management
	ListServices(ctx context.Context, clusterID, namespace string) ([]*K8sService, error)
	GetService(ctx context.Context, clusterID, namespace, name string) (*K8sService, error)
	CreateService(ctx context.Context, clusterID string, service interface{}) error
	DeleteService(ctx context.Context, clusterID, namespace, name string) error

	// Pod Management
	ListPods(ctx context.Context, clusterID, namespace string) ([]*K8sPod, error)
	GetPod(ctx context.Context, clusterID, namespace, name string) (*K8sPod, error)
	DeletePod(ctx context.Context, clusterID, namespace, name string) error
	GetPodLogs(ctx context.Context, clusterID, namespace, podName, containerName string, tail int) (string, error)
	ExecPodCommand(ctx context.Context, clusterID, namespace, podName, containerName string, command []string) (string, error)

	// Node Management
	ListNodes(ctx context.Context, clusterID string) ([]*K8sNode, error)
	GetNode(ctx context.Context, clusterID, name string) (*K8sNode, error)
}
