package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type wsRequest struct {
	Room        string `form:"room" binding:"required"`
	HistoryTime string `form:"time"`
}

// wsEnpoint godoc
// @Summary      WebSocket chat
// @Description  websocket chat endpoint to connect specified room
// @Param        room  query  string  true   "Room to connect via ws"
// @Param        time  query  string  false  "Time to load history of messages (1h, 5m etc.)"
// @Tags         chat
// @Produce      json
// @Router       /ws [get]
func (s *Server) wsEnpoint(ctx *gin.Context) {
	// parse params
	var req wsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	historyTime, err := time.ParseDuration(req.HistoryTime)
	if err != nil {
		historyTime = time.Minute
	}

	// create room in db
	_, err = s.store.CreateRoom(ctx, req.Room)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
			default:
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))

				return
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))

			return
		}
	}

	// upgrade request
	s.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := s.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("couldnt upgrade: %s", err)

		return
	}

	// start hub if not exist
	if _, ok := s.hubs.Load(req.Room); !ok {
		hub := newHub(req.Room)
		s.hubs.Store(req.Room, hub)

		go hub.run()
	}

	client := &Client{hub: s.hubs.LoadForce(req.Room), conn: ws, send: make(chan Message, 256), time: historyTime}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go s.writePump(client)
	go s.readPump(client)

	log.Printf("Client Connected with req: %s\n", req)
}
