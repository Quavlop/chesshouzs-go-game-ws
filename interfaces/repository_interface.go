package interfaces

import "ingenhouzs.com/chesshouzs/go-game/models"

type Repository interface {
	MatchRepository
	UserRepository
}

type MatchRepository interface {
	GetUnderMatchmakingPlayers(params models.PoolParams) ([]models.PlayerPool, error)
}

type UserRepository interface {
	GetUserDataByID(id string) (models.User, error)
}
