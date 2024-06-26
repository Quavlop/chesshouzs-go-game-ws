package controllers

import (
	"ingenhouzs.com/chesshouzs/go-game/helpers"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

func (c *Controller) HandleMatchmaking(client models.WebSocketClientData) (models.WebSocketResponse, error) {
	serviceParams := models.HandleMatchmakingParams{}

	err := helpers.BindParams(client.Data, &serviceParams)
	if err != nil {
		return helpers.ErrorWebSocketResponseWrap(client.Event, "INVALID PAYLOAD"), err
	}

	data, err := c.WebSocketService.HandleMatchmaking(client, serviceParams)
	if err != nil {
		return helpers.ErrorWebSocketResponseWrap(client.Event, "Failed to handle matchmaking : "+err.Error()), err
	}

	return helpers.SuccessWebSocketResponseWrap(client.Event, data), nil
}

func (c *Controller) HandleRecoverMatchSocketConnection(client models.WebSocketClientData) (models.WebSocketResponse, error) {
	serviceParams := models.HandleRecoverMatchSocketConnectionParams{}

	err := helpers.BindParams(client.Data, &serviceParams)
	if err != nil {
		return helpers.ErrorWebSocketResponseWrap(client.Event, "INVALID PAYLOAD"), err
	}

	data, err := c.WebSocketService.HandleRecoverMatchSocketConnection(client, serviceParams)
	if err != nil {
		return helpers.ErrorWebSocketResponseWrap(client.Event, "Failed to recover match connection : "+err.Error()), err
	}

	return helpers.SuccessWebSocketResponseWrap(client.Event, data), nil
}
