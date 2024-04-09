package models

import "github.com/gorilla/websocket"

type WebSocketChannel struct {
	Source *websocket.Conn
	Target *websocket.Conn
	Data   interface{}
}

type WebSocketResponse struct {
	Source string `json:"source"`
	Event  string `json:"event"`
	Data   string `json:"data"`
}
