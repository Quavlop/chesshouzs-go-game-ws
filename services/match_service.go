package services

import (
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"ingenhouzs.com/chesshouzs/go-game/constants"
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

	// check if player is already in game
	if s.wsConnections.IsClientInRoom(constants.WS_ROOM_TYPE_GAME, user.ID.String()) {
		return result, errs.ERR_PLAYER_IN_GAME
	}

	// check if player is registered on pool
	playerPoolData, err := s.repository.GetPlayerPoolData(models.PlayerPoolParams{
		PoolParams: models.PoolParams{
			Type:        params.Type,
			TimeControl: params.TimeControl,
		},
		User: user,
	})
	if err != nil && err != errs.ERR_REDIS_DATA_NOT_FOUND {
		return result, err
	}

	if len(playerPoolData) > 0 {
		return result, errs.ERR_PLAYER_IN_POOL
	}

	eloBounds := s.BaseService.HttpService.CalculateEloBounds(user)

	eligibleOpponents, err := s.FilterEligibleOpponent(client, models.FilterEligibleOpponentParams{
		Filter: models.PoolParams{
			Type:        params.Type,
			TimeControl: params.TimeControl,
			UpperBound:  eloBounds.Upper,
			LowerBound:  eloBounds.Lower,
		},
		Client: models.PlayerPool{
			User: user,
		},
	})
	if err != nil && err != errs.ERR_NO_AVAILABLE_PLAYERS {
		helpers.LogErrorCallStack(*client.Context, err)
		return result, err
	}

	joinTime := time.Now()
	poolParams := models.PlayerPoolParams{
		PoolParams: models.PoolParams{
			Type:        params.Type,
			TimeControl: params.TimeControl,
		},
		User: user,
	}

	if err == errs.ERR_NO_AVAILABLE_PLAYERS {
		err = s.repository.WithRedisTrx((*(client.Context)).Request().Context(), []string{helpers.GetPlayerPoolCloneKey(poolParams), helpers.GetPoolKey(poolParams.PoolParams)}, func(pipe redis.Pipeliner) error {
			err = s.repository.InsertPlayerIntoPool(poolParams, joinTime)
			if err != nil {
				return err
			}

			err = s.repository.InsertPlayerOnPoolDataToRedis(poolParams, joinTime)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return result, err
		}

		result = models.HandleMatchmakingResponse{
			ID: user.ID.String(),
			Opponent: models.PlayerMatchmakingResponse{
				JoinTime: joinTime.Format(time.RFC3339),
				User:     user,
			},
		}

		return result, nil
	}

	// if found match then remove the enemy from pool and insert both into game data
	opponent := eligibleOpponents.Player
	result.Opponent = models.PlayerMatchmakingResponse{
		JoinTime: joinTime.Format(time.RFC3339),
		User:     user,
	}

	redisPoolKey := helpers.GetPoolKey(models.PoolParams{
		Type:        params.Type,
		TimeControl: params.TimeControl,
	})

	moveCacheID := uuid.New()

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

		// delete redis hash key for player data clone
		err = s.repository.DeletePlayerOnPoolDataToRedis(poolParams, joinTime)
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
	gameID := uuid.New()
	err = s.repository.InsertGameData(models.InsertGameParams{
		ID:                gameID,
		WhitePlayerID:     client.User.ID,
		BlackPlayerID:     opponent.User.ID,
		GameTypeVariantID: gameTypeVariantID,
		MovesCacheRef:     moveCacheID,
		CreatedAt:         time.Now().Format(time.RFC3339),
	})
	if err != nil {
		return result, err
	}

	s.wsConnections.EmitOneOnOne(models.WebSocketChannel{
		Source:       client.User.ID.String(),
		TargetClient: opponent.User.ID.String(),
		Event:        constants.WS_EVENT_EMIT_START_GAME,
		Data: models.GameData{
			ID: gameID.String(),
		},
	})

	room := s.wsConnections.CreateRoom(&models.GameRoom{
		Name: gameID.String(),
		Type: constants.WS_ROOM_TYPE_GAME,
	})

	room.AddClient(client.User.ID.String())
	room.AddClient(opponent.User.ID.String())

	result.ID = gameID.String()

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

	// playerPool, err = s.FilterOutOpponents(client, playerPool)
	// if err != nil {
	// 	helpers.LogErrorCallStack(*client.Context, err)
	// 	return result, err
	// }

	// if len(playerPool) <= 0 {
	// 	err = errs.ERR_NO_AVAILABLE_PLAYERS
	// 	helpers.LogErrorCallStack(*client.Context, err)
	// 	return result, err
	// }

	playerPool, err = s.SortPlayerPool(client, playerPool)
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
