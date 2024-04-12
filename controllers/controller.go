package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	"ingenhouzs.com/chesshouzs/go-game/helpers"
	"ingenhouzs.com/chesshouzs/go-game/interfaces"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

type Controller struct {
	HttpService      interfaces.HttpService
	WebSocketService interfaces.WebsocketService
}

func NewController(e *echo.Echo, httpService interfaces.HttpService, webSocketService interfaces.WebsocketService) *Controller {
	controller := &Controller{
		HttpService:      httpService,
		WebSocketService: webSocketService,
	}

	WebsocketRoutes(e, controller)
	return controller

}

func WebsocketRoutes(e *echo.Echo, controller *Controller) {
	route := e.Group("/")

	route.GET("", func(c echo.Context) error {
		return helpers.HttpResponse(c, http.StatusOK, "WKWKsssssssss", models.Response{Status: 2})
	})
}
