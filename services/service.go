package services

import (
	"ingenhouzs.com/chesshouzs/go-game/config/websocket"
	"ingenhouzs.com/chesshouzs/go-game/interfaces"
)

type BaseService struct {
	WebSocketService interfaces.WebsocketService
	HttpService      interfaces.HttpService
}
type httpService struct {
	repository  interfaces.Repository
	baseService BaseService
}

type webSocketService struct {
	repository    interfaces.Repository
	wsConnections *websocket.Connections
	baseService   BaseService
}

type gameRoomService struct {
	room interfaces.WebSocketRoom
}

func NewBaseService(webSocketService interfaces.WebsocketService, httpService interfaces.HttpService) *BaseService {
	return &BaseService{WebSocketService: webSocketService, HttpService: httpService}
}

func NewHttpService(repository interfaces.Repository, baseService BaseService) interfaces.HttpService {
	return &httpService{repository, baseService}
}

func NewWebSocketService(repository interfaces.Repository, wsConnections *websocket.Connections, baseService BaseService) interfaces.WebsocketService {
	return &webSocketService{repository, wsConnections, baseService}
}
