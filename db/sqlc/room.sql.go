// Code generated by sqlc. DO NOT EDIT.
// source: room.sql

package db

import (
	"context"
)

const createRoom = `-- name: CreateRoom :one
INSERT INTO rooms (
  name
) VALUES (
  $1
) RETURNING name, created_at
`

func (q *Queries) CreateRoom(ctx context.Context, name string) (Room, error) {
	row := q.db.QueryRowContext(ctx, createRoom, name)
	var i Room
	err := row.Scan(&i.Name, &i.CreatedAt)
	return i, err
}

const getRooms = `-- name: GetRooms :many
SELECT name, created_at FROM rooms
ORDER BY name
`

func (q *Queries) GetRooms(ctx context.Context) ([]Room, error) {
	rows, err := q.db.QueryContext(ctx, getRooms)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Room{}
	for rows.Next() {
		var i Room
		if err := rows.Scan(&i.Name, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
