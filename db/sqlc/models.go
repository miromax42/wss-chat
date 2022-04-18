// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"time"
)

type Message struct {
	ID        int32     `json:"id"`
	Sender    string    `json:"sender"`
	Room      string    `json:"room"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

type Room struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
