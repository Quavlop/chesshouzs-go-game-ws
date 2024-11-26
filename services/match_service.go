package services

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/redis/go-redis/v9"
	"ingenhouzs.com/chesshouzs/go-game/constants"
	"ingenhouzs.com/chesshouzs/go-game/helpers"
	"ingenhouzs.com/chesshouzs/go-game/helpers/errs"
	"ingenhouzs.com/chesshouzs/go-game/models"
	"ingenhouzs.com/chesshouzs/go-game/services/rpc/pb"
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
	// TODO add validation for ingame database psql
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
			err = s.repository.InsertPlayerIntoPool(poolParams, joinTime, pipe)
			if err != nil {
				return err
			}

			err = s.repository.InsertPlayerOnPoolDataToRedis(poolParams, joinTime, pipe)
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

	skills, err := s.repository.GetGameSkills(models.GameSkill{})
	if err != nil {
		helpers.LogErrorCallStack(*client.Context, err)
		return result, err
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

	moveCacheKey := helpers.GetGameMoveCacheKey(models.MoveCache{
		ID: moveCacheID,
	})

	poolCloneKey := helpers.GetPlayerPoolCloneKey(poolParams)

	gameID := uuid.New()

	defaultBuffSkillState := make(map[string]models.SkillState)
	defaultDebuffSkillState := make(map[string]models.SkillState)

	for _, skill := range skills {
		if skill.Type == constants.SKILL_TYPE_BUFF && !skill.Permanent {
			defaultBuffSkillState[skill.Name] = models.SkillState{
				DurationLeft: 0,
				List:         []models.SkillStatus{},
			}
		} else if skill.Type == constants.SKILL_TYPE_DEBUFF && !skill.Permanent {
			defaultDebuffSkillState[skill.Name] = models.SkillState{
				DurationLeft: 0,
				List:         []models.SkillStatus{},
			}
		}
	}

	// insert player state
	state := models.PlayerState{
		PlayerID:    user.ID.String(),
		GameID:      gameID.String(),
		BuffState:   defaultBuffSkillState,
		DebuffState: defaultDebuffSkillState,
	}

	err = s.repository.InsertPlayerState(state)
	if err != nil {
		return result, err
	}

	state.PlayerID = opponent.User.ID.String()
	err = s.repository.InsertPlayerState(state)
	if err != nil {
		return result, err
	}

	// TODO : remove duplicate param definition below and above
	err = s.repository.WithRedisTrx((*(client.Context)).Request().Context(), []string{redisPoolKey, poolCloneKey, moveCacheKey}, func(pipe redis.Pipeliner) error {
		err = s.repository.DeletePlayerFromPool(models.PlayerPoolParams{
			PoolParams: models.PoolParams{
				Type:        params.Type,
				TimeControl: params.TimeControl,
			},
			Player: opponent,
		}, pipe)
		if err != nil {
			return err
		}

		// delete redis hash key for player data clone
		err = s.repository.DeletePlayerOnPoolDataToRedis(models.PlayerPoolParams{
			User: models.User{
				ID: opponent.User.ID,
			},
		}, joinTime, pipe)
		if err != nil {
			return err
		}

		// insert game cache move ref to redis
		err = s.repository.InsertMoveCacheIdentifier(models.MoveCache{
			ID:                 moveCacheID,
			Turn:               true,
			LastMovement:       time.Now(),
			BlackTotalDuration: 0,
			WhiteTotalDuration: 0,
		}, pipe)
		if err != nil {
			return err
		}

		err = s.repository.InsertMatchSkillCount(models.InitMatchSkillStats{
			ID:         user.ID,
			GameSkills: skills,
		}, pipe)
		if err != nil {
			return err
		}

		err = s.repository.InsertMatchSkillCount(models.InitMatchSkillStats{
			ID:         opponent.User.ID,
			GameSkills: skills,
		}, pipe)
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

	room := s.wsConnections.CreateRoom(&models.GameRoom{
		Type: constants.WS_ROOM_TYPE_GAME,
	}, gameID.String())

	room.AddClient(client.User.ID.String())
	room.AddClient(opponent.User.ID.String())

	// insert to mysql
	// roomID, err := uuid.Parse(room.GetRoomID())
	// if err != nil {
	// 	return result, err
	// }

	err = s.repository.InsertGameData(models.GameActiveData{
		ID:                gameID,
		WhitePlayerID:     client.User.ID,
		BlackPlayerID:     opponent.User.ID,
		GameTypeVariantID: gameTypeVariantID,
		RoomID:            gameID,
		MovesCacheRef:     moveCacheID,
		StartTime:         time.Now().Format(time.RFC3339),
	})
	if err != nil {
		return result, err
	}

	s.wsConnections.EmitToRoom(models.WebSocketChannel{
		Source:       client.User.ID.String(),
		TargetClient: opponent.User.ID.String(),
		Event:        constants.WS_EVENT_EMIT_START_GAME,
		Data: models.GameData{
			ID: gameID.String(),
		},
		TargetRoom: room.GetRoomID(),
	})

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

func (s *webSocketService) CleanMatchupState(c echo.Context, user models.User) error {
	/*
		CASE 1 : player initiator is waiting for opponent then leaves page / refresh
				(no enemy has accepted the matchup)
			- delete player data from redis pool
			- delete player data copy (hash key-val)
	*/

	/*
		CASE 2 : found matchup (start game already)
			Note :
			- handle refresh (if possible, disable refresh)
			- handle move to another page (if move to another page by manually url rewrite, then redirect back to the game url)
			- disable going back to previous page
			- game invalidation / end only occurs when one of the player either leave the game or win the game

			Step to recover disconnection while game (run when init connection)
			-> get game data (exists on the first step)
			-> keep the room data (do not invalidate), but delete the old connection (i guess this is implemented already)
			-> reinitialize connection, insert player to room again
	*/

	if user.ID == uuid.Nil {
		return nil
	}

	_, err := s.repository.GetPlayerCurrentGameState(user.ID.String())
	if err != nil && err != errs.ERR_ACTIVE_GAME_NOT_FOUND {
		helpers.LogErrorCallStack(c, err)
	}

	if err == errs.ERR_ACTIVE_GAME_NOT_FOUND {
		// find pool data clone first

		playerPoolData, err := s.repository.GetPlayerPoolData(models.PlayerPoolParams{
			User: models.User{
				ID: user.ID,
			},
		})

		if err != nil {
			return err
		}

		poolPlayerCloneKey := helpers.GetPlayerPoolCloneKey(models.PlayerPoolParams{
			User: models.User{
				ID: user.ID,
			},
		})
		poolKey := helpers.GetPoolKey(models.PoolParams{
			Type:        playerPoolData["game-type"],
			TimeControl: playerPoolData["game-time-control"],
		})

		err = s.repository.WithRedisTrx(c.Request().Context(), []string{poolKey, poolPlayerCloneKey}, func(pipe redis.Pipeliner) error {

			// delete redis hash key for player data clone
			err = s.repository.DeletePlayerOnPoolDataToRedis(models.PlayerPoolParams{
				User: models.User{
					ID: user.ID,
				},
			}, time.Time{}, pipe)
			if err != nil {
				return err
			}

			err = s.repository.DeleteMatchSkillCount(models.InitMatchSkillStats{
				ID: user.ID,
			}, pipe)
			if err != nil {
				return err
			}

			joinTime, _ := time.Parse(time.RFC3339, playerPoolData["game-join-time"])

			err = s.repository.DeletePlayerFromPool(models.PlayerPoolParams{
				PoolParams: models.PoolParams{
					Type:        playerPoolData["game-type"],
					TimeControl: playerPoolData["game-time-control"],
				},
				Player: models.PlayerPool{
					User: models.User{
						ID:        user.ID,
						EloPoints: user.EloPoints,
					},
					JoinTime: joinTime,
				},
			}, pipe)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			helpers.LogErrorCallStack(c, err)
			return err
		}
		return nil
	}

	return nil
}

func (s *webSocketService) HandleConnectMatchSocketConnection(client models.WebSocketClientData, params models.HandleConnectMatchSocketConnectionParams) (models.HandleConnectMatchSocketConnectionResponse, error) {

	/*
		Step to recover disconnection while game (run when init connection)
		-> get game data (exists on the first step)
		-> keep the room data (do not invalidate), but delete the old connection (i guess this is implemented already)
		-> reinitialize connection, insert player to room again


		// CASE 1 : one player leaves, the other stays
		- check if the room still exists
		- if room still exists it means that the other stays,
		- rejoin immediately by adding the connection to the room
		- if room does not exist anymore, go to CASE 2


		// CASE 2 : both player leaves
		this cause the room to be deleted
		- recreate the room with the room id from the currentGameData
		- *the other clients will join this room
	*/

	var result models.HandleConnectMatchSocketConnectionResponse
	user := client.User

	currentGameData, err := s.repository.GetPlayerCurrentGameState(user.ID.String())
	if err != nil && err != errs.ERR_ACTIVE_GAME_NOT_FOUND {
		return result, err
	}

	gameRoom := s.wsConnections.GetRoomByID(currentGameData.ID.String())
	if gameRoom == nil {
		// CASE 2
		gameRoom = s.wsConnections.CreateRoom(&models.GameRoom{
			Type: constants.WS_ROOM_TYPE_GAME,
		}, currentGameData.ID.String())

		gameRoom.AddClient(user.ID.String())

		return result, err
	}

	// CASE 1
	gameRoom.AddClient(user.ID.String())

	return result, nil
}

func (s *webSocketService) HandleGamePublishAction(client models.WebSocketClientData, params models.HandleGamePublishActionParams) (models.HandleGamePublishActionResponse, error) {

	user := client.User

	game, err := s.repository.GetPlayerCurrentGameState(user.ID.String())
	if err != nil && err != errs.ERR_ACTIVE_GAME_NOT_FOUND {
		helpers.LogErrorCallStack(*client.Context, err)
		return models.HandleGamePublishActionResponse{}, err
	}

	// TODO : validate new state, send state on end game from client

	// if client color identifier is black then make turn true for white
	// true -> white
	// false -> black

	var opponentID uuid.UUID
	if game.WhitePlayerID == user.ID {
		opponentID = game.BlackPlayerID
	} else {
		opponentID = game.WhitePlayerID
	}

	playerState, err := s.repository.GetPlayerState(models.PlayerState{
		PlayerID: user.ID.String(),
		GameID:   game.ID.String(),
	})
	if err != nil {
		return models.HandleGamePublishActionResponse{}, err
	}

	gameMove, err := s.repository.GetMoveCacheIdentifier(models.MoveCache{
		ID: game.MovesCacheRef,
	})
	if err != nil {
		return models.HandleGamePublishActionResponse{}, err
	}

	oldState := gameMove["move"]
	oldStateToArr := helpers.ConvertNotationToArray(oldState)
	if user.ID == game.BlackPlayerID {
		oldStateToArr = helpers.TransformBoard(oldStateToArr)
	}
	oldState = helpers.ConvertArrayToNotation(oldStateToArr)

	newState := params.State
	newStateToArr := helpers.ConvertNotationToArray(newState)
	if user.ID == game.BlackPlayerID {
		newStateToArr = helpers.TransformBoard(newStateToArr)
	}
	newState = helpers.ConvertArrayToNotation(newStateToArr)

	if os.Getenv("VALIDATE_MOVE") == "ON" {
		validator, err := s.rpcClient.MatchServiceRpc.ValidateMove((*(client.Context)).Request().Context(), &pb.ValidateMoveReq{
			OldState: oldState,
			NewState: newState,
		})
		if err != nil {
			return models.HandleGamePublishActionResponse{}, err
		}

		if !validator.Valid {
			return models.HandleGamePublishActionResponse{}, errs.ERR_INVALID_MOVE
		}
	}

	var newTurn bool
	if gameMove["turn"] == "1" {
		newTurn = false
	} else {
		newTurn = true
	}

	// decrement all skill state duration by one
	if len(playerState.BuffState) > 0 {
		for name, skill := range playerState.BuffState {
			if skill.DurationLeft > 0 {
				skill.DurationLeft -= 1
			}
			for i, effects := range skill.List {
				if effects.DurationLeft > 0 {
					skill.List[i].DurationLeft -= 1
				}
			}
			playerState.BuffState[name] = skill
		}
	}

	if len(playerState.DebuffState) > 0 {
		for name, skill := range playerState.DebuffState {
			if skill.DurationLeft > 0 {
				skill.DurationLeft -= 1
			}
			for i, effects := range skill.List {
				if effects.DurationLeft > 0 {
					skill.List[i].DurationLeft -= 1
				}
			}
			playerState.DebuffState[name] = skill
		}
	}

	if len(playerState.BuffState) > 0 || len(playerState.DebuffState) > 0 {
		err = s.repository.UpdatePlayerState(models.PlayerState{
			PlayerID:    user.ID.String(),
			GameID:      game.ID.String(),
			BuffState:   playerState.BuffState,
			DebuffState: playerState.DebuffState,
		})
		if err != nil {
			return models.HandleGamePublishActionResponse{}, err
		}
	}

	// if turn is 1, then last movement timestamp belongs to black
	// if turn is 0, then last movement timestamp belongs to white

	previousMoveTimestamp, err := time.Parse(time.RFC3339, gameMove["last_movement"])
	if err != nil {
		return models.HandleGamePublishActionResponse{}, err
	}

	currentTime := time.Now()
	countDuration := currentTime.Sub(previousMoveTimestamp)

	var whiteSpentDuration int64
	var blackSpentDuration int64

	if newTurn { // black is moving
		intCurrentCumulativeDuration, err := strconv.ParseInt(gameMove["black_total_duration"], 10, 64)
		if err != nil {
			return models.HandleGamePublishActionResponse{}, err
		}
		intEnemyCumulativeDuration, err := strconv.ParseInt(gameMove["white_total_duration"], 10, 64)
		if err != nil {
			return models.HandleGamePublishActionResponse{}, err
		}
		cumulativeDuration := intCurrentCumulativeDuration + int64(countDuration.Seconds()) - game.Increment
		if cumulativeDuration > game.Duration {
			return models.HandleGamePublishActionResponse{}, errs.ERR_GAME_TIMEOUT
		}

		whiteSpentDuration = intEnemyCumulativeDuration
		blackSpentDuration = cumulativeDuration

		err = s.repository.InsertMoveCacheIdentifier(models.MoveCache{
			ID:                 game.MovesCacheRef,
			Move:               params.State,
			Turn:               newTurn,
			LastMovement:       currentTime,
			BlackTotalDuration: cumulativeDuration,
			WhiteTotalDuration: intEnemyCumulativeDuration,
		}, nil)
	} else { // white is moving
		intCurrentCumulativeDuration, err := strconv.ParseInt(gameMove["white_total_duration"], 10, 64)
		if err != nil {
			return models.HandleGamePublishActionResponse{}, err
		}
		intEnemyCumulativeDuration, err := strconv.ParseInt(gameMove["black_total_duration"], 10, 64)
		if err != nil {
			return models.HandleGamePublishActionResponse{}, err
		}
		cumulativeDuration := intCurrentCumulativeDuration + int64(countDuration.Seconds()) - game.Increment
		if cumulativeDuration > game.Duration {
			return models.HandleGamePublishActionResponse{}, errs.ERR_GAME_TIMEOUT
		}

		blackSpentDuration = intEnemyCumulativeDuration
		whiteSpentDuration = cumulativeDuration

		err = s.repository.InsertMoveCacheIdentifier(models.MoveCache{
			ID:                 game.MovesCacheRef,
			Move:               params.State,
			Turn:               newTurn,
			LastMovement:       currentTime,
			WhiteTotalDuration: cumulativeDuration,
			BlackTotalDuration: intEnemyCumulativeDuration,
		}, nil)
	}
	if err != nil {
		return models.HandleGamePublishActionResponse{}, err
	}

	gameRoom := s.wsConnections.GetRoomByID(game.ID.String())
	if gameRoom == nil {
		gameRoom = s.wsConnections.CreateRoom(&models.GameRoom{
			Type: constants.WS_ROOM_TYPE_GAME,
		}, game.ID.String())
	}

	gameRoom.AddClient(user.ID.String())

	s.wsConnections.EmitToRoom(models.WebSocketChannel{
		Source:       user.ID.String(),
		TargetClient: opponentID.String(),
		Event:        constants.WS_EVENT_EMIT_UPDATE_GAME_STATE,
		Data: models.HandleGamePublishActionResponse{
			State:              params.State,
			Turn:               user.ID == game.BlackPlayerID,
			WhiteSpentDuration: whiteSpentDuration,
			BlackSpentDuration: blackSpentDuration,
		},
		TargetRoom: gameRoom.GetRoomID(),
	})
	fmt.Println(whiteSpentDuration, blackSpentDuration)
	return models.HandleGamePublishActionResponse{
		State:              params.State,
		Turn:               newTurn,
		WhiteSpentDuration: whiteSpentDuration,
		BlackSpentDuration: blackSpentDuration,
	}, nil
}

func (s *webSocketService) ApplySkillEffects(gameID uuid.UUID, skillId uuid.UUID, playerID uuid.UUID, opponentID uuid.UUID, position models.Position) error {

	// skill info mysql
	skillDataList, err := s.repository.GetGameSkills(models.GameSkill{
		ID: skillId,
	})
	if err != nil {
		return err
	}

	var skillInfo models.GameSkill
	if len(skillDataList) <= 0 {
		return errs.ERR_GAME_SKILL_DATA_NOT_FOUND
	}
	skillInfo = skillDataList[0]

	if skillInfo.Permanent {
		return nil
	}

	var executionTarget uuid.UUID
	if skillInfo.Type == constants.SKILL_TYPE_BUFF {
		executionTarget = playerID
	} else {
		executionTarget = opponentID
	}

	playerState, err := s.repository.GetPlayerState(models.PlayerState{
		PlayerID: executionTarget.String(),
		GameID:   gameID.String(),
	})
	if err != nil {
		return err
	}

	var updateState map[string]models.SkillState
	if skillInfo.Type == constants.SKILL_TYPE_BUFF {
		updateState = playerState.BuffState
	} else {
		updateState = playerState.DebuffState
	}

	if updateState == nil {
		updateState = make(map[string]models.SkillState)
	}

	// logic
	state := updateState[skillInfo.Name]
	if !skillInfo.AutoTrigger {
		state.List = append(state.List, models.SkillStatus{
			Position: models.SkillPosition{
				Row: position.Row,
				Col: position.Col,
			},
			DurationLeft: skillInfo.Duration,
		})
		updateState[skillInfo.Name] = state
	} else {
		state.DurationLeft = skillInfo.Duration
		updateState[skillInfo.Name] = state
	}

	params := models.PlayerState{
		PlayerID: executionTarget.String(),
		GameID:   gameID.String(),
	}

	if skillInfo.Type == constants.SKILL_TYPE_BUFF {
		params.BuffState = updateState
	} else {
		params.DebuffState = updateState
	}

	err = s.repository.UpdatePlayerState(params)
	if err != nil {
		return err
	}

	return nil
}
