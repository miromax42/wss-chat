-- name: CreateMessage :one
INSERT INTO messages (
  sender,
  room,
  payload
) VALUES (
  $1,$2,$3
) RETURNING *;
