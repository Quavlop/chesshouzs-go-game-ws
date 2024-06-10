package services

import (
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"ingenhouzs.com/chesshouzs/go-game/helpers"
	"ingenhouzs.com/chesshouzs/go-game/helpers/errs"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

// Websocket services
func (s *webSocketService) HandleMatchmaking(client models.WebSocketClientData, params models.HandleMatchmakingParams) (models.HandleMatchmakingResponse, error) {
	var result models.HandleMatchmakingResponse

	user := client.User

	timeControlParsedStr := strings.Split(params.TimeControl, "-")
	if len(timeControlParsedStr) != 2 {
		return result, errs.ERR_ERROR_TIME_CONTROL_PARSE
	}

	timeControlSecond, err := strconv.Atoi(timeControlParsedStr[0])
	if err != nil {
		return result, err
	}
	timeControlIncrement, err := strconv.Atoi(timeControlParsedStr[1])
	if err != nil {
		return result, err
	}

	isGameParamValid, err := s.BaseService.HttpService.IsValidGameType(models.GameTypeVariant{
		Name:      params.Type,
		Duration:  int32(timeControlSecond),
		Increment: int32(timeControlIncrement),
	})
	if err != nil {
		return result, err
	}

	if !isGameParamValid {
		return result, errs.ERR_INVALID_GAME_TYPE
	}

	eligibleOpponents, err := s.BaseService.WebSocketService.FilterEligibleOpponent(client, models.FilterEligibleOpponentParams{
		Filter: models.PoolParams{
			Type:        params.Type,
			TimeControl: params.TimeControl,
		},
		Client: models.PlayerPool{
			User: user,
		},
	})
	if err != nil && err != errs.ERR_NO_AVAILABLE_PLAYERS {
		helpers.LogErrorCallStack(*client.Context, err)
		return result, err
	}

	if err == errs.ERR_NO_AVAILABLE_PLAYERS {
		err = s.repository.InsertPlayerIntoPool(models.PlayerPoolParams{
			PoolParams: models.PoolParams{
				Type:        params.Type,
				TimeControl: params.TimeControl,
			},
			User: user,
		})
		if err != nil {
			return result, err
		}
		return result, nil
	}

	// if found match then remove the enemy from pool and insert both into game data
	opponent := eligibleOpponents.Player
	result.Opponent = opponent

	redisPoolKey := helpers.GetPoolKey(models.PoolParams{
		Type:        params.Type,
		TimeControl: params.TimeControl,
	})

	moveCacheID, _ := uuid.Parse(uuid.NewString())

	err = s.repository.WithRedisTrx((*(client.Context)).Request().Context(), []string{redisPoolKey}, func(pipe redis.Pipeliner) error {
		err = s.repository.DeletePlayerFromPool(models.PlayerPoolParams{
			PoolParams: models.PoolParams{
				Type:        params.Type,
				TimeControl: params.TimeControl,
			},
			Player: opponent,
		})
		if err != nil {
			return err
		}
		// insert game cache move ref to redis
		err = s.repository.InsertMoveCacheIdentifier(models.MoveCache{
			ID: moveCacheID,
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return result, err
	}

	gameTypeVariant, err := s.repository.GetGameTypeVariant(models.GameTypeVariant{
		Duration:  int32(timeControlSecond),
		Increment: int32(timeControlIncrement),
	})
	if err != nil {
		return result, err
	}

	var gameTypeVariantID uuid.UUID
	if len(gameTypeVariant) > 0 {
		gameTypeVariantID = gameTypeVariant[0].ID
	}

	// insert to mysql
	err = s.repository.InsertGameData(models.InsertGameParams{
		WhitePlayerID:     client.User.ID,
		BlackPlayerID:     opponent.User.ID,
		GameTypeVariantID: gameTypeVariantID,
		MovesCacheRef:     moveCacheID,
	})
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *webSocketService) FilterEligibleOpponent(client models.WebSocketClientData, params models.FilterEligibleOpponentParams) (models.FilterEligibleOpponentResponse, error) {
	var result models.FilterEligibleOpponentResponse

	playerPool, err := s.repository.GetUnderMatchmakingPlayers(params.Filter)
	if err != nil {
		helpers.LogErrorCallStack(*client.Context, err)
		return result, err
	}

	if len(playerPool) <= 0 {
		err = errs.ERR_NO_AVAILABLE_PLAYERS
		helpers.LogErrorCallStack(*client.Context, err)
		return result, err
	}

	playerPool, err = s.BaseService.WebSocketService.FilterOutOpponents(client, playerPool)
	if err != nil {
		helpers.LogErrorCallStack(*client.Context, err)
		return result, err
	}

	if len(playerPool) <= 0 {
		err = errs.ERR_NO_AVAILABLE_PLAYERS
		helpers.LogErrorCallStack(*client.Context, err)
		return result, err
	}

	playerPool, err = s.BaseService.WebSocketService.SortPlayerPool(client, playerPool)
	if err != nil {
		helpers.LogErrorCallStack(*client.Context, err)
		return result, err
	}

	result.Player = playerPool[0]

	return result, nil
}

func (s *webSocketService) SortPlayerPool(client models.WebSocketClientData, pool []models.PlayerPool) ([]models.PlayerPool, error) {
	for i := 0; i < len(pool)-1; i++ {
		minIdx := i
		for idx := i + 1; idx < len(pool); idx++ {
			if s.BaseService.WebSocketService.PlayerSortFilter(pool[minIdx], pool[idx]) {
				minIdx = idx
			}
		}
		temp := pool[i]
		pool[i] = pool[minIdx]
		pool[minIdx] = temp
	}

	return pool, nil
}

func (s *webSocketService) FilterOutOpponents(client models.WebSocketClientData, pool []models.PlayerPool) ([]models.PlayerPool, error) {
	var result []models.PlayerPool

	// get player data first
	player := models.PlayerPool{
		User: client.User,
	}

	for _, opponent := range pool {
		if s.BaseService.WebSocketService.IsMatchmakingEligible(player, opponent) {
			result = append(result, opponent)
		}
	}
	return result, nil
}

func (s *webSocketService) IsMatchmakingEligible(player models.PlayerPool, opponent models.PlayerPool) bool {
	threshold, err := strconv.Atoi(os.Getenv("ELO_GAP_THRESHOLD"))
	if err != nil {
		return math.Abs(float64(player.User.EloPoints-opponent.User.EloPoints)) <= 150
	}
	return math.Abs(float64(player.User.EloPoints-opponent.User.EloPoints)) <= float64(threshold)
}

func (s *webSocketService) PlayerSortFilter(playerOne models.PlayerPool, playerTwo models.PlayerPool) bool {
	if playerOne.User.EloPoints != playerTwo.User.EloPoints {
		return playerOne.User.EloPoints > playerTwo.User.EloPoints
	}
	return playerOne.JoinTime.After(playerTwo.JoinTime)
}
