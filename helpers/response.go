package helpers

import (
	"github.com/labstack/echo"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

func HttpResponse(e echo.Context, status int, message string, data interface{}) error {
	if status >= 500 {
		return e.JSON(status, models.Response{
			Status:    status,
			Error:     message,
			ErrorData: data,
		})
	}

	return e.JSON(status, models.Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}
