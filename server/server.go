package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	db "github.com/miromax42/wss-chat/db/sqlc"
	"github.com/miromax42/wss-chat/util"
)

type Server struct {
	config   util.Config
	store    *db.Queries
	upgrader websocket.Upgrader
}

func New(config util.Config, store *db.Queries) (*Server, error) { //nolint:gocritic //config
	return &Server{
		config: config,
		store:  store,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}, nil
}

func (s *Server) Start(port int) error {
	s.configureRoutes()

	return http.ListenAndServe(fmt.Sprintf(":%d", s.config.ServerPort), nil)
}

func (s *Server) configureRoutes() {
	http.HandleFunc("/", s.homePage)
	http.HandleFunc("/ws", s.wsEnpoint)
}
