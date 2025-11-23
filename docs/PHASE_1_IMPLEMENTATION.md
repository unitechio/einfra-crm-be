# Phase 1 Implementation Plan: Core Docker Management

## Overview
Implement core Docker container management features including exec, stats, pause/resume, and stack deployment.

## Timeline: 2-3 weeks

---

## Task 1: Container Operations Enhancement (Week 1)

### 1.1 Container Exec (Interactive Shell)
**Files to create/modify:**
- `internal/usecase/docker_exec_usecase.go`
- `internal/http/handler/docker_handler.go` (update)
- `pkg/docker/exec.go`

**Implementation:**
```go
// pkg/docker/exec.go
type ExecConfig struct {
    ContainerID string
    Cmd         []string
    AttachStdin bool
    AttachStdout bool
    AttachStderr bool
    Tty         bool
    Env         []string
    WorkingDir  string
}

func (c *Client) ContainerExec(ctx context.Context, config ExecConfig) (io.ReadCloser, error)
```

**API Endpoints:**
- `POST /api/v1/docker/containers/:id/exec` - Create exec instance
- `POST /api/v1/docker/exec/:execId/start` - Start exec
- `GET /api/v1/docker/exec/:execId/resize` - Resize TTY

### 1.2 Container Stats (Real-time)
**Files to create/modify:**
- `internal/usecase/docker_stats_usecase.go`
- `pkg/docker/stats.go`
- WebSocket handler for streaming

**Implementation:**
```go
type ContainerStats struct {
    CPU     float64
    Memory  uint64
    MemoryLimit uint64
    NetworkRx uint64
    NetworkTx uint64
    BlockRead uint64
    BlockWrite uint64
}

func (c *Client) ContainerStats(ctx context.Context, containerID string, stream bool) (<-chan ContainerStats, error)
```

**API Endpoints:**
- `GET /api/v1/docker/containers/:id/stats` - Get stats (stream via WebSocket)

### 1.3 Pause/Resume/Unpause
**Files to modify:**
- `internal/usecase/docker_usecase.go`
- `internal/http/handler/docker_handler.go`

**API Endpoints:**
- `POST /api/v1/docker/containers/:id/pause`
- `POST /api/v1/docker/containers/:id/unpause`

### 1.4 Commit Container to Image
**Files to create:**
- `internal/usecase/docker_commit_usecase.go`

**API Endpoints:**
- `POST /api/v1/docker/containers/:id/commit`

---

## Task 2: Image Build & Push (Week 1)

### 2.1 Build Image from Dockerfile
**Files to create:**
- `internal/usecase/docker_build_usecase.go`
- `pkg/docker/build.go`

**Implementation:**
```go
type BuildConfig struct {
    Dockerfile  string
    Context     io.Reader
    Tags        []string
    BuildArgs   map[string]string
    Labels      map[string]string
    NoCache     bool
}

func (c *Client) ImageBuild(ctx context.Context, config BuildConfig) (io.ReadCloser, error)
```

**API Endpoints:**
- `POST /api/v1/docker/images/build` - Build image (upload Dockerfile or tar context)

### 2.2 Push Image to Registry
**Files to modify:**
- `internal/usecase/docker_usecase.go`

**API Endpoints:**
- `POST /api/v1/docker/images/:id/push`

---

## Task 3: Docker Compose/Stack Management (Week 2)

### 3.1 Stack Domain Models
**Files to create:**
- `internal/domain/docker_stack.go`

```go
type DockerStack struct {
    ID          string
    Name        string
    ComposeFile string // YAML content
    EnvVars     map[string]string
    Services    []StackService
    Status      StackStatus
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type StackService struct {
    Name       string
    Image      string
    Replicas   int
    Status     string
}
```

### 3.2 Stack Repository
**Files to create:**
- `internal/repository/docker_stack_repository.go`

### 3.3 Stack Usecase
**Files to create:**
- `internal/usecase/docker_stack_usecase.go`

**Methods:**
- `DeployStack(ctx, name, composeFile, envVars)`
- `UpdateStack(ctx, stackID, composeFile, envVars)`
- `RemoveStack(ctx, stackID)`
- `ListStacks(ctx)`
- `GetStackServices(ctx, stackID)`
- `GetStackLogs(ctx, stackID, serviceFilter)`

### 3.4 Stack Handler
**Files to create:**
- `internal/http/handler/docker_stack_handler.go`

**API Endpoints:**
- `POST /api/v1/docker/stacks` - Deploy stack
- `GET /api/v1/docker/stacks` - List stacks
- `GET /api/v1/docker/stacks/:id` - Get stack details
- `PUT /api/v1/docker/stacks/:id` - Update stack
- `DELETE /api/v1/docker/stacks/:id` - Remove stack
- `GET /api/v1/docker/stacks/:id/services` - Get stack services
- `GET /api/v1/docker/stacks/:id/logs` - Get stack logs

