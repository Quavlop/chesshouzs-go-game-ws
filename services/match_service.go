package services

import (
	"ingenhouzs.com/chesshouzs/go-game/constants"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

// Websocket services
func (s *webSocketService) HandleMatchmaking(channel models.WebSocketChannel) (models.WebSocketResponse, error) {
	return models.WebSocketResponse{
		Status: constants.WS_SERVER_RESPONSE_SUCCESS,
		Source: "WKKWKWKKWKW",
	}, nil
}
