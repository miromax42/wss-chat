// Code generated by sqlc. DO NOT EDIT.
// source: message.sql

package db

import (
	"context"
)

const createMessage = `-- name: CreateMessage :one
INSERT INTO messages (
  sender,
  room,
  payload
) VALUES (
  $1,$2,$3
) RETURNING id, sender, room, payload, created_at
`

type CreateMessageParams struct {
	Sender  string `json:"sender"`
	Room    string `json:"room"`
	Payload string `json:"payload"`
}

func (q *Queries) CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error) {
	row := q.db.QueryRowContext(ctx, createMessage, arg.Sender, arg.Room, arg.Payload)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.Sender,
		&i.Room,
		&i.Payload,
		&i.CreatedAt,
	)
	return i, err
}