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

// K8sConfigMap represents a Kubernetes ConfigMap
// @Description Kubernetes ConfigMap
type K8sConfigMap struct {
	Name      string            `json:"name" example:"game-config"`
	Namespace string            `json:"namespace" example:"default"`
	ClusterID string            `json:"cluster_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Data      map[string]string `json:"data"`
	CreatedAt time.Time         `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

// K8sSecret represents a Kubernetes Secret
// @Description Kubernetes Secret
type K8sSecret struct {
	Name      string            `json:"name" example:"db-secret"`
	Namespace string            `json:"namespace" example:"default"`
	ClusterID string            `json:"cluster_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Type      string            `json:"type" example:"Opaque"`
	Data      map[string]string `json:"data,omitempty"` // Base64 encoded values
	CreatedAt time.Time         `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

// K8sIngress represents a Kubernetes Ingress
// @Description Kubernetes Ingress
type K8sIngress struct {
	Name      string           `json:"name" example:"web-ingress"`
	Namespace string           `json:"namespace" example:"default"`
	ClusterID string           `json:"cluster_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Rules     []K8sIngressRule `json:"rules"`
	CreatedAt time.Time        `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

type K8sIngressRule struct {
	Host  string           `json:"host" example:"example.com"`
	Paths []K8sIngressPath `json:"paths"`
}

type K8sIngressPath struct {
	Path        string `json:"path" example:"/api"`
	PathType    string `json:"path_type" example:"Prefix"`
	ServiceName string `json:"service_name" example:"api-service"`
	ServicePort int32  `json:"service_port" example:"80"`
}

// K8sStatefulSet represents a Kubernetes StatefulSet
// @Description Kubernetes StatefulSet
type K8sStatefulSet struct {
	Name          string    `json:"name" example:"web"`
	Namespace     string    `json:"namespace" example:"default"`
	ClusterID     string    `json:"cluster_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Replicas      int32     `json:"replicas" example:"3"`
	ReadyReplicas int32     `json:"ready_replicas" example:"3"`
	ServiceName   string    `json:"service_name" example:"nginx"`
	CreatedAt     time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

// K8sDaemonSet represents a Kubernetes DaemonSet
// @Description Kubernetes DaemonSet
type K8sDaemonSet struct {
	Name                   string    `json:"name" example:"fluentd"`
	Namespace              string    `json:"namespace" example:"kube-system"`
	ClusterID              string    `json:"cluster_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	DesiredNumberScheduled int32     `json:"desired_number_scheduled" example:"5"`
	NumberReady            int32     `json:"number_ready" example:"5"`
	CreatedAt              time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

// K8sJob represents a Kubernetes Job
// @Description Kubernetes Job
type K8sJob struct {
	Name        string    `json:"name" example:"pi"`
	Namespace   string    `json:"namespace" example:"default"`
	ClusterID   string    `json:"cluster_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Completions int32     `json:"completions" example:"1"`
	Succeeded   int32     `json:"succeeded" example:"1"`
	Failed      int32     `json:"failed" example:"0"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

// K8sCronJob represents a Kubernetes CronJob
// @Description Kubernetes CronJob
type K8sCronJob struct {
	Name             string     `json:"name" example:"hello"`
	Namespace        string     `json:"namespace" example:"default"`
	ClusterID        string     `json:"cluster_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Schedule         string     `json:"schedule" example:"*/1 * * * *"`
	Suspend          bool       `json:"suspend" example:"false"`
	LastScheduleTime *time.Time `json:"last_schedule_time,omitempty"`
	CreatedAt        time.Time  `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

// K8sPV represents a Kubernetes PersistentVolume
// @Description Kubernetes PersistentVolume
type K8sPV struct {
	Name          string    `json:"name" example:"pv0003"`
	ClusterID     string    `json:"cluster_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Capacity      string    `json:"capacity" example:"5Gi"`
	AccessModes   []string  `json:"access_modes" example:"ReadWriteOnce"`
	ReclaimPolicy string    `json:"reclaim_policy" example:"Recycle"`
	Status        string    `json:"status" example:"Available"`
	ClaimRef      string    `json:"claim_ref,omitempty" example:"default/myclaim"`
	CreatedAt     time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

// K8sPVC represents a Kubernetes PersistentVolumeClaim
// @Description Kubernetes PersistentVolumeClaim
type K8sPVC struct {
	Name        string    `json:"name" example:"myclaim"`
	Namespace   string    `json:"namespace" example:"default"`
	ClusterID   string    `json:"cluster_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Status      string    `json:"status" example:"Bound"`
	Volume      string    `json:"volume" example:"pv0003"`
	Capacity    string    `json:"capacity" example:"5Gi"`
	AccessModes []string  `json:"access_modes" example:"ReadWriteOnce"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
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

	// ConfigMap Management
	ListConfigMaps(ctx context.Context, clusterID, namespace string) ([]*K8sConfigMap, error)
	GetConfigMap(ctx context.Context, clusterID, namespace, name string) (*K8sConfigMap, error)
	CreateConfigMap(ctx context.Context, clusterID string, configMap interface{}) error
	DeleteConfigMap(ctx context.Context, clusterID, namespace, name string) error

	// Secret Management
	ListSecrets(ctx context.Context, clusterID, namespace string) ([]*K8sSecret, error)
	GetSecret(ctx context.Context, clusterID, namespace, name string) (*K8sSecret, error)
	CreateSecret(ctx context.Context, clusterID string, secret interface{}) error
	DeleteSecret(ctx context.Context, clusterID, namespace, name string) error

	// Ingress Management
	ListIngresses(ctx context.Context, clusterID, namespace string) ([]*K8sIngress, error)
	GetIngress(ctx context.Context, clusterID, namespace, name string) (*K8sIngress, error)
	CreateIngress(ctx context.Context, clusterID string, ingress interface{}) error
	DeleteIngress(ctx context.Context, clusterID, namespace, name string) error

	// StatefulSet Management
	ListStatefulSets(ctx context.Context, clusterID, namespace string) ([]*K8sStatefulSet, error)
	GetStatefulSet(ctx context.Context, clusterID, namespace, name string) (*K8sStatefulSet, error)
	CreateStatefulSet(ctx context.Context, clusterID string, statefulSet interface{}) error
	DeleteStatefulSet(ctx context.Context, clusterID, namespace, name string) error

	// DaemonSet Management
	ListDaemonSets(ctx context.Context, clusterID, namespace string) ([]*K8sDaemonSet, error)
	GetDaemonSet(ctx context.Context, clusterID, namespace, name string) (*K8sDaemonSet, error)
	CreateDaemonSet(ctx context.Context, clusterID string, daemonSet interface{}) error
	DeleteDaemonSet(ctx context.Context, clusterID, namespace, name string) error

	// Job Management
	ListJobs(ctx context.Context, clusterID, namespace string) ([]*K8sJob, error)
	GetJob(ctx context.Context, clusterID, namespace, name string) (*K8sJob, error)
	CreateJob(ctx context.Context, clusterID string, job interface{}) error
	DeleteJob(ctx context.Context, clusterID, namespace, name string) error

	// CronJob Management
	ListCronJobs(ctx context.Context, clusterID, namespace string) ([]*K8sCronJob, error)
	GetCronJob(ctx context.Context, clusterID, namespace, name string) (*K8sCronJob, error)
	CreateCronJob(ctx context.Context, clusterID string, cronJob interface{}) error
	DeleteCronJob(ctx context.Context, clusterID, namespace, name string) error

	// PV/PVC Management
	ListPVs(ctx context.Context, clusterID string) ([]*K8sPV, error)
	GetPV(ctx context.Context, clusterID, name string) (*K8sPV, error)
	ListPVCs(ctx context.Context, clusterID, namespace string) ([]*K8sPVC, error)
	GetPVC(ctx context.Context, clusterID, namespace, name string) (*K8sPVC, error)
}
