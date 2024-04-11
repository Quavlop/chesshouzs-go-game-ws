package services

import (
	"ingenhouzs.com/chesshouzs/go-game/constants"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

// Websocket services
func (s *webSocketService) HandleMatchmaking(conn models.WebSocketClientConnection) (models.WebSocketResponse, error) {
	// s.wsConnections.EmitGlobalBroadcast(models.WebSocketChannel{
	// 	Source: conn.Token,
	// 	Event:  "MATCH",
	// 	Data:   "BROADCASTEEEEDDDDDDDDD FROM " + conn.Token,
	// })
	return models.WebSocketResponse{
		Status: constants.WS_SERVER_RESPONSE_SUCCESS,
	}, nil
}
