package repositories

import (
	"context"
	"encoding/json"
	"time"

	"ingenhouzs.com/chesshouzs/go-game/helpers"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

func (r *Repository) GetUnderMatchmakingPlayers(params models.PoolParams) ([]models.PlayerPool, error) {
	var data []models.PlayerPool
	key := helpers.GetPoolKey(params)

	timeout := helpers.GetTimeoutThreshold("DATABASE_QUERY_TIMEOUT")
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	redisClient := r.redis.LRange(ctx, key, 0, -1)
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
