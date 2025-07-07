package router

import (
	"tic-tac-toe-game/internal/hanlder"

	"github.com/gin-gonic/gin"
)

type GameRouter struct {
	GameHandler *hanlder.GameHandler
}

func NewGameRouter(gameHandler *hanlder.GameHandler) *GameRouter {
	return &GameRouter{GameHandler: gameHandler}
}

func (r *GameRouter) InitRouter(rg *gin.RouterGroup) {
	rg.GET("/rooms", r.GameHandler.GetAllRooms)
}
