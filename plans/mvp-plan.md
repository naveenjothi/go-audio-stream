# 📅 30-Day Backend Roadmap — Personal Audio Streaming Platform (Go + gRPC + WebRTC + Kafka)

> Goal → Build a backend that allows users to sign into the app, connect their phone to the web, and stream audio stored *locally on their device*.

---

## Week 1 — Core Backend & Auth Foundation

### 🎯 Objective  
Build the base backend service with authentication, database structure, and user session management.



### Tasks
- [ ] Setup repo structure: `cmd/`, `internal/`, `pkg/`, `proto/`, `api/`
- [ ] Implement **User Auth**
  - [ ] Signup/Login API
  - [ ] JWT-based session tokens
  - [ ] Refresh token support
- [ ] Database Schema
  - Users
  - Device pairings
  - Song metadata (no file storage yet)
- [ ] Redis integration for session/cache
- [ ] Health check + logging + config env structure
- [ ] CI/CD basic pipeline

### Deliverables
- 🔐 Login + token generation working
- 🗄 Basic DB models + migrations
- 🧠 API documentation for auth + metadata

---

## Week 2 — Mobile Upload + gRPC Metadata Sync

### 🎯 Objective  
Allow the mobile app to scan and sync audio metadata to server.

### Tasks
- [ ] Define Protobuf schemas
  - `SongMetadata`
  - `DeviceInfo`
  - `SyncRequest/SyncResponse`
- [ ] Implement gRPC service for metadata sync
- [ ] Build `grpc-gateway` REST layer (optional)
- [ ] Store synced tracks in DB linked by user ID
- [ ] API to fetch user songs via browser

### Deliverables
- 📡 gRPC endpoints for metadata sync
- 🗂 User song library visible in web UI
- 🔄 Automatic resync when library changes

---

## Week 3 — WebRTC Streaming + Device Pairing

### 🎯 Objective  
Enable streaming from phone → browser using WebRTC.

### Tasks
- [ ] Implement signaling server (WebSockets or gRPC streams)
- [ ] Device pairing:
  - [ ] QR scan or link-code authentication
  - [ ] Store active connections
- [ ] Establish WebRTC PeerConnection
- [ ] Implement audio chunk capture & stream from phone
- [ ] Browser client receives & plays audio

### Deliverables
- 🔗 Phone connects to browser securely
- 🎶 Song streams realtime with <500ms latency
- 🧪 Ability to play single track end-to-end

---

## Week 4 — Kafka + Events + Remote Controller

### 🎯 Objective  
Add multi-device sync + playback controls.

### Tasks
- [ ] Kafka topics:
  - playback.events
  - device.status
  - now.playing
- [ ] Produce/consume events:
  - play / pause
  - seek
  - device online/offline
- [ ] Browser acts as **remote control**
  - Next track
  - Volume change
  - Pause/resume
- [ ] Playback status mirrored on both devices

### Deliverables
- 🛰 Real-time event sync (Kafka)
- 🎛 Full remote playback controls
- 📱 Web controls reflect instantly on phone

---

## Stretch Goals (if time left)

| Feature | Value |
|--------|-------|
| Upload library → S3 | true personal cloud streaming |
| AI auto-tagging metadata | album → genre → artist indexing |
| Listen-stats like Spotify Wrapped | fun analytics layer |
| Multiple device streaming | 1 mobile → 3 browsers |

---

## Final Result After 1 Month

| Capability | Status |
|---|---|
| Auth & DB | 🟢 Completed |
| gRPC metadata sync | 🟢 Completed |
| Device pairing | 🟢 Completed |
| WebRTC audio streaming | 🟢 Completed |
| Kafka playback sync | 🟢 Completed |

> You walk away with production-level experience in  
**Go, gRPC, WebRTC, Kafka, realtime media streaming.**

---
