package middlewares

import (
	"encoding/json"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"ingenhouzs.com/chesshouzs/go-game/helpers"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

func LogRequest(c echo.Context, requestBody []byte) models.RequestResponseBridge {
	startTime := time.Now()
	requestID := uuid.NewString()
	data := models.RequestLogData{
		Level:     "INFO",
		Type:      "REQUEST",
		RequestID: requestID,
		Header:    helpers.ParseHeadersToString(c.Request().Header),
		Time:      time.Now().Format(os.Getenv("TIME_FORMAT")),
		Host:      c.Request().Host,
		Method:    c.Request().Method,
		URI:       c.Request().URL.String(),
		Body:      string(requestBody),
		RemoteIP:  c.Request().RemoteAddr,
		BytesIn:   len(requestBody),
	}

	stringData, err := json.Marshal(data)
	if err != nil {
		message := "Failed to write request log : " + err.Error()
		helpers.WriteErrLog(message)
		c.Logger().Debug(message)
		return models.RequestResponseBridge{RequestID: requestID, StartTime: startTime}
	}

	message := string(stringData)
	helpers.WriteOutLog(message)
	c.Logger().Debug(message)
	return models.RequestResponseBridge{RequestID: requestID, StartTime: startTime}
}

func LogResponse(c echo.Context, requestMetadata models.RequestResponseBridge, responseBody []byte) {
	data := models.ResponseLogData{
		Level:        "INFO",
		Type:         "RESPONSE",
		RequestID:    requestMetadata.RequestID,
		Header:       helpers.ParseHeadersToString(c.Request().Header),
		Time:         time.Now().Format(os.Getenv("TIME_FORMAT")),
		URI:          c.Request().URL.String(),
		Status:       c.Response().Status,
		Response:     string(responseBody),
		LatencyHuman: time.Since(requestMetadata.StartTime).String(),
		BytesOut:     len(responseBody),
	}

	stringData, err := json.Marshal(data)
	if err != nil {
		message := "Failed to write response log : " + err.Error()
		helpers.WriteErrLog(message)
		c.Logger().Debug(message)
	}

	message := string(stringData)
	helpers.WriteOutLog(message)
	c.Logger().Debug(message)
}

func Logger(c echo.Context, requestBody []byte, responseBody []byte) {
	data := LogRequest(c, requestBody)
	LogResponse(c, data, responseBody)
}
