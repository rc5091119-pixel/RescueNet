-- name: CreateMessage :one
INSERT INTO messages (
    room_id,
    sender_id,
    content
)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetRoomMessages :many
SELECT *
FROM messages
WHERE room_id = $1
ORDER BY created_at ASC;

-- name: IsRoomMember :one
SELECT EXISTS(
    SELECT 1
    FROM room_members
    WHERE room_id = $1
    AND user_id = $2
);