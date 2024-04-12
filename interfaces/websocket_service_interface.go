package interfaces

import "ingenhouzs.com/chesshouzs/go-game/models"

type WebsocketService interface {
	HandleMatchmaking(client models.WebSocketClientData, params models.HandleMatchmakingParams) (models.HandleMatchmakingResponse, error)
}
