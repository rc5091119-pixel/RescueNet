# RescueNet

A backend system for real-time emergency alerting, built in Go. When a user triggers an SOS, RescueNet looks up their location from a dedicated `user_locations` table, finds nearby users within a 1 km radius using the Haversine formula, and returns the alert along with the list of nearby responders. Other users can accept the alert — up to a cap of 15 responders per emergency.

> Work in progress — chat and resolution features are planned (see [RoadMap.md](./RoadMap.md)).

---

## How It Works

1. User updates their location via `PATCH /api/location` (stored in a separate `user_locations` table)
2. User triggers an SOS via `POST /api/alerts` — the server reads their location from `user_locations`, creates an alert record, and queries all other user locations
3. The Haversine function filters users within 1 km and returns them as `nearby_users`
4. Nearby users can accept the alert via `PATCH /api/alerts/{id}/accept` — limited to 15 responders per alert

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go (Golang) |
| HTTP | `net/http` (standard library) |
| Auth | JWT (JSON Web Tokens) |
| Database | PostgreSQL |
| Query gen | sqlc |
| Location math | Haversine (custom implementation) |
| Config | godotenv |

---

## Project Structure

```
rescuenet/
├── internal/
│   └── database/               # sqlc-generated DB queries and types
├── sql/                        # SQL schema and migration files
├── authmiddleware.go           # JWT auth middleware, injects userID into context
├── handler_create_users.go     # POST /api/users — register
├── handler_login.go            # POST /api/login — login, returns JWT
├── handler_update_locatio.go   # PATCH /api/location — upsert user location
├── handler_get_user_locations.go  # POST /api/alerts — create alert + find nearby users
├── handlelr_accept_alert.go    # PATCH /api/alerts/{id}/accept — accept an alert (max 15)
├── handler_test_protected.go   # GET /api/test — auth check endpoint
├── haversine.go                # Distance calculation between two coordinates
├── json.go                     # respondWithJSON / respondWithError helpers
├── main.go                     # Server setup, routes, DB init
├── sqlc.yaml
├── go.mod
└── go.sum
```

---

## API Endpoints

### Auth

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/api/users` | No | Register a new user |
| `POST` | `/api/login` | No | Login and receive a JWT token |

### Location

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `PATCH` | `/api/location` | Yes | Update the authenticated user's location |

### Alerts

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/api/alerts` | Yes | Trigger an SOS — creates alert from user's stored location, returns nearby users |
| `PATCH` | `/api/alerts/{id}/accept` | Yes | Accept an alert as a responder (max 15 per alert) |

> All protected routes require `Authorization: Bearer <token>` in the request header.

---

## Database Design

Locations are stored in a **separate `user_locations` table** rather than on the `users` table, so they can be updated independently and queried efficiently.

```sql
-- Users
CREATE TABLE users (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT NOT NULL,
    email      TEXT UNIQUE NOT NULL,
    password   TEXT NOT NULL,   -- bcrypt hashed
    created_at TIMESTAMP DEFAULT now()
);

-- User locations (separate table, updated via PATCH /api/location)
CREATE TABLE user_locations (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    latitude   FLOAT NOT NULL,
    longitude  FLOAT NOT NULL,
    updated_at TIMESTAMP DEFAULT now()
);

-- Alerts (emergency SOS events)
CREATE TABLE alerts (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID REFERENCES users(id),
    latitude   FLOAT NOT NULL,
    longitude  FLOAT NOT NULL,
    status     TEXT DEFAULT 'active',
    created_at TIMESTAMP DEFAULT now()
);

-- Alert responses (who accepted the alert)
CREATE TABLE alert_responses (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    alert_id   UUID REFERENCES alerts(id),
    user_id    UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT now()
);
```

---

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 14+

### Run Locally

```bash
# Clone the repo
git clone https://github.com/rc5091119-pixel/RescueNet.git
cd RescueNet

# Install dependencies
go mod tidy

# Set environment variables
cp .env.example .env
# Fill in DB_URL and JWT_SECRET in .env

# Apply SQL migrations from sql/ folder

# Start the server
go run .
```

Server starts at `http://localhost:8080`

---

## Environment Variables

```env
DB_URL=postgres://user:password@localhost:5432/rescuenet?sslmode=disable
JWT_SECRET=your_secret_key
```

---

## Example Requests

**Register**
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Ravi", "email": "ravi@example.com", "password": "secret"}'
```

**Update Location**
```bash
curl -X PATCH http://localhost:8080/api/location \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"lat": 26.8467, "lng": 80.9462}'
```

**Trigger SOS**
```bash
curl -X POST http://localhost:8080/api/alerts \
  -H "Authorization: Bearer <token>"
```

**Accept Alert**
```bash
curl -X PATCH http://localhost:8080/api/alerts/<alert-id>/accept \
  -H "Authorization: Bearer <token>"
```

---

## Roadmap

See [RoadMap.md](./RoadMap.md) for planned features.

---

## Author

**Ravindra Choudhary**
B.Tech — Electronics and Communication Engineering, NIT Agartala

- GitHub: [rc5091119-pixel](https://github.com/rc5091119-pixel)
- Email: rc5091119@gmail.com