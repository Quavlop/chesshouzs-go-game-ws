package services

import "ingenhouzs.com/chesshouzs/go-game/interfaces"

type service struct {
	repository interfaces.Repository
}

func NewService(repository interfaces.Repository) interfaces.Service {
	return &service{repository}
}
