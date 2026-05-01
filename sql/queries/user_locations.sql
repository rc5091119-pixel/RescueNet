-- name: UpdateUserLocations :exec
INSERT INTO user_locations(user_id,latitude,longitude,updated_at)
VALUES(
    $1,
    $2,
    $3,
    NOW()
)
ON CONFLICT (user_id)
DO UPDATE SET
latitude = EXCLUDED.latitude,
longitude = EXCLUDED.longitude,
updated_at = NOW();

-- name: GetUserLocationByUserID :one
SELECT user_id,latitude,longitude
FROM user_locations
WHERE user_id = $1;

-- name: GetUserLocations :many
SELECT user_id, latitude, longitude
FROM user_locations
WHERE latitude BETWEEN sqlc.arg(lat)::float8 - 0.01 AND sqlc.arg(lat)::float8 + 0.01
AND longitude BETWEEN sqlc.arg(lng)::float8 - 0.01 AND sqlc.arg(lng)::float8 + 0.01;
