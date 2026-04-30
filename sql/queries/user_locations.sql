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

-- name: GetUserLocations :many
SELECT user_id,latitude,longitude
FROM user_locations
WHERE latitude BETWEEN $1 - 0.01 AND $1 + 0.01
AND longitude BETWEEN $2 - 0.01 AND $2 + 0.01;
