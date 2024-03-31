package services

import "ingenhouzs.com/chesshouzs/go-game/interfaces"

type httpService struct {
	repository interfaces.Repository
}

type webSocketService struct {
	repository interfaces.Repository
}

func NewHttpService(repository interfaces.Repository) interfaces.HttpService {
	return &httpService{repository}
}

func NewWebSocketService(repository interfaces.Repository) interfaces.WebsocketService {
	return &webSocketService{repository}
}
