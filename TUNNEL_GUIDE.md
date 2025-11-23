SSH Tunnel and Kubeconfig - Complete Implementation
Overview
Fully implemented SSH tunnel support and kubeconfig management for secure connections to private infrastructure (Servers, Docker, Kubernetes).

Files Created
1. Core Infrastructure
pkg/ssh/tunnel.go
 - SSH tunnel with port forwarding
internal/domain/tunnel.go
 - TunnelConfig and KubeConfig models
2. Handlers
internal/http/handler/tunnel_handler.go
 - Tunnel management API
internal/http/handler/kubeconfig_handler.go
 - Kubeconfig management API
3. Usecase Updates
internal/usecase/server_usecase.go
 - Updated with tunnel manager
internal/usecase/server_tunnel_helper.go
 - Tunnel connection helpers
4. Tests
internal/usecase/server_usecase_test.go
 - Unit tests
tests/integration/tunnel_test.go
 - Integration tests
5. Database
internal/repository/migrations/020_add_tunnel_support.up.sql
internal/repository/migrations/020_add_tunnel_support.down.sql
6. Domain Updates
internal/domain/server.go
 - Added tunnel fields
API Endpoints
Tunnel Management
POST   /api/v1/tunnels              - Create tunnel
GET    /api/v1/tunnels              - List active tunnels
GET    /api/v1/tunnels/:id/stats    - Get tunnel stats
DELETE /api/v1/tunnels/:id          - Stop tunnel
POST   /api/v1/tunnels/stop-all     - Stop all tunnels
Kubeconfig Management
POST   /api/v1/kubeconfigs          - Upload kubeconfig
GET    /api/v1/kubeconfigs          - List kubeconfigs
GET    /api/v1/kubeconfigs/:id      - Get kubeconfig
DELETE /api/v1/kubeconfigs/:id      - Delete kubeconfig
POST   /api/v1/kubeconfigs/:id/test - Test connection
Usage Examples
1. Create Tunnel via API
curl -X POST http://localhost:8080/api/v1/tunnels \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "mysql-tunnel",
    "ssh_host": "bastion.example.com",
    "ssh_port": 22,
    "ssh_user": "tunnel-user",
    "ssh_key_path": "/keys/bastion.pem",
    "local_addr": "localhost:3307",
    "remote_addr": "10.0.1.100:3306"
  }'
2. Connect to Private Server
// Server with tunnel enabled
server := &domain.Server{
    Name:          "private-server",
    IPAddress:     "10.0.1.100",
    TunnelEnabled: true,
    TunnelHost:    "bastion.example.com",
    TunnelPort:    22,
    TunnelUser:    "tunnel-user",
    TunnelKeyPath: "/keys/bastion.pem",
}
// Usecase automatically creates tunnel
result, err := serverUsecase.ExecuteCommand(ctx, server.ID, "hostname")
3. Upload Kubeconfig
# Base64 encode kubeconfig
CONFIG_B64=$(cat kubeconfig.yaml | base64)
curl -X POST http://localhost:8080/api/v1/kubeconfigs \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Production Cluster",
    "cluster_id": "cluster-uuid",
    "config_type": "file",
    "config_data": "'$CONFIG_B64'",
    "context_name": "prod-context",
    "is_default": true
  }'
Integration with main.go
Add to dependency injection:

// Initialize tunnel manager
tunnelManager := ssh.NewTunnelManager()
// Update server usecase
serverUsecase := usecase.NewServerUsecase(serverRepo, tunnelManager)
// Initialize handlers
tunnelHandler := handler.NewTunnelHandler(tunnelManager)
kubeconfigHandler := handler.NewKubeconfigHandler()
// Pass to router
router := router.InitRouter(
    cfg,
    // ... existing handlers ...
    tunnelHandler,
    kubeconfigHandler,
)
Testing
Run Unit Tests
go test ./internal/usecase/... -v
Run Integration Tests
go test ./tests/integration/... -v
Skip Integration Tests
go test ./... -short
Database Migration
# Run migration
migrate -path internal/repository/migrations \
        -database "postgres://user:pass@localhost:5432/db?sslmode=disable" \
        up
# Rollback
migrate -path internal/repository/migrations \
        -database "postgres://user:pass@localhost:5432/db?sslmode=disable" \
        down 1
Key Features
✅ SSH Tunnel - Port forwarding through bastion hosts
✅ Tunnel Manager - Multiple concurrent tunnels
✅ Auto Connection - Usecase handles tunnel creation
✅ Kubeconfig Import - Upload and store K8s configs
✅ Real-time Metrics - Collect server metrics via SSH
✅ Connection Caching - Reuse SSH clients
✅ Comprehensive Tests - Unit and integration tests

Security Notes
SSH Keys - Store encrypted, never in plaintext
Bastion Hardening - Limit access, key-only auth
Tunnel Cleanup - Auto-close on server deletion
Kubeconfig Encryption - Encrypt at rest
RBAC - Restrict tunnel/kubeconfig access
Status: ✅ Complete
All components implemented and ready for production use!