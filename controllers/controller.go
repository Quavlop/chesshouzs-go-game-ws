package controllers

import (
	"net/http"
	"net/http/pprof"

	"github.com/labstack/echo"
	"ingenhouzs.com/chesshouzs/go-game/interfaces"
)

type Controller struct {
	HttpService      interfaces.HttpService
	WebSocketService interfaces.WebsocketService
	Repository       interfaces.Repository
}

func NewController(e *echo.Echo, httpService interfaces.HttpService, webSocketService interfaces.WebsocketService, repository interfaces.Repository) *Controller {
	controller := &Controller{
		HttpService:      httpService,
		WebSocketService: webSocketService,
		Repository:       repository,
	}

	HttpRoutes(e, controller)
	PprofRoutes(e, controller)
	return controller

}

func HttpRoutes(e *echo.Echo, controller *Controller) {
	// route := e.Group("/rest")

	// middlewares
	// route.Use(middlewares.Auth(controller.Repository))

	// route.GET("", func(c echo.Context) error {
	// 	return helpers.HttpResponse(c, http.StatusOK, "HTTP ROUTES BABY", models.Response{Status: 2})
	// })
}

func PprofRoutes(e *echo.Echo, controller *Controller) {
	route := e.Group("/debug/pprof")

	route.GET("", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	route.GET("/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))
	route.GET("/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
}
