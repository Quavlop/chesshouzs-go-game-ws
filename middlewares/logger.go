package middlewares

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"ingenhouzs.com/chesshouzs/go-game/helpers"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

func SetRequestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := helpers.GenerateRequestID(c)
		c.Set("request_id", requestID)
		return next(c)
	}
}

func LogRequest(c echo.Context, requestBody []byte) models.RequestResponseBridge {
	location := helpers.GetLocalTimeZone()
	startTime := time.Now().In(location)
	requestID := c.Get("request_id").(string)
	data := models.RequestLogData{
		Level:     "INFO",
		Type:      "REQUEST",
		RequestID: requestID,
		Header:    helpers.ParseHeadersToString(c.Request().Header),
		Time:      time.Now().In(location).Format(os.Getenv("TIME_FORMAT")),
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
		log.Errorf(message)
		return models.RequestResponseBridge{RequestID: requestID, StartTime: startTime}
	}

	message := string(stringData)
	helpers.WriteOutLog(message)
	log.Info(message)

	return models.RequestResponseBridge{RequestID: requestID, StartTime: startTime}
}

func LogResponse(c echo.Context, requestMetadata models.RequestResponseBridge, responseBody []byte) {
	location := helpers.GetLocalTimeZone()
	logLevel := helpers.MapStatusResponseToLogLevel(c.Response().Status)
	status := c.Response().Status
	if len(responseBody) <= 0 {
		status = http.StatusInternalServerError
	}
	data := models.ResponseLogData{
		Level:        logLevel,
		Type:         "RESPONSE",
		RequestID:    requestMetadata.RequestID,
		Header:       helpers.ParseHeadersToString(c.Request().Header),
		Time:         time.Now().In(location).Format(os.Getenv("TIME_FORMAT")),
		URI:          c.Request().URL.String(),
		Status:       status,
		Response:     string(responseBody),
		LatencyHuman: time.Since(requestMetadata.StartTime).String(),
		BytesOut:     len(responseBody),
	}

	stringData, err := json.Marshal(data)
	if err != nil {
		message := "Failed to write response log : " + err.Error()
		helpers.WriteErrLog(message)
		log.Errorf(message)
		return
	}

	message := string(stringData)
	helpers.WriteOutLog(message)
	if logLevel == "ERROR" {
		helpers.WriteErrLog(message)
		log.Errorf(message)
	} else {
		log.Info(message)
	}
}

func PanicLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer helpers.PanicRecover(c)
		return next(c)
	}
}
func Logger(c echo.Context, requestBody []byte, responseBody []byte) {
	data := LogRequest(c, requestBody)
	if c.Response().Status >= 500 {
		helpers.LogErrorCallStack(c, nil)
	}
	LogResponse(c, data, responseBody)
}
