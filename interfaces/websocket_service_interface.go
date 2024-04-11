package interfaces

import "ingenhouzs.com/chesshouzs/go-game/models"

type WebsocketService interface {
	HandleMatchmaking(conn models.WebSocketClientConnection) (models.WebSocketResponse, error)
}
