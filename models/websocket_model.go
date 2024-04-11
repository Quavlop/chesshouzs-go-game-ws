package models

import (
	"github.com/gorilla/websocket"
)

type WebSocketClientConnection struct {
	Connection *websocket.Conn
	Token      string
}

type WebSocketClientMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type WebSocketChannel struct {
	Source       string
	TargetClient string
	TargetRoom   string
	Event        string
	Data         interface{}
}

type WebSocketResponse struct {
	Status string      `json:"status,omitempty"`
	Event  string      `json:"event,omitempty"`
	Data   interface{} `json:"data"`
}

type GameRoom struct {
	id   string
	Name string // METADATA
	Type string // METADATA
	// etc ...

	// key : user's session token
	// [BETA] value : boolean value to show if user is in the room ()
	clients map[string]bool
}