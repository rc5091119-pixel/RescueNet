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
SELECT room_id
FROM room_members
WHERE user_id = $1;