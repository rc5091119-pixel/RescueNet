# 🚨 RescueNet

> A real-time emergency alert and coordination backend built in Go — connecting people in crisis with nearby responders instantly.

---

## What is RescueNet?

RescueNet is a production-ready backend system for community emergency response. When a user triggers an SOS, the system uses the **Haversine formula** to find every user within 1 km, sends them alert notifications, and automatically wires everyone into a **live WebSocket room** for real-time chat and location sharing.

The entire system is built on Go's standard `net/http` library — no framework, no magic. JWT auth, PostgreSQL via sqlc, goroutine-safe WebSocket hub.

---

## How it works

```
User triggers SOS (POST /api/alerts)
        │
        ▼
Server reads user's stored GPS coordinates
        │
        ▼
Haversine algorithm scans all users within 1 km radius
        │
        ▼
Alert notifications created for each nearby user
        │
        ▼
Nearby users see alert via GET /api/notifications
        │
        ▼
A user accepts the alert (POST /api/alerts/{id}/accept)
        │
        ▼
Chat room is created — all responders join via WebSocket
        │
        ▼
Real-time chat + live location sharing until emergency resolved
```

---

## Tech stack

| Layer | Technology |
|---|---|
| Language | Go 1.21+ |
| Web framework | `net/http` standard library |
| Authentication | JWT (golang-jwt) |
| Password hashing | bcrypt |
| Database | PostgreSQL 14+ |
| DB query layer | sqlc (type-safe generated queries) |
| WebSocket | gorilla/websocket |
| UUID | google/uuid |
| Config | godotenv |

---

## API reference

### Public endpoints

| Method | Route | Description |
|---|---|---|
| `POST` | `/api/users` | Register a new user |
| `POST` | `/api/login` | Login — returns JWT token |

### Protected endpoints (require `Authorization: Bearer <token>`)

| Method | Route | Description |
|---|---|---|
| `POST` | `/api/location` | Update your current GPS coordinates |
| `POST` | `/api/alerts` | Trigger an SOS — finds nearby users and sends notifications |
| `GET` | `/api/notifications` | Get pending alert notifications for the current user |
| `POST` | `/api/alerts/{id}/accept` | Accept an alert (max 3 responders per alert) |
| `GET` | `/api/my-rooms` | Get all chat rooms the current user belongs to |
| `GET` | `/api/rooms/{roomID}/info` | Get info about a specific room |
| `GET` | `/api/rooms/{roomID}/messages` | Fetch message history for a room |
| `POST` | `/api/rooms/{roomID}/messages` | Send a message to a room (REST fallback) |
| `GET` | `/api/rooms/{roomID}/locations` | Get live locations of all room members |

### WebSocket

| Route | Description |
|---|---|
| `GET /ws/rooms/{roomID}?token=<jwt>` | Connect to a room's real-time channel |

Once connected, send JSON messages:

```json
// Chat message — saved to DB and broadcast to all room members
{ "type": "chat", "content": "I'm 2 minutes away" }

// Location update — broadcast to all room members (not persisted)
{ "type": "location", "latitude": 26.9124, "longitude": 75.7873 }
```

---

## Database schema

```sql
-- Users
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          TEXT,
    email         TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMP DEFAULT now()
);

-- User locations (updated separately via /api/location)
CREATE TABLE user_locations (
    user_id   UUID REFERENCES users(id),
    latitude  FLOAT NOT NULL,
    longitude FLOAT NOT NULL,
    updated_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (user_id)
);

-- Alerts (SOS events)
CREATE TABLE alerts (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID REFERENCES users(id),
    latitude   FLOAT NOT NULL,
    longitude  FLOAT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

-- Alert responses (who accepted)
CREATE TABLE alert_responses (
    alert_id UUID REFERENCES alerts(id),
    user_id  UUID REFERENCES users(id),
    PRIMARY KEY (alert_id, user_id)
);

-- Alert notifications (who was notified)
CREATE TABLE alert_notifications (
    alert_id   UUID REFERENCES alerts(id),
    user_id    UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (alert_id, user_id)
);

-- Chat rooms
CREATE TABLE rooms (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP DEFAULT now()
);

-- Room members
CREATE TABLE room_members (
    room_id UUID REFERENCES rooms(id),
    user_id UUID REFERENCES users(id),
    PRIMARY KEY (room_id, user_id)
);

-- Messages
CREATE TABLE messages (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id    UUID REFERENCES rooms(id),
    sender_id  UUID REFERENCES users(id),
    content    TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);
```

