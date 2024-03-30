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

	GameRoutes(e, controller)

}

func GameRoutes(e *echo.Echo, controller *Controller) {
	route := e.Group("/")

	route.GET("", func(c echo.Context) error {
		return helpers.HttpResponse(c, http.StatusOK, "WKWK", models.Response{Status: 2})
	})
}
