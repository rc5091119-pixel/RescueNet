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
SELECT
    m.id,
    m.room_id,
    m.sender_id,
    u.name AS sender_name,
    m.content,
    m.created_at
FROM messages m
JOIN users u
    ON m.sender_id = u.id
WHERE m.room_id = $1
ORDER BY m.created_at ASC;

-- name: IsRoomMember :one
SELECT EXISTS(
    SELECT 1
    FROM room_members
    WHERE room_id = $1
    AND user_id = $2
);

-- name: GetMessageByID :one
SELECT
    m.id,
    m.room_id,
    m.sender_id,
    u.name AS sender_name,
    m.content,
    m.created_at
FROM messages m
JOIN users u
    ON m.sender_id = u.id
WHERE m.id = $1;