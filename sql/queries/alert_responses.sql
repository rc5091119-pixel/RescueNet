-- name: CreateAlertResponse :exec
INSERT INTO alert_responses (alert_id, user_id, status)
VALUES ($1, $2, 'accepted');

-- name: CountAcceptedUsers :one
SELECT COUNT(*)
FROM alert_responses
WHERE alert_id = $1 AND status = 'accepted';

-- name: GetAcceptedUsers :many
SELECT user_id
FROM alert_responses
WHERE alert_id = $1 AND status = 'accepted';