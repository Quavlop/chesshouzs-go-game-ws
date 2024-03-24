package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

	e.Use(middleware.BodyDump(middlewares.Logger))

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
		// e.Logger.Fatal("Failed to connect PostgreSQL : " + err.Error())
	}

	redis, err := repositories.ConnectRedis(redisConnection)
	if err != nil {
		// e.Logger.Fatal("Failed to connect Redis : " + err.Error())
	}

	repository := repositories.NewRepository(psql, redis)
	service := services.NewService(repository)
	controllers.NewController(e, service)

	e.Logger.Fatal(e.Start(":" + os.Getenv("SERVICE_PORT")))
}
