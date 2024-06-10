package interfaces

import (
	"context"

	"github.com/redis/go-redis/v9"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

type Repository interface {
	MatchRepository
	UserRepository
	GameRepository

	Transaction
}

type MatchRepository interface {
	GetUnderMatchmakingPlayers(params models.PoolParams) ([]models.PlayerPool, error)
	InsertPlayerIntoPool(params models.PlayerPoolParams) error
	DeletePlayerFromPool(params models.PlayerPoolParams) error
	InsertMoveCacheIdentifier(params models.MoveCache) error
	InsertGameData(params models.InsertGameParams) error
}

type UserRepository interface {
	GetUserDataByID(id string) (models.User, error)
}

type GameRepository interface {
	GetGameTypeVariant(params models.GameTypeVariant) ([]models.GameTypeVariant, error)
}

type Transaction interface {
	WithRedisTrx(ctx context.Context, keys []string, fn func(pipe redis.Pipeliner) error) error
}
