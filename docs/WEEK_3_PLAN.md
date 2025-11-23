# Week 3 Plan: Advanced Logging & Monitoring 🚀

## 🎯 Objectives
Build a robust monitoring and logging system for the infrastructure.

## 📅 Tasks

### Day 1: Log Streaming System
- [ ] **Log Streamer Package** (`pkg/logstream`)
  - WebSocket based log streaming
  - Support for following logs (tail -f)
  - Filtering by stdout/stderr
  - Timestamp support
- [ ] **Log Usecase & Handler**
  - Integrate with Docker Client
  - WebSocket endpoint: `/api/v1/logs/stream`

### Day 2: Docker Events Monitoring
- [ ] **Event Listener**
  - Listen to Docker daemon events
  - Filter relevant events (die, oom, health_status)
- [ ] **Event Broadcaster**
  - Broadcast events to frontend via WebSocket
  - Store critical events in database (Audit log)

### Day 3: Resource Alerting
- [ ] **Metrics Collector**
  - Background job to collect stats
- [ ] **Alert Rules Engine**
  - Define rules (e.g., CPU > 80% for 5m)
  - Trigger notifications (Email/Slack)

## 🛠️ Technical Stack
- **WebSocket**: `github.com/gorilla/websocket` (already used)
- **Docker Events**: `types.EventsOptions`
- **Concurrency**: Go routines & Channels

---

**Status:** 🚀 Starting Day 1...
