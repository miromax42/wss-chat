-- name: CreateMessage :one
INSERT INTO messages (
  sender,
  room,
  payload
) VALUES (
  $1,$2,$3
) RETURNING *;

-- name: GetMessages :many
SELECT * 
FROM messages 
WHERE room=$1
AND created_at >= @ago
ORDER BY created_at;

