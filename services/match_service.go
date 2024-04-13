package services

import (
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

	return models.HandleMatchmakingResponse{
		ID: "1",
	}, nil
}
