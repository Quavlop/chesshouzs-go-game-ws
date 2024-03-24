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

func LogRequest(c echo.Context, requestBody []byte) string {
	data := models.RequestLogData{
		Level:     "INFO",
		Type:      "REQUEST",
		RequestID: uuid.NewString(),
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
		return data.RequestID
	}

	message := string(stringData)
	helpers.WriteOutLog(message)
	c.Logger().Debug(message)
	return data.RequestID
}

func LogResponse(c echo.Context, requestID string, responseBody []byte) {
	data := models.ResponseLogData{
		Level:     "INFO",
		Type:      "RESPONSE",
		RequestID: requestID,
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
	requestID := LogRequest(c, requestBody)
	LogResponse(c, requestID, responseBody)
}
