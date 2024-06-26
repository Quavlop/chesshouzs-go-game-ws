package interfaces

import (
	"github.com/labstack/echo"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

type WebsocketService interface {
	MatchService
	PlayerService
}

type MatchService interface {
	HandleMatchmaking(client models.WebSocketClientData, params models.HandleMatchmakingParams) (models.HandleMatchmakingResponse, error)
	HandleRecoverMatchSocketConnection(client models.WebSocketClientData, params models.HandleRecoverMatchSocketConnectionParams) (models.HandleRecoverMatchSocketConnectionResponse, error)
	FilterEligibleOpponent(client models.WebSocketClientData, params models.FilterEligibleOpponentParams) (models.FilterEligibleOpponentResponse, error)
	SortPlayerPool(client models.WebSocketClientData, pool []models.PlayerPool) ([]models.PlayerPool, error)
	FilterOutOpponents(client models.WebSocketClientData, pool []models.PlayerPool) ([]models.PlayerPool, error)
	IsMatchmakingEligible(player models.PlayerPool, opponent models.PlayerPool) bool
	PlayerSortFilter(playerOne models.PlayerPool, playerTwo models.PlayerPool) bool
	CleanMatchupState(c echo.Context, user models.User) error
}

type PlayerService interface {
}
