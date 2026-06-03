-- name: CreateAlertNotification :exec
INSERT INTO alert_notifications (alert_id,user_id)
VALUES ($1, $2)
ON CONFLICT (alert_id, user_id) DO NOTHING;

-- name: GetPendingAlertsForUser :many
SELECT *
FROM alert_notifications
WHERE user_id = $1
AND status = 'pending';

-- name: MarkNotificationAccepted :exec
UPDATE alert_notifications
SET status = 'accepted'
WHERE alert_id = $1
AND user_id = $2;