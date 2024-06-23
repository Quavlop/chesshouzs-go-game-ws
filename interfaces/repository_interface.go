package interfaces

import (
	"context"
	"time"

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
	InsertPlayerIntoPool(params models.PlayerPoolParams, joinTime time.Time, pipe redis.Pipeliner) error
	DeletePlayerFromPool(params models.PlayerPoolParams, pipe redis.Pipeliner) error
	InsertMoveCacheIdentifier(params models.MoveCache, pipe redis.Pipeliner) error
	DeleteMoveCacheIdentifier(params models.MoveCache, pipe redis.Pipeliner) error
	InsertGameData(params models.GameActiveData) error
	InsertPlayerOnPoolDataToRedis(params models.PlayerPoolParams, joinTime time.Time, pipe redis.Pipeliner) error
	DeletePlayerOnPoolDataToRedis(params models.PlayerPoolParams, joinTime time.Time, pipe redis.Pipeliner) error
	GetPlayerPoolData(params models.PlayerPoolParams) (map[string]string, error)
	GetPlayerCurrentGameState(token string) (models.GameActiveData, error)
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
