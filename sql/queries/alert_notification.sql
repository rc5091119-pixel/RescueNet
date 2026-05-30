-- name: CreateAlertNotification :exec
INSERT INTO alert_notifications (alert_id,user_id)
VALUES ($1, $2)
ON CONFLICT (alert_id, user_id) DO NOTHING;