package helpers

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

func WriteOutLog(content string) error {
	err := WriteToFile(os.Getenv("LOG_OUT_PATH"), content+"\n")
	if err != nil {
		return err
	}
	return nil
}

func WriteErrLog(content string) error {
	err := WriteToFile(os.Getenv("LOG_ERR_PATH"), content+"\n")
	if err != nil {
		return err
	}
	return nil
}

func ParseHeadersToString(header http.Header) string {
	var result string

	for headerName, headerValueList := range header {
		result += headerName + ":"
		for _, headerValue := range headerValueList {
			result += headerValue
			if len(headerValueList) > 1 {
				result += ","
			}
		}
		result += "|"
	}

	return result
}

func GenerateRequestID(c echo.Context) string {
	if c.Request().Header.Get("X-Request-ID") != "" {
		return c.Request().Header.Get("X-Request-ID")
	}
	return uuid.NewString()
}

func MapStatusResponseToLogLevel(status int) string {
	if status >= 100 && status < 500 {
		return "INFO"
	}
	return "ERROR"
}

func CaptureStackTrace() string {
	stackBuf := make([]byte, 8192)
	stackSize := runtime.Stack(stackBuf, false)
	return string(stackBuf[:stackSize])
}

func LogErrorCallStack(c echo.Context, err error) {
	var errMessage string
	var data models.LogErrorCallStack
	if err != nil {
		errMessage = err.Error()
	} else {
		errMessage = CaptureStackTrace()
	}

	event, _ := c.Get("ws-event").(string)

	if event != "" {
		data = models.LogErrorCallStack{
			Level:     "ERROR",
			Type:      "WS_INTERNAL",
			Event:     event,
			RequestID: c.Get("request_id").(string),
			Time:      time.Now().Format(os.Getenv("TIME_FORMAT")),
			Message:   errMessage,
			URI:       c.Request().URL.String(),
		}
	} else {
		data = models.LogErrorCallStack{
			Level:     "ERROR",
			Type:      "HTTP_INTERNAL",
			RequestID: c.Get("request_id").(string),
			Time:      time.Now().Format(os.Getenv("TIME_FORMAT")),
			Message:   errMessage,
			URI:       c.Request().URL.String(),
		}
	}

	stringData, err := json.Marshal(data)
	if err != nil {
		message := "Failed to write response log : " + err.Error()
		WriteErrLog(message)
		log.Errorf(message)
		return
	}

	message := string(stringData)
	WriteErrLog(message)
	WriteOutLog(message)
	log.Errorf(message)
}
