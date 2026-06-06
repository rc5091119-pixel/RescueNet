-- name: CreateRoom :exec
INSERT INTO rooms(id, alert_id)
VALUES($1, $2);

-- name: AddRoomMember :exec
INSERT INTO room_members(room_id, user_id)
VALUES($1, $2);

-- name: GetRoomMembers :many
SELECT user_id
FROM room_members
WHERE room_id = $1;

-- name: GetUserRooms :many
SELECT
    r.id AS room_id,
    u.name AS creator_name,
    a.latitude,
    a.longitude
FROM room_members rm
JOIN rooms r
    ON rm.room_id = r.id
JOIN alerts a
    ON r.alert_id = a.id
JOIN users u
    ON a.user_id = u.id
WHERE rm.user_id = $1;

-- name: GetRoomInfo :one
SELECT
    r.id,
    r.alert_id,
    u.name AS creator_name,
    a.latitude,
    a.longitude
FROM rooms r
JOIN alerts a
    ON r.alert_id = a.id
JOIN users u
    ON a.user_id = u.id
WHERE r.id = $1;