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

	HttpRoutes(e, controller)
	WebsocketRoutes(e, controller)
	return controller

}

func WebsocketRoutes(e *echo.Echo, controller *Controller) {
	route := e.Group("/ws")

	route.GET("", func(c echo.Context) error {
		return helpers.HttpResponse(c, http.StatusOK, "WKWKsssssssss", models.Response{Status: 2})
	})
}

func HttpRoutes(e *echo.Echo, controller *Controller) {
	route := e.Group("/rest")

	route.GET("", func(c echo.Context) error {
		return helpers.HttpResponse(c, http.StatusOK, "HTTP ROUTES BABY", models.Response{Status: 2})
	})
}