---

## Task 4: Volume File Browser (Week 2)

### 4.1 File Browser Domain
**Files to create:**
- `internal/domain/file_browser.go`

```go
type FileInfo struct {
    Name        string
    Path        string
    Size        int64
    IsDir       bool
    ModTime     time.Time
    Permissions string
}
```

### 4.2 File Browser Usecase
**Files to create:**
- `internal/usecase/file_browser_usecase.go`

**Methods:**
- `ListFiles(ctx, volumeName, path)`
- `UploadFile(ctx, volumeName, path, file)`
- `DownloadFile(ctx, volumeName, path)`
- `DeleteFile(ctx, volumeName, path)`
- `CreateFolder(ctx, volumeName, path)`
- `ReadFile(ctx, volumeName, path)` - For text files

### 4.3 File Browser Handler
**Files to create:**
- `internal/http/handler/file_browser_handler.go`

**API Endpoints:**
- `GET /api/v1/volumes/:name/browse` - List files
- `POST /api/v1/volumes/:name/upload` - Upload file
- `GET /api/v1/volumes/:name/download` - Download file
- `DELETE /api/v1/volumes/:name/files` - Delete file
- `POST /api/v1/volumes/:name/mkdir` - Create folder

---

## Task 5: Network Connect/Disconnect (Week 2)

### 5.1 Network Operations
**Files to modify:**
- `internal/usecase/docker_usecase.go`
- `internal/http/handler/docker_handler.go`

**API Endpoints:**
- `POST /api/v1/docker/networks/:id/connect` - Connect container
- `POST /api/v1/docker/networks/:id/disconnect` - Disconnect container

---

## Task 6: Database Migrations (Week 3)

### 6.1 Create Migration for Stacks
**Files to create:**
- `internal/repository/migrations/021_add_docker_stacks.up.sql`
- `internal/repository/migrations/021_add_docker_stacks.down.sql`

```sql
CREATE TABLE docker_stacks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    compose_file TEXT NOT NULL,
    env_vars JSONB,
    status VARCHAR(50),
    docker_host_id UUID REFERENCES docker_hosts(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stack_services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stack_id UUID REFERENCES docker_stacks(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    image VARCHAR(500),
    replicas INT DEFAULT 1,
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## Task 7: Testing (Week 3)

### 7.1 Unit Tests
**Files to create:**
- `internal/usecase/docker_exec_usecase_test.go`
- `internal/usecase/docker_stats_usecase_test.go`
- `internal/usecase/docker_stack_usecase_test.go`
- `internal/usecase/file_browser_usecase_test.go`

### 7.2 Integration Tests
**Files to create:**
- `tests/integration/docker_operations_test.go`
- `tests/integration/docker_stack_test.go`
- `tests/integration/file_browser_test.go`

---

## Task 8: Documentation (Week 3)

### 8.1 API Documentation
- Update Swagger documentation
- Add usage examples

### 8.2 User Guide
**Files to create:**
- `docs/DOCKER_MANAGEMENT.md`
- `docs/STACK_DEPLOYMENT.md`
- `docs/FILE_BROWSER.md`

---

## Deliverables Checklist

### Week 1:
- [x] Container exec (interactive shell)
- [x] Container stats (real-time)
- [x] Pause/Resume/Unpause
- [x] Commit container to image
- [x] Build image from Dockerfile
- [x] Push image to registry

### Week 2:
- [x] Stack deployment (docker-compose)
- [x] Stack update/remove
- [x] Stack logs viewer
- [x] Volume file browser
- [x] Network connect/disconnect

### Week 3:
- [x] Database migrations
- [x] Unit tests
- [x] Integration tests
- [x] Documentation
- [x] Swagger updates

---

## Dependencies

### Go Packages:
```bash
go get github.com/docker/docker/client
go get github.com/docker/docker/api/types
go get github.com/docker/docker/api/types/container
go get github.com/docker/docker/api/types/network
go get github.com/docker/docker/pkg/archive
go get github.com/gorilla/websocket
```

### Docker Compose:
```bash
# For stack deployment
go get github.com/compose-spec/compose-go
```

---

## Success Criteria

1. ✅ Users can execute commands in running containers
2. ✅ Real-time container stats visible via WebSocket
3. ✅ Containers can be paused/resumed
4. ✅ Containers can be committed to images
5. ✅ Images can be built from Dockerfile via UI
6. ✅ Images can be pushed to registries
7. ✅ Docker Compose stacks can be deployed
8. ✅ Stack logs can be viewed
9. ✅ Files in volumes can be browsed/uploaded/downloaded
10. ✅ Containers can be connected/disconnected from networks

---

## Next Steps

After Phase 1 completion, proceed to:
- **Phase 2**: Monitoring & Security (Dashboard, Vulnerability Scanning, RBAC)
