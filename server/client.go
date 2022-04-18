package server

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	db "github.com/miromax42/wss-chat/db/sqlc"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan Message

	time time.Duration
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (s *Server) readPump(c *Client) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))                                                           //nolint:errcheck //sample
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil }) //nolint:errcheck,nlreturn,lll //sample

	arg := db.GetMessagesParams{
		Room: c.hub.name,
		Ago:  time.Now().Add(-c.time),
	}
	messages, err := s.store.GetMessages(s.ctx, arg)
	if err != nil { //nolint:wsl //ok
		log.Printf("error: %v", err)

		return
	}

	for _, message := range messages {
		c.hub.broadcast <- Message{
			Sender:  message.Sender,
			Payload: message.Payload,
		}
	}

	c.hub.broadcast <- Message{
		Sender:  "admin",
		Payload: "client connected",
	}

	for {
		var message Message

		err := c.conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}

			break
		}

		// add to db
		_, err = s.store.CreateMessage(s.ctx, db.CreateMessageParams{
			Sender:  message.Sender,
			Room:    c.hub.name,
			Payload: message.Payload,
		})
		if err != nil {
			log.Printf("insert: %v", err)

			continue
		}

		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (s *Server) writePump(c *Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) //nolint:errcheck //sample

			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{}) //nolint:errcheck //sample

				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				log.Println(err)

				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) //nolint:errcheck //sample

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
