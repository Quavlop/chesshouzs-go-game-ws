package helpers

import (
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/labstack/echo"
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
