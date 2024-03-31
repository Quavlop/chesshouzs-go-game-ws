package interfaces

import "ingenhouzs.com/chesshouzs/go-game/models"

type WebsocketService interface {
	HandleMatchmaking(channel models.WebSocketChannel) (models.WebSocketResponse, error)
}
