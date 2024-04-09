package services

import (
	"ingenhouzs.com/chesshouzs/go-game/models"
)

// Websocket services
func (s *webSocketService) HandleMatchmaking(channel models.WebSocketChannel) (models.WebSocketResponse, error) {
	return models.WebSocketResponse{
		Source: "WKKWKWKKWKW",
	}, nil
}
