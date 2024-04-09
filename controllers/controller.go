package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	"ingenhouzs.com/chesshouzs/go-game/helpers"
	"ingenhouzs.com/chesshouzs/go-game/interfaces"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

type Controller struct {
	Service interfaces.HttpService
}

func NewController(e *echo.Echo, service interfaces.HttpService) {
	controller := &Controller{
		Service: service,
	}

	WebsocketRoutes(e, controller)

}

func WebsocketRoutes(e *echo.Echo, controller *Controller) {
	route := e.Group("/")

	route.GET("", func(c echo.Context) error {
		return helpers.HttpResponse(c, http.StatusOK, "WKWKsssssssss", models.Response{Status: 2})
	})
}
