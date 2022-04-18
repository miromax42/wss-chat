-- name: CreateRoom :one
INSERT INTO rooms (
  name
) VALUES (
  $1
) RETURNING *;

-- name: GetRooms :many
SELECT * FROM rooms
ORDER BY name;