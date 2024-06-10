package repositories

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"ingenhouzs.com/chesshouzs/go-game/helpers"
	"ingenhouzs.com/chesshouzs/go-game/helpers/errs"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

func (r *Repository) GetUnderMatchmakingPlayers(params models.PoolParams) ([]models.PlayerPool, error) {
	var data []models.PlayerPool
	key := helpers.GetPoolKey(params)

	timeout := helpers.GetTimeoutThreshold("DATABASE_QUERY_TIMEOUT")
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	redisClient := r.redis.ZRange(ctx, key, 0, -1)
	pool, err := redisClient.Result()
	if err != nil {
		return data, err
	}

	for _, player := range pool {
		var playerData models.PlayerPool
		if err := json.Unmarshal([]byte(player), &playerData); err != nil {
			return data, err
		}
		data = append(data, playerData)
	}

	return data, nil
}

func (r *Repository) InsertPlayerIntoPool(params models.PlayerPoolParams) error {
	key := helpers.GetPoolKey(params.PoolParams)

	timeout := helpers.GetTimeoutThreshold("DATABASE_QUERY_TIMEOUT")
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	data, err := json.Marshal(models.PlayerPool{
		User: models.User{
			ID:        params.User.ID,
			EloPoints: params.User.EloPoints,
		},
		JoinTime: time.Now(),
	})
	if err != nil {
		return err
	}

	result := r.redis.ZAdd(ctx, key, redis.Z{
		Score:  float64(params.EloPoints),
		Member: data,
	})
	if err := result.Err(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeletePlayerFromPool(params models.PlayerPoolParams) error {
	key := helpers.GetPoolKey(params.PoolParams)

	timeout := helpers.GetTimeoutThreshold("DATABASE_QUERY_TIMEOUT")
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	result := r.redis.ZRem(ctx, key, redis.Z{
		Score:  float64(params.Player.User.EloPoints),
		Member: params.Player.JoinTime,
	})
	if err := result.Err(); err != nil {
		return err
	}

	if result.Val() <= 0 {
		return errs.ERR_REDIS_DATA_NOT_FOUND
	}

	return nil
}

func (r *Repository) InsertMoveCacheIdentifier(params models.MoveCache) error {
	key := helpers.GetGameMoveCacheKey(params)

	timeout := helpers.GetTimeoutThreshold("DATABASE_QUERY_TIMEOUT")
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	result := r.redis.HMSet(ctx, key, map[string]interface{}{
		"move": "",
		"turn": true,
	})

	if err := result.Err(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) InsertGameData(params models.InsertGameParams) error {

	query := `
		INSERT INTO game_active 
			(
				white_player_id, 
				black_player_id, 
				game_type_variant_id, 
				moves_cache_ref
			)
			VALUES 
			(
				?, 
				?, 
				?,
				?
			)
	`
	var empty []interface{}
	return r.postgres.Raw(
		query,
		params.WhitePlayerID.String(),
		params.BlackPlayerID.String(),
		params.GameTypeVariantID.String(),
		params.MovesCacheRef.String(),
	).Scan(&empty).Error
}
