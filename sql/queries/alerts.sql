-- name: CreateAlert :one
INSERT INTO alerts (user_id, latitude, longitude, status)
VALUES ($1, $2, $3, 'active')
RETURNING id, user_id, latitude, longitude, status, created_at;