package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"ingenhouzs.com/chesshouzs/go-game/config/websocket"
	"ingenhouzs.com/chesshouzs/go-game/controllers"
	"ingenhouzs.com/chesshouzs/go-game/middlewares"
	"ingenhouzs.com/chesshouzs/go-game/models"
	"ingenhouzs.com/chesshouzs/go-game/repositories"
	"ingenhouzs.com/chesshouzs/go-game/services"
)

func main() {
	e := echo.New()

	if err := godotenv.Load(); err != nil {
		e.Logger.Fatal(err.Error())
	}

	// middlewares
	e.Use(middlewares.SetRequestID)
	e.Use(middleware.Recover())
	e.Use(middleware.BodyDump(middlewares.Logger))
	e.Use(middlewares.PanicLogger)

	postgresConnection := models.SqlConnection{
		Driver:   os.Getenv("POSTGRES_DB_DRIVER"),
		Host:     os.Getenv("POSTGRES_DB_HOST"),
		Port:     os.Getenv("POSTGRES_DB_PORT"),
		User:     os.Getenv("POSTGRES_DB_USER"),
		Password: os.Getenv("POSTGRES_DB_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB_DATABASE"),
	}

	redisConnection := models.RedisConnection{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	psql, err := repositories.ConnectPostgreSQL(postgresConnection)
	if err != nil {
		e.Logger.Fatal("Failed to connect PostgreSQL : " + err.Error())
	}

	redis, err := repositories.ConnectRedis(redisConnection)
	if err != nil {
		e.Logger.Fatal("Failed to connect Redis : " + err.Error())
	}

	// individual connection per client
	// key : user's session token
	// value : client connection metadata
	wsConnections := &websocket.Connections{} // TODO -> get user list from db

	// rooms connection
	// key : room_id
	// value : room_object consisting room_id and client list (map)
	wsGameRooms := make(map[string]*models.GameRoom) // TODO -> get room list from db
	wsConnections.Init()

	repository := repositories.NewRepository(psql, redis)

	// auth middleware
	e.Use(middlewares.Auth(repository))

	httpService := services.NewHttpService(repository, services.BaseService{})
	websocketService := services.NewWebSocketService(repository, wsConnections, services.BaseService{})

	baseService := services.NewBaseService(websocketService, httpService)
	baseService.WebSocketService = websocketService
	baseService.HttpService = httpService

	controller := controllers.NewController(e, httpService, websocketService)
	websocket.NewWebSocketHandler(e, controller, wsConnections, wsGameRooms)

	e.Logger.Fatal(e.Start(":" + os.Getenv("SERVICE_PORT")))
}
