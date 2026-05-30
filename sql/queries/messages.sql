-- name: CreateMessage :exec
INSERT INTO messages (
    id,
    room_id,
    sender_id,
    content
)
VALUES (
    $1,
    $2,
    $3,
    $4
);

-- name: GetRoomMessages :many
SELECT *
FROM messages
WHERE room_id = $1
ORDER BY created_at ASC;

-- name: GetRecentRoomMessages :many
SELECT *
FROM messages
WHERE room_id = $1
ORDER BY created_at DESC
LIMIT 50;