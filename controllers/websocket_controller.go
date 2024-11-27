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

func (c *Controller) HandleConnectMatchSocketConnection(client models.WebSocketClientData) (models.WebSocketResponse, error) {
	serviceParams := models.HandleConnectMatchSocketConnectionParams{}

	err := helpers.BindParams(client.Data, &serviceParams)
	if err != nil {
		return helpers.ErrorWebSocketResponseWrap(client.Event, "INVALID PAYLOAD"), err
	}

	data, err := c.WebSocketService.HandleConnectMatchSocketConnection(client, serviceParams)
	if err != nil {
		return helpers.ErrorWebSocketResponseWrap(client.Event, "Failed to recover match connection : "+err.Error()), err
	}

	return helpers.SuccessWebSocketResponseWrap(client.Event, data), nil
}

func (c *Controller) HandleGamePublishAction(client models.WebSocketClientData) (models.WebSocketResponse, error) {
	serviceParams := models.HandleGamePublishActionParams{}

	err := helpers.BindParams(client.Data, &serviceParams)
	if err != nil {
		return helpers.ErrorWebSocketResponseWrap(client.Event, "INVALID PAYLOAD"), err
	}

	data, err := c.WebSocketService.HandleGamePublishAction(client, serviceParams)
	if err != nil {
		return helpers.ErrorWebSocketResponseWrap(client.Event, "Failed to recover match connection : "+err.Error()), err
	}

	return helpers.SuccessWebSocketResponseWrap(client.Event, data), nil
}

func (c *Controller) GetGameTimeDuration(client models.WebSocketClientData) (models.WebSocketResponse, error) {
	serviceParams := models.GetGameTimeDurationParams{}

	err := helpers.BindParams(client.Data, &serviceParams)
	if err != nil {
		return helpers.ErrorWebSocketResponseWrap(client.Event, "INVALID PAYLOAD"), err
	}

	data, err := c.WebSocketService.GetGameTimeDuration(client, serviceParams)
	if err != nil {
		return helpers.ErrorWebSocketResponseWrap(client.Event, "Failed to recover match connection : "+err.Error()), err
	}

	return helpers.SuccessWebSocketResponseWrap(client.Event, data), nil
}
