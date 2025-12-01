````md
# 📄 API Contracts — Audio Streaming Backend (Go + gRPC + WebRTC)

## Auth Service (REST)

---

### 🔐 Auth Mechanism (Firebase)

Client must attach Google-signed Firebase ID Token

``Authorization: Bearer <firebase_token>``

Go backend validates token:

```go
idToken := strings.TrimPrefix(authHeader, "Bearer ")
token, err := firebaseAuth.VerifyIDToken(ctx, idToken)
if err != nil { return 401 }
userID := token.UID
```

---
## User Songs + Library Metadata (REST + gRPC)

---

### GET `/songs`

Returns user’s synced songs.

#### Response (200)

```json
[
  {
    "id": "uuid",
    "title": "Something Just Like This",
    "artist": "Coldplay",
    "album": "Kaleidoscope",
    "duration": 230000,
    "deviceId": "mobile-uuid"
  }
]
```

---

### gRPC Service `SongSyncService`

File: `proto/song_sync.proto`

```protobuf
service SongSyncService {
  rpc SyncMetadata(stream SongMetadataRequest) returns (SyncMetadataResponse);
  rpc GetSongs(UserRequest) returns (SongListResponse);
}

message SongMetadataRequest {
  string deviceId = 1;
  string songId = 2;
  string title = 3;
  string album = 4;
  string artist = 5;
  int32 duration = 6; // milliseconds
}

message SyncMetadataResponse {
  int32 totalImported = 1;
  string status = 2; // success | partial | error
}

message UserRequest { string userId = 1; }

message SongListResponse { repeated SongMetadataRequest songs = 1; }
```

---

## Device Pairing + WebRTC Signaling API

---

### POST `/devices/pair/initiate`

Generate pairing QR or code.

#### Request

```json
{
  "deviceType": "mobile"
}
```

#### Response

```json
{
  "pairCode": "89XK-P2",
  "expiresIn": 300
}
```

---

### POST `/devices/pair/verify`

Used by browser to complete pairing.

```json
{
  "pairCode": "89XK-P2"
}
```

#### Response

```json
{
  "status": "paired",
  "mobileDeviceId": "uuid",
  "browserDeviceId": "uuid"
}
```

---

## WebSocket Signaling (for WebRTC)

Endpoint: `ws://api/rtc/signal`

### Outgoing Messages from Client → Server

```json
{
  "type": "offer|answer|candidate|sync",
  "from": "device_id",
  "to": "device_id",
  "sdp": "...optional...",
  "candidate": "...optional..."
}
```

### Incoming Server → Client

```json
{
  "type": "offer|answer|candidate",
  "from": "device_id",
  "payload": { "sdp": "...", "candidate": "..." }
}
```

---

## Kafka Playback Event Contracts

Topic: `playback.events`

```json
{
  "eventType": "PLAY|PAUSE|SEEK|NEXT|PREV|VOLUME",
  "userId": "uuid",
  "deviceId": "uuid",
  "songId": "uuid",
  "timestamp": 1735739273,
  "metadata": {
    "seekMs": 136000,
    "volume": 0.85
  }
}
```

---

## Playback Status Sync API (Web & Mobile use)

WebSocket channel → `ws://api/playback/state`

### Example Stream Messages

```json
{ "type": "PLAY", "songId": "uuid", "positionMs": 3000 }
{ "type": "PAUSE", "songId": "uuid" }
{ "type": "UPDATE_POSITION", "positionMs": 92000 }
```

---

# System Flow Summary

```
📱 Mobile App
     ⬇ Sync song metadata via gRPC
🖥 Web App
     ⬆ fetches songs via REST
🔗 Pair device → QR or code
🔄 WebRTC signaling via WS
🎶 Stream phone → browser P2P
📡 Kafka sync → play/pause/seek
```

---

```
This contract set gives you a fully spec’d backend interface for 1-month execution.
```
