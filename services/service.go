package services

import (
	"ingenhouzs.com/chesshouzs/go-game/config/websocket"
	"ingenhouzs.com/chesshouzs/go-game/interfaces"
)

type httpService struct {
	repository interfaces.Repository
}

type webSocketService struct {
	repository    interfaces.Repository
	wsConnections *websocket.Connections
}

type gameRoomService struct {
	room interfaces.WebSocketRoom
}

func NewHttpService(repository interfaces.Repository) interfaces.HttpService {
	return &httpService{repository}
}

func NewWebSocketService(repository interfaces.Repository, wsConnections *websocket.Connections) interfaces.WebsocketService {
	return &webSocketService{repository, wsConnections}
}
