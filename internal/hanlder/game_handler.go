package hanlder

import (
	"tic-tac-toe-game/internal/model"
	"tic-tac-toe-game/internal/service"
	"time"

	socketio "github.com/googollee/go-socket.io"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	GamerService service.GameService
}

func NewGameHandler(gamerService service.GameService) *GameHandler {
	return &GameHandler{GamerService: gamerService}
}

func (g *GameHandler) GetAllRooms(c *gin.Context) {
	AllRooms := g.GamerService.GetAllRooms()
	c.JSON(200, gin.H{"rooms": AllRooms})
}

func (g *GameHandler) CreateGameRoom(c *gin.Context) {
	gameRoom := model.Game{
		RoomID:    "",
		Players:   []model.Player{},
		Board:     [3][3]string{},
		Winner:    "",
		Turn:      "",
		IsOver:    false,
		IsDraw:    false,
		CreatedAt: time.Now(),
	}

	createdRoom := g.GamerService.CreateGameRoom(gameRoom)

	c.JSON(200, gin.H{"room": createdRoom})
}

func (g *GameHandler) JoinGameRoom(ClientID string, roomID string, s *socketio.Server) error {
	err := g.GamerService.JoinGameRoom(roomID, ClientID)
	if err != nil {
		return err
	}

	game := g.GamerService.FindGameRoom(roomID)
	if game.Turn == "" && len(game.Players) > 0 {
		game.Turn = game.Players[0].ClientID // เริ่มเล่นที่คนแรก
		g.GamerService.UpdateGameRoom(game)  // ฟังก์ชันนี้ต้องเขียนเพิ่มใน service
	}

	if g.GamerService.GameReady(roomID) {
		payload := model.StartGamePayload{
			RoomID:  game.RoomID,
			Players: game.Players,
			Turn:    game.Turn,
		}
		s.BroadcastToRoom("/", roomID, "start-game", payload)
	}
	return nil
}

func (g *GameHandler) LeaveGameRoom(clientID string, roomID string, s *socketio.Server, conn socketio.Conn) error {
	err := g.GamerService.LeaveGameRoom(roomID, clientID)
	if err != nil {
		return err
	}

	conn.Leave(roomID)

	s.BroadcastToRoom("/", roomID, "opponent-left", gin.H{
		"clientID": clientID,
	})
	return nil
}
