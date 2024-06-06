package models

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

type WebSocketClientConnection struct {
	Connection *websocket.Conn
	Token      string
}

type WebSocketClientData struct {
	Connection *websocket.Conn
	Token      string
	Event      string
	Context    *echo.Context
	User       User
	Data       interface{}
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
