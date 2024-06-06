package helpers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

func PanicRecover(c echo.Context) {
	if r := recover(); r != nil {
		stackTrace := CaptureStackTrace()

		data := models.ResponseLogData{
			Level:     "ERROR",
			Type:      "INTERNAL_PANIC",
			RequestID: c.Get("request_id").(string),
			Header:    ParseHeadersToString(c.Request().Header),
			Time:      time.Now().Format(os.Getenv("TIME_FORMAT")),
			URI:       c.Request().URL.String(),
			Status:    http.StatusInternalServerError,
			Message:   stackTrace,
		}

		stringData, err := json.Marshal(data)
		if err != nil {
			message := "Failed to write response log : " + err.Error()
			WriteErrLog(message)
			log.Errorf(message)
			return
		}

		message := string(stringData)
		WriteOutLog(message)
		WriteErrLog(message)
	}
}
