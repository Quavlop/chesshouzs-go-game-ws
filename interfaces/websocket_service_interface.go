package interfaces

import (
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

type WebsocketService interface {
	MatchService
	PlayerService
}

type MatchService interface {
	HandleMatchmaking(client models.WebSocketClientData, params models.HandleMatchmakingParams) (models.HandleMatchmakingResponse, error)
	HandleConnectMatchSocketConnection(client models.WebSocketClientData, params models.HandleConnectMatchSocketConnectionParams) (models.HandleConnectMatchSocketConnectionResponse, error)
	HandleGamePublishAction(client models.WebSocketClientData, params models.HandleGamePublishActionParams) (models.HandleGamePublishActionResponse, error)
	FilterEligibleOpponent(client models.WebSocketClientData, params models.FilterEligibleOpponentParams) (models.FilterEligibleOpponentResponse, error)
	SortPlayerPool(client models.WebSocketClientData, pool []models.PlayerPool) ([]models.PlayerPool, error)
	FilterOutOpponents(client models.WebSocketClientData, pool []models.PlayerPool) ([]models.PlayerPool, error)
	IsMatchmakingEligible(player models.PlayerPool, opponent models.PlayerPool) bool
	PlayerSortFilter(playerOne models.PlayerPool, playerTwo models.PlayerPool) bool
	CleanMatchupState(c echo.Context, user models.User) error
	ApplySkillEffects(gameID uuid.UUID, skillId uuid.UUID, playerID uuid.UUID, opponentID uuid.UUID, position models.Position) error
}

type PlayerService interface {
}
