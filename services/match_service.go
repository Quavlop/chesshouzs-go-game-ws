package services

import (
	"errors"

	"ingenhouzs.com/chesshouzs/go-game/models"
)

// Websocket services
func (s *webSocketService) HandleMatchmaking(client models.WebSocketClientData, params models.HandleMatchmakingParams) (models.HandleMatchmakingResponse, error) {

	// playerPool, err := s.repository.GetUnderMatchmakingPlayers(models.PoolParams{
	// 	Type:        params.Type,
	// 	TimeControl: params.TimeControl,
	// })
	// if err != nil {
	// 	return models.WebSocketResponse{
	// 		Status: constants.WS_SERVER_RESPONSE_ERROR,
	// 		Event:  "MM",
	// 		Data:   "WKWWKKW",
	// 	}, nil
	// }

	// LOGGER level callstack
	// MAKE ALL LOGGING AS MIDDLEWARE FORMAT
	return models.HandleMatchmakingResponse{
		ID: "1",
	}, errors.New("KKWK")
}
