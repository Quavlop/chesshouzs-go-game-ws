package interfaces

import "ingenhouzs.com/chesshouzs/go-game/models"

type HttpService interface {
	GameService
}

type GameService interface {
	IsValidGameType(params models.GameTypeVariant) (bool, error)
	CalculateEloBounds(params models.User) models.EloBounds
}
