package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// getRooms godoc
// @Summary      Get all rooms
// @Description  get all rooms from db
// @Tags         room
// @Produce      json
// @Success      200  {array}   db.Room
// @failure      500  {object}  Error
// @Router       /rooms [get]
func (s *Server) getRooms(ctx *gin.Context) {
	rooms, err := s.store.GetRooms(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, rooms)
}