---

## Getting started

### Prerequisites

- Go 1.21+
- PostgreSQL 14+
- sqlc (`go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`)

### Setup

```bash
# 1. Clone
git clone https://github.com/rc5091119-pixel/RescueNet.git
cd RescueNet

# 2. Install dependencies
go mod tidy

# 3. Configure environment
cp temp.env .env
# Edit .env — set DB_URL, JWT_SECRET, ALLOWED_ORIGIN

# 4. Run DB migrations (apply SQL files in /sql)

# 5. Start server
go run .
```

Server starts at `http://localhost:8080`

---

## Environment variables

```env
DB_URL=postgres://user:password@localhost:5432/rescuenet?sslmode=disable
JWT_SECRET=your_secret_key_here
ALLOWED_ORIGIN=http://localhost:5173
```

---

## Project structure

```
RescueNet/
├── main.go                          # Server entry point, all routes registered here
├── authmiddleware.go                # JWT validation middleware + room membership check
├── corsMiddleware.go                # CORS headers middleware
├── hub.go                           # WebSocket hub — room map with RWMutex
├── haversine.go                     # Haversine distance formula (km)
├── json.go                          # respondWithJSON / respondWithError helpers
├── handler_create_users.go          # POST /api/users
├── handler_login.go                 # POST /api/login
├── handler_update_location.go       # POST /api/location
├── handler_create_alert.go          # POST /api/alerts
├── handler_accept_alert.go          # POST /api/alerts/{id}/accept
├── handler_get_notifications.go     # GET /api/notifications
├── handler_websocket.go             # GET /ws/rooms/{roomID} — WebSocket handler
├── handler_create_messages.go       # POST /api/rooms/{roomID}/messages
├── handler_get_messages.go          # GET /api/rooms/{roomID}/messages
├── handler_get_room_info.go         # GET /api/rooms/{roomID}/info
├── handler_get_room_members_locations.go  # GET /api/rooms/{roomID}/locations
├── handler_get_users_rooms.go       # GET /api/my-rooms
├── internal/
│   ├── auth/                        # JWT creation and validation, bcrypt hashing
│   └── database/                    # sqlc-generated type-safe DB queries
├── sql/                             # SQL migration files and sqlc query definitions
├── sqlc.yaml                        # sqlc configuration
└── go.mod
```

---

## WebSocket architecture

The hub is a single in-memory struct holding all active rooms:

```go
type Hub struct {
    Rooms map[uuid.UUID]map[*Client]bool
    mu    sync.RWMutex
}
```

Each WebSocket connection is a `Client` with its `UserID` and `RoomID`. When a message arrives, the handler acquires a read lock and broadcasts to all other clients in the same room. Disconnects are cleaned up with `defer` — if a room empties, it's removed from the map.

Authentication for WebSocket connections happens via `?token=<jwt>` query parameter (since browsers cannot set custom headers on WebSocket upgrades). Room membership is verified against the database before the upgrade proceeds.

---

## Key design decisions

**Why `net/http` instead of Gin/Echo?** Zero external dependencies for the router. Go 1.22's pattern matching (`GET /api/rooms/{roomID}/messages`) handles everything needed.

**Why sqlc?** Type-safe SQL queries with zero reflection. DB schema is the source of truth — changing a query in SQL automatically updates the Go types.

**Why in-memory hub instead of Redis pub/sub?** Single-server deployment doesn't need distributed messaging. The `sync.RWMutex` pattern handles concurrent reads (broadcasts) efficiently.

**Why Haversine in Go instead of PostGIS?** Keeps the DB dependency minimal. At the scale of 1 km radius lookups the full table scan with Go-side filtering is fast enough. Can be replaced with a PostGIS `ST_DWithin` index if scale demands it.

---

## Author

**Ravindra Choudhary**
B.Tech — Electronics and Communication Engineering
National Institute of Technology Agartala | GPA: 8.72

- Email: rc5091119@gmail.com
- GitHub: [github.com/rc5091119-pixel](https://github.com/rc5091119-pixel)