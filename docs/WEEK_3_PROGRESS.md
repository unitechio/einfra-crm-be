# Week 3 Progress: Log Streaming & Events 📡

## ✅ Completed (Day 1 & 2):

### 1. Log Streaming
- **Core Logic**: `pkg/logstream` (Parsing, Timestamps, Stdout/Stderr)
- **Docker Integration**: `pkg/docker/logs.go`
- **API**: `GET /api/v1/logs/containers/:id/stream` (WebSocket)

### 2. Docker Events
- **Core Logic**: `pkg/docker/events.go` (Event stream wrapper)
- **Business Logic**: `internal/usecase/event_usecase.go` (Monitoring, Broadcasting)
- **API**: `GET /api/v1/events/stream` (WebSocket)
- **Handler**: `internal/http/handler/event_handler.go`

### 3. Integration
- ✅ Router updated with both Log and Event routes
- ✅ Main.go wired with all new components
- ✅ Lint errors fixed

---

## 🎯 Next: Resource Alerting (Day 3)

**Goal**: Alert when resources (CPU/RAM) exceed thresholds.

**Tasks:**
1. Create `AlertUsecase`
2. Define Alert Rules (e.g., CPU > 80%)
3. Background Job for checking stats
4. Integrate with Email/Notification system

**Status:** 🟢 Ready for Day 3
