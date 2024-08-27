package main

import (
	"context"
	"os"
	"strconv"
	"strings"

	_ "net/http/pprof"

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

// individual connection per client
// key : user's session token
// value : client connection metadata
// TODO -> get user list from db
var wsConnections *websocket.Connections = &websocket.Connections{}

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

	cassandraConnection := models.CassandraConnection{
		Host:     os.Getenv("CASSANDRA_HOST"),
		Keyspace: os.Getenv("CASSANDRA_KEYSPACE"),
	}

	psql, err := repositories.ConnectPostgreSQL(postgresConnection)
	if err != nil {
		e.Logger.Fatal("Failed to connect PostgreSQL : " + err.Error())
	}

	redis, err := repositories.ConnectRedis(redisConnection)
	if err != nil {
		e.Logger.Fatal("Failed to connect Redis : " + err.Error())
	}

	cassandra, err := repositories.ConnectCassandra(cassandraConnection)
	if err != nil {
		e.Logger.Fatal("Failed to connect Apache Cassandra : " + err.Error())
	}

	rpcClient, err := services.NewRpcClient(os.Getenv("RPC_SERVER"))
	if err != nil {
		e.Logger.Fatal("Failed to create rpc client : " + err.Error())
	}

	// rooms connection
	// key : room_id
	// value : room_object consisting room_id and client list (map)
	wsConnections.Init()

	repository := repositories.NewRepository(psql, redis, cassandra)

	httpService := services.NewHttpService(repository, &services.BaseService{}, &rpcClient)
	websocketService := services.NewWebSocketService(repository, wsConnections, &services.BaseService{}, &rpcClient)

	baseService := services.NewBaseService(websocketService, httpService)
	baseService.WebSocketService = websocketService
	baseService.HttpService = httpService

	httpService = services.NewHttpService(repository, baseService, &rpcClient)
	websocketService = services.NewWebSocketService(repository, wsConnections, baseService, &rpcClient)

	controller := controllers.NewController(e, httpService, websocketService, repository)
	websocket.NewWebSocketHandler(e, controller, wsConnections)

	kafkaImpl := services.NewKafkaImpl(repository, wsConnections, baseService, context.Background())

	consumerWorker, err := strconv.Atoi(os.Getenv("KAFKA_CONSUMER_WORKER"))
	if err != nil {
		consumerWorker = 3
	}

	for i := 0; i < consumerWorker; i++ {
		go services.NewKafkaConsumer(services.KafkaConsumerConfig{
			Host:            os.Getenv("KAFKA_HOST"),
			GroupId:         os.Getenv("KAFKA_CONSUMER_GROUP_ID"),
			Topics:          strings.Split(os.Getenv("KAFKA_TOPICS"), ","),
			AutoOffsetReset: os.Getenv("KAFKA_CONFIG_AUTO_OFFSET_RESET"),
		}, kafkaImpl)
	}

	e.Logger.Fatal(e.Start(":" + os.Getenv("SERVICE_PORT")))
}
