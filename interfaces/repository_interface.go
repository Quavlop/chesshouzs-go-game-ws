package interfaces

import "ingenhouzs.com/chesshouzs/go-game/models"

type Repository interface {
	GameRepository
}

type GameRepository interface {
	GetUnderMatchmakingPlayers(params models.PoolParams) ([]models.PlayerPool, error)
}
