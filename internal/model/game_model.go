package model

import "time"

type Player struct {
	ClientID string `json:"clientID"`
	Symbol   string `json:"symbol"`
}

type Game struct {
	RoomID    string       `json:"roomID"`
	Players   []Player     `json:"players"`
	Board     [3][3]string `json:"board"`
	Winner    string       `json:"winner"`
	Turn      string       `json:"turn"`
	IsOver    bool         `json:"isOver"`
	IsDraw    bool         `json:"isDraw"`
	CreatedAt time.Time    `json:"createdAt"`
}

type Move struct {
	Player string `json:"player"`
	RoomID string `json:"roomId"`
	Row    int    `json:"row"`
	Col    int    `json:"col"`
}

type StartGamePayload struct {
	RoomID  string   `json:"roomId"`
	Players []Player `json:"players"`
	Turn    string   `json:"turn"`
}
