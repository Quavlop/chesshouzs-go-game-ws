package services

import (
	"os"
	"strconv"

	"ingenhouzs.com/chesshouzs/go-game/helpers"
	"ingenhouzs.com/chesshouzs/go-game/helpers/errs"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

// Websocket services
func (s *webSocketService) HandleMatchmaking(client models.WebSocketClientData, params models.HandleMatchmakingParams) (models.HandleMatchmakingResponse, error) {
	var result models.HandleMatchmakingResponse

	// TODO : Get player Data, JWT Middleware

	// get player data first
	// eligibleOpponents, err := s.FilterEligibleOpponent(client, models.FilterEligibleOpponentParams{
	// 	Filter: models.PoolParams{
	// 		Type:        params.Type,
	// 		TimeControl: params.TimeControl,
	// 	},
	// 	Client: models.PlayerPool{
	// 		EloPoints: 32,
	// 	},
	// })
	// if err != nil {
	// 	helpers.LogErrorCallStack(*client.Context, err)
	// 	return result, nil
	// }

	// if no one available then insert into pool
	// if found match then remove the enemy from pool and insert both into game data

	return result, nil
}

func (s *webSocketService) FilterEligibleOpponent(client models.WebSocketClientData, params models.FilterEligibleOpponentParams) (models.FilterEligibleOpponentResponse, error) {
	var result models.FilterEligibleOpponentResponse

	playerPool, err := s.repository.GetUnderMatchmakingPlayers(params.Filter)
	if err != nil {
		helpers.LogErrorCallStack(*client.Context, err)
		return result, err
	}

	playerPool, err = s.FilterOutOpponents(client, playerPool)
	if err != nil {
		helpers.LogErrorCallStack(*client.Context, err)
		return result, err
	}

	playerPool, err = s.SortPlayerPool(client, playerPool)
	if err != nil {
		helpers.LogErrorCallStack(*client.Context, err)
		return result, err
	}

	if len(playerPool) <= 0 {
		err = errs.ERR_NO_AVAILABLE_PLAYERS
		helpers.LogErrorCallStack(*client.Context, err)
		return result, err
	}

	return result, nil
}

func (s *webSocketService) SortPlayerPool(client models.WebSocketClientData, pool []models.PlayerPool) ([]models.PlayerPool, error) {
	for i, elemI := range pool {
		for j, elemJ := range pool {
			if s.PlayerSortFilter(elemI, elemJ) {
				temp := pool[i]
				pool[i] = pool[j]
				pool[j] = temp
			}
		}
	}

	return pool, nil
}

func (s *webSocketService) FilterOutOpponents(client models.WebSocketClientData, pool []models.PlayerPool) ([]models.PlayerPool, error) {
	var result []models.PlayerPool
	// get player data first
	for _, player := range pool {
		if s.IsMatchmakingEligible(player, player) {
			result = append(result, player)
		}
	}
	return result, nil
}

func (s *webSocketService) IsMatchmakingEligible(player models.PlayerPool, opponent models.PlayerPool) bool {
	threshold, err := strconv.Atoi(os.Getenv("ELO_GAP_THRESHOLD"))
	if err != nil {
		return player.EloPoints-opponent.EloPoints <= 150
	}
	return player.EloPoints-opponent.EloPoints <= int32(threshold)
}

func (s *webSocketService) PlayerSortFilter(playerOne models.PlayerPool, playerTwo models.PlayerPool) bool {
	if playerOne.EloPoints != playerTwo.EloPoints {
		return playerOne.EloPoints > playerTwo.EloPoints
	}
	return playerOne.JoinTime.After(playerTwo.JoinTime)
}
