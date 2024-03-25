package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"ingenhouzs.com/chesshouzs/go-game/interfaces"
	"ingenhouzs.com/chesshouzs/go-game/models"
	"ingenhouzs.com/chesshouzs/go-game/repositories"
)

type Controller struct {
	Service interfaces.Service
}

func NewController(e *echo.Echo, service interfaces.Service) {
	controller := &Controller{
		Service: service,
	}

	GameRoutes(e, controller)

}

func GameRoutes(e *echo.Echo, controller *Controller) {
	route := e.Group("/")

	route.GET("", func(c echo.Context) error {
		_, err := repositories.ConnectPostgreSQL(models.SqlConnection{})
		if err != nil {
			log.Errorf("wkwkw")
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, "sr")
	})
	route.POST("", func(c echo.Context) error { return c.JSON(http.StatusInternalServerError, "sssss") })
	// route.POST("", func(c echo.Context) error { return errors.New("ERROR TEST") })

}
