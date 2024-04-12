package helpers

import (
	"encoding/json"

	"ingenhouzs.com/chesshouzs/go-game/constants"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

func BindParams(data interface{}, target interface{}) error {
	// convert from map to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// convert from JSON to corresponding struct param
	err = json.Unmarshal(jsonData, target)
	if err != nil {
		return err
	}
	return nil
}

func ErrorWebSocketResponseWrap(event string, data interface{}) models.WebSocketResponse {
	return models.WebSocketResponse{
		Status: constants.WS_SERVER_RESPONSE_ERROR,
		Event:  event,
		Data:   data,
	}
}

func SuccessWebSocketResponseWrap(event string, data interface{}) models.WebSocketResponse {
	return models.WebSocketResponse{
		Status: constants.WS_SERVER_RESPONSE_SUCCESS,
		Event:  event,
		Data:   data,
	}
}
