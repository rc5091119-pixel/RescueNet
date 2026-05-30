-- name: CreateAlertNotification :exec
INSERT INTO alert_notifications (alert_id,user_id)
VALUES ($1, $2)
ON CONFLICT (alert_id, user_id) DO NOTHING;

-- name: GetPendingAlertsForUser :one
SELECT *
FROM alert_notifications
WHERE user_id = $1
AND status = 'pending';