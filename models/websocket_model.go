package models

import "github.com/gorilla/websocket"

type WebSocketChannel struct {
	Source *websocket.Conn
	Target *websocket.Conn
	Data   interface{}
}

type WebSocketResponse struct {
	Source string
	Event  string
	Data   string
}
