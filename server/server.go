package server

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	db "github.com/miromax42/wss-chat/db/sqlc"
	_ "github.com/miromax42/wss-chat/docs" //nolint:revive //for swaggo
	"github.com/miromax42/wss-chat/util"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"golang.org/x/net/context"
)

// @title           Simple websocket chat API
// @version         1.0
// @description     This is a sample server of wss-chat.

// @host      localhost:8080
// @BasePath  /
type Server struct {
	ctx      context.Context
	config   util.Config
	Srv      *http.Server
	store    *db.Queries
	router   *gin.Engine
	upgrader websocket.Upgrader
	hubs     *MapRW
}

func New(ctx context.Context, config util.Config, store *db.Queries) (*Server, error) { //nolint:gocritic //config
	server := &Server{
		ctx:    ctx,
		config: config,
		store:  store,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		hubs: &MapRW{
			m: make(map[string]*Hub),
		},
	}

	server.setupRoutes()

	portString := fmt.Sprintf(":%d", server.config.ServerPort)
	server.Srv = &http.Server{
		Addr:    portString,
		Handler: server.router,
	}

	return server, nil
}

func (s *Server) setupRoutes() {
	router := gin.Default()

	router.GET("/ws", s.wsEnpoint)
	router.GET("/rooms", s.getRooms)

	router.Use(static.Serve("/", static.LocalFile("./public", true)))
	router.NoRoute(func(c *gin.Context) {
		c.File("./public/index.html")
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.router = router
}

type Error struct {
	Error string `json:"error"`
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
