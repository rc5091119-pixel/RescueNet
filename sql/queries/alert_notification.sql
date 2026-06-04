-- name: CreateAlertNotification :exec
INSERT INTO alert_notifications (alert_id,user_id)
VALUES ($1, $2)
ON CONFLICT (alert_id, user_id) DO NOTHING;

-- name: GetPendingAlertsForUser :many
SELECT
    an.id,
    an.alert_id,
    an.status,
    an.user_id,
    u.name AS creator_name
FROM alert_notifications an
JOIN alerts a
    ON an.alert_id = a.id
JOIN users u
    ON a.user_id = u.id
WHERE an.user_id = $1
AND an.status = 'pending';

-- name: MarkNotificationAccepted :exec
UPDATE alert_notifications
SET status = 'accepted'
WHERE alert_id = $1
AND user_id = $2;