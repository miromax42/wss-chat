-- name: CreateRoom :one
INSERT INTO rooms (
  name,
  creator
) VALUES (
  $1, $2
) RETURNING *;