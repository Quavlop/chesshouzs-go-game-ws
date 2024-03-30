package services

import "ingenhouzs.com/chesshouzs/go-game/interfaces"

type service struct {
	repository interfaces.Repository
}

func NewHttpService(repository interfaces.Repository) interfaces.HttpService {
	return &service{repository}
}

func NewWebSocketService(repository interfaces.Repository) interfaces.WebsocketService {
	return &service{repository}
}
