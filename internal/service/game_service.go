package service

import (
	"errors"
	"tic-tac-toe-game/internal/model"
	"tic-tac-toe-game/internal/utility"
)

type GameService struct{}

func NewGameService() *GameService {
	return &GameService{}
}

var GameRooms = []model.Game{}

func (s *GameService) GetAllRooms() []model.Game {
	return GameRooms
}

func (s *GameService) CreateGameRoom(gameRoom model.Game) model.Game {
	var RoomID string
	for {
		RoomID = utility.GenerateRoomID()
		if s.FindGameRoom(RoomID).RoomID == "" {
			break
		}
	}
	gameRoom.RoomID = RoomID
	GameRooms = append(GameRooms, gameRoom)
	return gameRoom
}

func (s *GameService) FindGameRoom(roomID string) model.Game {
	for _, room := range GameRooms {
		if room.RoomID == roomID {
			return room
		}
	}
	return model.Game{}
}

func (s *GameService) DeleteGameRoom(roomID string) {
	for i, room := range GameRooms {
		if room.RoomID == roomID {
			GameRooms = append(GameRooms[:i], GameRooms[i+1:]...)
		}
	}
}

func (s *GameService) JoinGameRoom(roomID string, ClientID string) error {
	for i, room := range GameRooms {
		if room.RoomID == roomID {
			if len(room.Players) >= 2 {
				return errors.New("room is full")
			}
			if len(room.Players) == 0 {
				// ผู้เล่นแรกได้ X
				room.Players = append(room.Players, model.Player{ClientID: ClientID, Symbol: "X"})
				GameRooms[i] = room
			} else if len(room.Players) == 1 {
				// ผู้เล่นที่สองได้ O
				room.Players = append(room.Players, model.Player{ClientID: ClientID, Symbol: "O"})
				GameRooms[i] = room
			}

		}
	}
	return nil
}

func (s *GameService) GameReady(roomID string) bool {
	for _, room := range GameRooms {
		if room.RoomID == roomID && len(room.Players) == 2 {
			return true
		}
	}
	return false
}

func (s *GameService) LeaveGameRoom(roomID string, clientID string) error {
	for i, room := range GameRooms {
		if room.RoomID == roomID {
			// ลบ player ออกจาก room
			newPlayers := []model.Player{}
			for _, p := range room.Players {
				if p.ClientID != clientID {
					newPlayers = append(newPlayers, p)
				}
			}

			room.Players = newPlayers
			GameRooms[i] = room

			// ถ้าไม่มีผู้เล่นแล้ว ให้ลบห้อง
			if len(room.Players) == 0 {
				s.DeleteGameRoom(roomID)
			}
			return nil
		}
	}
	return errors.New("room not found")
}

func (s *GameService) UpdateGameRoom(updatedGame model.Game) error {
	for i, room := range GameRooms {
		if room.RoomID == updatedGame.RoomID {
			GameRooms[i] = updatedGame
			return nil
		}
	}
	return errors.New("room not found")
}
