package services

import (
	"math"
	"os"
	"strconv"

	"ingenhouzs.com/chesshouzs/go-game/helpers"
	"ingenhouzs.com/chesshouzs/go-game/helpers/errs"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

// Websocket services
func (s *webSocketService) HandleMatchmaking(client models.WebSocketClientData, params models.HandleMatchmakingParams) (models.HandleMatchmakingResponse, error) {
	var result models.HandleMatchmakingResponse

	user := client.User

	eligibleOpponents, err := s.baseService.WebSocketService.FilterEligibleOpponent(client, models.FilterEligibleOpponentParams{
		Filter: models.PoolParams{
			Type:        params.Type,
			TimeControl: params.TimeControl,
		},
		Client: models.PlayerPool{
			EloPoints: user.EloPoints,
		},
	})
	if err != nil {
		helpers.LogErrorCallStack(*client.Context, err)
		return result, err
	}

	// if no one available then insert into pool
	// if found match then remove the enemy from pool and insert both into game data

	opponent := eligibleOpponents.Player
	result.Opponent = opponent

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

	playerPool, err = s.baseService.WebSocketService.FilterOutOpponents(client, playerPool)
	if err != nil {
		helpers.LogErrorCallStack(*client.Context, err)
		return result, err
	}

	playerPool, err = s.baseService.WebSocketService.SortPlayerPool(client, playerPool)
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
			if s.baseService.WebSocketService.PlayerSortFilter(pool[minIdx], pool[idx]) {
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
		EloPoints: client.Data.(models.Player).EloPoints,
	}

	for _, opponent := range pool {
		if s.baseService.WebSocketService.IsMatchmakingEligible(player, opponent) {
			result = append(result, opponent)
		}
	}
	return result, nil
}

func (s *webSocketService) IsMatchmakingEligible(player models.PlayerPool, opponent models.PlayerPool) bool {
	threshold, err := strconv.Atoi(os.Getenv("ELO_GAP_THRESHOLD"))
	if err != nil {
		return math.Abs(float64(player.EloPoints-opponent.EloPoints)) <= 150
	}
	return math.Abs(float64(player.EloPoints-opponent.EloPoints)) <= float64(threshold)
}

func (s *webSocketService) PlayerSortFilter(playerOne models.PlayerPool, playerTwo models.PlayerPool) bool {
	if playerOne.EloPoints != playerTwo.EloPoints {
		return playerOne.EloPoints > playerTwo.EloPoints
	}
	return playerOne.JoinTime.After(playerTwo.JoinTime)
}
