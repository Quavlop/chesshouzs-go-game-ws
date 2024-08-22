package repositories

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
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

	redisClient := r.redis.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: strconv.Itoa(int(params.LowerBound)),
		Max: strconv.Itoa(int(params.UpperBound)),
	})
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

func (r *Repository) InsertPlayerIntoPool(params models.PlayerPoolParams, joinTime time.Time, pipe redis.Pipeliner) error {
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

	var result *redis.IntCmd
	if pipe != nil {
		result = pipe.ZAdd(ctx, key, redis.Z{
			Score:  float64(params.User.EloPoints),
			Member: data,
		})
	} else {
		result = r.redis.ZAdd(ctx, key, redis.Z{
			Score:  float64(params.User.EloPoints),
			Member: data,
		})
	}

	if err := result.Err(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeletePlayerFromPool(params models.PlayerPoolParams, pipe redis.Pipeliner) error {
	key := helpers.GetPoolKey(params.PoolParams)
	helpers.WriteOutLog("GET POOL_REDIS : " + key)

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

	var result *redis.IntCmd
	if pipe != nil {
		result = pipe.ZRem(ctx, key, string(data))
	} else {
		result = r.redis.ZRem(ctx, key, string(data))
	}

	if err := result.Err(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) InsertMoveCacheIdentifier(params models.MoveCache, pipe redis.Pipeliner) error {
	key := helpers.GetGameMoveCacheKey(params)
	helpers.WriteOutLog("INSERT MOVE_CACHE_KEY : " + key)

	timeout := helpers.GetTimeoutThreshold("DATABASE_QUERY_TIMEOUT")
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	if params.Move == "" {
		params.Move = helpers.GameNotationBuilder(14)
	}

	var result *redis.BoolCmd
	if pipe != nil {
		result = pipe.HMSet(ctx, key, map[string]interface{}{
			"move": params.Move,
			"turn": params.Turn,
		})
	} else {
		result = r.redis.HMSet(ctx, key, map[string]interface{}{
			"move": params.Move,
			"turn": params.Turn,
		})
	}

	if err := result.Err(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) InsertGameData(params models.GameActiveData) error {

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
		params.StartTime,
	).Scan(&empty).Error
}

func (r *Repository) InsertPlayerOnPoolDataToRedis(params models.PlayerPoolParams, joinTime time.Time, pipe redis.Pipeliner) error {
	key := helpers.GetPlayerPoolCloneKey(params)

	timeout := helpers.GetTimeoutThreshold("DATABASE_QUERY_TIMEOUT")
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	var result *redis.BoolCmd
	if pipe != nil {
		result = pipe.HMSet(ctx, key, "game-type", params.Type, "game-time-control", params.TimeControl, "game-join-time", joinTime)
	} else {
		result = r.redis.HMSet(ctx, key, "game-type", params.Type, "game-time-control", params.TimeControl, "game-join-time", joinTime)
	}
	if err := result.Err(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeletePlayerOnPoolDataToRedis(params models.PlayerPoolParams, joinTime time.Time, pipe redis.Pipeliner) error {
	key := helpers.GetPlayerPoolCloneKey(params)

	timeout := helpers.GetTimeoutThreshold("DATABASE_QUERY_TIMEOUT")
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	var result *redis.IntCmd
	if pipe != nil {
		result = pipe.Del(ctx, key)
	} else {
		result = r.redis.Del(ctx, key)
	}

	if err := result.Err(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetPlayerPoolData(params models.PlayerPoolParams) (map[string]string, error) {
	key := helpers.GetPlayerPoolCloneKey(params)
	helpers.WriteOutLog("GET PLAYER_POOL_CLONE : " + key)

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

func (r *Repository) GetPlayerCurrentGameState(token string) (models.GameActiveData, error) {
	var data models.GameActiveData

	db := r.postgres.Table("game_active ga").
		Select("*")

	if token != "" {
		db = db.Where("(white_player_id = ? OR black_player_id = ?)", token, token)
	}

	result := db.Take(&data)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return data, result.Error
	}

	if result.Error == gorm.ErrRecordNotFound {
		return data, errs.ERR_ACTIVE_GAME_NOT_FOUND
	}

	return data, nil
}

func (r *Repository) DeleteMoveCacheIdentifier(params models.MoveCache, pipe redis.Pipeliner) error {
	key := helpers.GetGameMoveCacheKey(params)
	helpers.WriteOutLog("REMOVE (INVALIDATION) MOVE_CACHE_KEY : " + key)

	timeout := helpers.GetTimeoutThreshold("DATABASE_QUERY_TIMEOUT")
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	var result *redis.IntCmd
	if pipe != nil {
		result = pipe.Del(ctx, key)
	} else {
		result = r.redis.Del(ctx, key)
	}

	if err := result.Err(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) InsertMatchSkillCount(params models.InitMatchSkillStats, pipe redis.Pipeliner) error {
	key := helpers.GetPlayerMatchSkillState(params)

	timeout := helpers.GetTimeoutThreshold("DATABASE_QUERY_TIMEOUT")
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	var args []interface{}
	if len(params.GameSkillMap) > 0 {
		for key, val := range params.GameSkillMap {
			args = append(args, key, val)
		}
	} else {
		for _, skill := range params.GameSkills {
			args = append(args, skill.ID.String(), skill.UsageCount)
		}
	}

	var result *redis.BoolCmd
	if pipe != nil {
		result = pipe.HMSet(ctx, key, args...)
	} else {
		result = r.redis.HMSet(ctx, key, args...)
	}
	if err := result.Err(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteMatchSkillCount(params models.InitMatchSkillStats, pipe redis.Pipeliner) error {
	key := helpers.GetPlayerMatchSkillState(params)

	timeout := helpers.GetTimeoutThreshold("DATABASE_QUERY_TIMEOUT")
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	var result *redis.IntCmd
	if pipe != nil {
		result = pipe.Del(ctx, key)
	} else {
		result = r.redis.Del(ctx, key)
	}
	if err := result.Err(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetPlayerSkillCountUsageData(params models.InitMatchSkillStats) (map[string]int, error) {
	var intResult map[string]int
	key := helpers.GetPlayerMatchSkillState(params)

	timeout := helpers.GetTimeoutThreshold("DATABASE_QUERY_TIMEOUT")
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	result, err := r.redis.HGetAll(ctx, key).Result()
	if err != nil {
		return intResult, err
	}

	if len(result) <= 0 {
		return intResult, errs.ERR_REDIS_DATA_NOT_FOUND
	}

	// Convert map[string]string to map[string]int
	intResult = make(map[string]int, len(result))
	for field, value := range result {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return intResult, err
		}
		intResult[field] = intValue
	}

	return intResult, nil
}

func (r *Repository) GetPlayerState(params models.PlayerState) (models.PlayerState, error) {
	query := `
		SELECT player_id, game_id, buff_state, debuff_state
		FROM chesshouzs.player_game_states
		WHERE player_id = ? AND game_id = ?
	`

	var playerState models.PlayerState

	err := r.cassandra.
		Query(query, params.PlayerID, params.GameID).
		Scan(&playerState.PlayerID, &playerState.GameID, &playerState.BuffState, &playerState.DebuffState)
	if err != nil {
		return models.PlayerState{}, err
	}

	return playerState, nil
}

func (r *Repository) InsertPlayerState(params models.PlayerState) error {

	query := `
		INSERT INTO chesshouzs.player_game_states (
			player_id, 
			game_id, 
			buff_state, 
			debuff_state
		) VALUES (?, ?, ?, ?)
	`

	err := r.cassandra.
		Query(query, params.PlayerID, params.GameID, params.BuffState, params.DebuffState).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdatePlayerState(params models.PlayerState) error {

	query := `
		UPDATE chesshouzs.player_game_states
		SET buff_state = ?, debuff_state = ?
		WHERE player_id = ? AND game_id = ?
	`

	err := r.cassandra.
		Query(query, params.BuffState, params.DebuffState, params.PlayerID, params.GameID).
		Exec()
	if err != nil {
		return err
	}

	return nil
}
