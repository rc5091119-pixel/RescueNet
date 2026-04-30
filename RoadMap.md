# 🚀 RescueNet – Complete Backend Workflow (Production-Ready)

## 🎯 Project Goal

User presses SOS → nearby users (within ~1 km) get notified → one user accepts → alert is resolved.

---

# 🗺️ Phase 0: Project Setup (Day 1–2)

* Initialize Go project (net/http or Gin)
* Setup PostgreSQL database
* Create `.env` for secrets (DB URL, JWT secret)
* Setup basic HTTP server
* Project structure:

/cmd
/internal
  /handlers
  /services
  /models
  /db
  /auth

---

# 🔐 Phase 1: Authentication System (Day 3–5)

## Features:

* User Signup → `POST /api/users`
* User Login → `POST /api/login`
* JWT token generation
* Auth middleware (protected routes)

## DB Schema: users

* id (UUID)
* email (unique)
* password_hash
* created_at

## Notes:

* Use bcrypt or Argon2 for hashing
* Store only hashed passwords

---

# 🌍 Phase 2: Location System (Day 6–8)

## Features:

* Update user location → `POST /api/location`

### Request Body:

{
"lat": 23.8,
"lng": 91.2
}

## DB Schema: user_locations

* user_id (FK)
* latitude
* longitude
* updated_at

## Notes:

* Keep latest location only
* Add index on (latitude, longitude)

---

# 🚨 Phase 3: Alert + Nearby System (CORE) (Day 9–14)

## Features:

* Trigger emergency → `POST /api/alerts`

## DB Schema: alerts

* id
* user_id
* latitude
* longitude
* status (active / accepted / resolved)
* created_at

---

## 🔥 Core Logic (VERY IMPORTANT)

1. User triggers alert
2. Backend fetches all user locations
3. Calculate distance using Haversine formula
4. Filter users within ≤ 1 km
5. Return nearby users (for testing phase)

---

# 📢 Phase 4: Notification System (Day 15–20)

## Step 1 (Testing):

* Return nearby users in API response

## Step 2 (Real-Time System):

* Implement WebSocket server

## Flow:

* Users connect via WebSocket
* Store active connections (user_id → connection)
* On alert:
  → Send event to nearby users instantly

---

# 🤝 Phase 5: Accept Help System (Day 21–24)

## Features:

* Accept alert → `POST /api/alerts/{id}/accept`

---

## DB Schema: alert_responses

* alert_id
* user_id
* status (accepted / rejected)
* responded_at

---

## 🔒 Critical Logic:

* First user to accept → wins
* Update alert status → accepted
* Lock alert (prevent multiple accepts)

---

# 🔐 Phase 6: Security & Validation (Day 25–28)

* JWT expiry handling
* Input validation (lat/lng, email, etc.)
* Rate limiting (prevent spam alerts)
* Error handling (important)

---

# ⚡ Phase 7: Optimization (Day 29–32)

* Add DB indexes:

  * user_locations (lat, lng)
  * alerts (status)

* Optimize queries

* Optional: Use PostGIS for geo queries

---

# 📄 Phase 8: Final Polish (Day 33–35)

* Clean README.md
* Add API documentation
* Add architecture diagram
* Add demo video (very important)

---

# 🧠 System Flow (Final Architecture)

User → Login → JWT
User → Update Location
User → Send Alert
Backend → Find Nearby Users (Haversine)
Backend → Send Notifications (WebSocket)
Nearby User → Accept Alert
Backend → Lock + Resolve Alert

---

# 🔧 Tech Stack

* Language: Go (Golang)
* Framework: net/http / Gin
* Database: PostgreSQL
* Auth: JWT
* Password Hashing: Argon2 / bcrypt
* Real-Time: WebSockets

---

# ⚠️ Rules to Follow

* Build MVP first (don’t overcomplicate)
* Do NOT skip authentication
* Do NOT jump to frontend early
* Always test APIs using Postman
* Keep code modular and clean

---

# 📅 Weekly Plan

Week 1 → Auth + DB
Week 2 → Location + Alerts
Week 3 → Nearby Logic
Week 4 → WebSocket Notifications
Week 5 → Optimization + README

---

# 🔥 One-Line Summary

Build authentication → store user location → trigger alert → find nearby users → send real-time notifications → allow one user to accept → resolve alert → optimize and polish.
