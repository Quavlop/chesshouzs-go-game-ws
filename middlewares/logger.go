package middlewares

import (
	"bytes"
	"io"
	"os"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func makeLogEntry(c echo.Context) *log.Entry {
	if c == nil {
		return log.WithFields(log.Fields{
			"at": time.Now().Format("2006-01-02 15:04:05"),
		})
	}

	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Errorf("Error reading request body: %s", err)
		return nil
	}
	defer c.Request().Body.Close()

	c.Request().Body = io.NopCloser(bytes.NewReader(bodyBytes))

	bodyString := string(bodyBytes)

	return log.WithFields(log.Fields{
		"at":     time.Now().Format(os.Getenv("TIME_FORMAT")),
		"method": c.Request().Method,
		"uri":    c.Request().URL.String(),
		"ip":     c.Request().RemoteAddr,
		"body":   bodyString,
	})
}

func MiddlewareLogging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		makeLogEntry(c).Info("REQUEST")
		return next(c)
	}
}
