package socket

import (
	"tic-tac-toe-game/internal/hanlder"
	"tic-tac-toe-game/internal/model"

	socketio "github.com/googollee/go-socket.io"
)

var (
	server *socketio.Server
)

func InitSocketServer(gameHandler *hanlder.GameHandler) {
	server = socketio.NewServer(nil)

	server.OnEvent("/", "joinRoom", func(conn socketio.Conn, roomID string) {
		conn.Join(roomID)
		clientID := conn.ID()

		err := gameHandler.JoinGameRoom(clientID, roomID, server)
		if err != nil {
			conn.Emit("error", err.Error())
		}
	})

	server.OnEvent("/", "leaveRoom", func(conn socketio.Conn, roomID string) {
		clientID := conn.ID()
		err := gameHandler.LeaveGameRoom(clientID, roomID, server, conn)
		if err != nil {
			conn.Emit("error", err.Error())
		}
	})

	server.OnEvent("/", "makeMove", func(conn socketio.Conn, data model.Move) {
		// สมมติส่ง roomID มาพร้อม Move
		game := gameHandler.GamerService.FindGameRoom(data.RoomID)
		if game.IsOver {
			conn.Emit("error", "Game is already over")
			return
		}

		// ตรวจสอบเป็นตาของผู้เล่นไหม
		if game.Turn != conn.ID() {
			conn.Emit("error", "Not your turn")
			return
		}

		// อัพเดตบอร์ด
		if data.Row < 0 || data.Row > 2 || data.Col < 0 || data.Col > 2 {
			conn.Emit("error", "Invalid move")
			return
		}

		if game.Board[data.Row][data.Col] != "" {
			conn.Emit("error", "Cell already taken")
			return
		}

		game.Board[data.Row][data.Col] = getPlayerSymbol(game.Players, conn.ID())

		// ตรวจผลชนะ หรือ เสมอ
		if winner := checkWinner(game.Board); winner != "" {
			game.Winner = winner
			game.IsOver = true
			game.Turn = ""
		} else if checkDraw(game.Board) {
			game.IsDraw = true
			game.IsOver = true
			game.Turn = ""
		} else {
			// เปลี่ยนเทิร์น
			game.Turn = getNextPlayerID(game.Players, conn.ID())
		}

		// อัพเดตเกม
		gameHandler.GamerService.UpdateGameRoom(game)

		// แจ้งอัพเดตบอร์ดให้ผู้เล่นทุกคนในห้อง
		server.BroadcastToRoom("/", game.RoomID, "updateBoard", game)
	})

	server.OnEvent("/", "createRoom", func(conn socketio.Conn) {
		createdRoom := gameHandler.CreateGameRoom(server)
		server.BroadcastToRoom("/", createdRoom.RoomID, "room-created", createdRoom)
	})

}

func GetSocketServer() *socketio.Server {
	return server
}

func getPlayerSymbol(players []model.Player, clientID string) string {
	for _, p := range players {
		if p.ClientID == clientID {
			return p.Symbol
		}
	}
	return ""
}

func getNextPlayerID(players []model.Player, currentID string) string {
	if len(players) < 2 {
		return ""
	}
	if players[0].ClientID == currentID {
		return players[1].ClientID
	}
	return players[0].ClientID
}

// ฟังก์ชันเช็คผลชนะ (ง่ายๆ)
func checkWinner(board [3][3]string) string {
	lines := [8][3][2]int{
		{{0, 0}, {0, 1}, {0, 2}},
		{{1, 0}, {1, 1}, {1, 2}},
		{{2, 0}, {2, 1}, {2, 2}},
		{{0, 0}, {1, 0}, {2, 0}},
		{{0, 1}, {1, 1}, {2, 1}},
		{{0, 2}, {1, 2}, {2, 2}},
		{{0, 0}, {1, 1}, {2, 2}},
		{{0, 2}, {1, 1}, {2, 0}},
	}

	for _, line := range lines {
		a, b, c := line[0], line[1], line[2]
		if board[a[0]][a[1]] != "" &&
			board[a[0]][a[1]] == board[b[0]][b[1]] &&
			board[b[0]][b[1]] == board[c[0]][c[1]] {
			return board[a[0]][a[1]]
		}
	}
	return ""
}

func checkDraw(board [3][3]string) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == "" {
				return false
			}
		}
	}
	return true
}
