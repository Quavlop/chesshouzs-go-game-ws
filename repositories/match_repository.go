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

func (r *Repository) InsertPlayerIntoPool(params models.PlayerPoolParams, joinTime time.Time) error {
	key := helpers.GetPoolKey(params.PoolParams)

	timeout := helpers.GetTimeoutThreshold("DATABASE_QUERY_TIMEOUT")
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	data, err := json.Marshal(models.PlayerPool{
		User: models.User{
			ID:        params.User.ID,
			EloPoints: params.User.EloPoints,
		},
		JoinTime: joinTime,
	})
	if err != nil {
		return err
	}

	result := r.redis.ZAdd(ctx, key, redis.Z{
		Score:  float64(params.User.EloPoints),
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

	data, err := json.Marshal(models.PlayerPool{
		User: models.User{
			ID:        params.Player.User.ID,
			EloPoints: params.Player.User.EloPoints,
		},
		JoinTime: params.Player.JoinTime,
	})
	if err != nil {
		return err
	}

	result := r.redis.ZRem(ctx, key, data)
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
	helpers.WriteOutLog("MOVE_CACHE_KEY : " + key)

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
				id,
				white_player_id, 
				black_player_id, 
				game_type_variant_id, 
				moves_cache_ref, 
				start_time
			)
			VALUES 
			(
				?,
				?, 
				?, 
				?,
				?,
				?
			)
	`
	var empty []interface{}
	return r.postgres.Raw(
		query,
		params.ID.String(),
		params.WhitePlayerID.String(),
		params.BlackPlayerID.String(),
		params.GameTypeVariantID.String(),
		params.MovesCacheRef.String(),
		params.CreatedAt,
	).Scan(&empty).Error
}

func (r *Repository) InsertPlayerOnPoolDataToRedis(params models.PlayerPoolParams, joinTime time.Time) error {
	key := helpers.GetPlayerPoolCloneKey(params)

	timeout := helpers.GetTimeoutThreshold("DATABASE_QUERY_TIMEOUT")
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	data, err := json.Marshal(models.PlayerPool{
		User: models.User{
			ID:        params.User.ID,
			EloPoints: params.User.EloPoints,
		},
		JoinTime: joinTime,
	})
	if err != nil {
		return err
	}

	result := r.redis.HSet(ctx, key, "data", data)
	if err := result.Err(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeletePlayerOnPoolDataToRedis(params models.PlayerPoolParams, joinTime time.Time) error {
	key := helpers.GetPlayerPoolCloneKey(params)

	timeout := helpers.GetTimeoutThreshold("DATABASE_QUERY_TIMEOUT")
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	result := r.redis.Del(ctx, key)
	if err := result.Err(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetPlayerPoolData(params models.PlayerPoolParams) (map[string]string, error) {
	key := helpers.GetPlayerPoolCloneKey(params)

	timeout := helpers.GetTimeoutThreshold("DATABASE_QUERY_TIMEOUT")
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	result, err := r.redis.HGetAll(ctx, key).Result()
	if err != nil {
		return result, err
	}

	if len(result) <= 0 {
		return result, errs.ERR_REDIS_DATA_NOT_FOUND
	}

	return result, nil
}
