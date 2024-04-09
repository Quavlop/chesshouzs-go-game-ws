package models

import (
	"github.com/gorilla/websocket"
)

type WebSocketClientConnection struct {
	Connection *websocket.Conn
}

type WebSocketChannel struct {
	source *websocket.Conn
	target *websocket.Conn
	Data   interface{}
}

func (wsChan WebSocketChannel) GetSource() *websocket.Conn {
	return wsChan.source
}

func (wsChan WebSocketChannel) GetTarget() *websocket.Conn {
	return wsChan.target
}

type WebSocketResponse struct {
	Status string      `json:"status,omitempty"`
	Source string      `json:"source,omitempty"`
	Target string      `json:"target,omitempty"`
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
