package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"ingenhouzs.com/chesshouzs/go-game/config/websocket"
	"ingenhouzs.com/chesshouzs/go-game/constants"
	"ingenhouzs.com/chesshouzs/go-game/helpers"
	"ingenhouzs.com/chesshouzs/go-game/interfaces"
	"ingenhouzs.com/chesshouzs/go-game/models"
	"ingenhouzs.com/chesshouzs/go-game/services/rpc/pb"
)

type BaseService struct {
	WebSocketService interfaces.WebsocketService
	HttpService      interfaces.HttpService
}
type httpService struct {
	repository  interfaces.Repository
	rpcClient   *rpcClient
	BaseService *BaseService
}

type webSocketService struct {
	repository    interfaces.Repository
	rpcClient     *rpcClient
	wsConnections *websocket.Connections
	BaseService   *BaseService
}

type gameRoomService struct {
	room interfaces.WebSocketRoom
}

type rpcClient struct {
	MatchServiceRpc pb.MatchServiceClient
}

type KafkaConsumerConfig struct {
	Host            string
	GroupId         string
	Topics          []string
	AutoOffsetReset string
}

type kafkaConsumer struct {
	context       context.Context
	repository    interfaces.Repository
	wsConnections *websocket.Connections
	BaseService   *BaseService
}

func NewBaseService(webSocketService interfaces.WebsocketService, httpService interfaces.HttpService) *BaseService {
	return &BaseService{WebSocketService: webSocketService, HttpService: httpService}
}

func NewHttpService(repository interfaces.Repository, baseService *BaseService, rpcClient *rpcClient) interfaces.HttpService {
	return &httpService{repository, rpcClient, baseService}
}

func NewWebSocketService(repository interfaces.Repository, wsConnections *websocket.Connections, baseService *BaseService, rpcClient *rpcClient) interfaces.WebsocketService {
	return &webSocketService{repository, rpcClient, wsConnections, baseService}
}

func NewRpcClient(serverHost string) (rpcClient, error) {
	var client rpcClient

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serverHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return client, err
	}
	defer conn.Close()

	client = rpcClient{
		MatchServiceRpc: pb.NewMatchServiceClient(conn),
	}

	return client, nil
}

func NewKafkaImpl(repository interfaces.Repository, wsConnections *websocket.Connections, service *BaseService, c context.Context) kafkaConsumer {
	return kafkaConsumer{
		context:       c,
		repository:    repository,
		wsConnections: wsConnections,
		BaseService:   service,
	}
}

func NewKafkaConsumer(config KafkaConsumerConfig, kafkaImpl kafkaConsumer) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.Host,
		"group.id":          config.GroupId,
		"auto.offset.reset": config.AutoOffsetReset,
	})
	if err != nil {
		helpers.WriteOutLog("[KAFKA CONSUMER] Error connecting kafka : " + err.Error())
		panic(err)
	}
	defer c.Close()

	c.SubscribeTopics(config.Topics, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err != nil {
			continue
		}

		rawMessage := string(msg.Value)
		cleanedMessage, err := strconv.Unquote(rawMessage)
		if err != nil {
			helpers.WriteOutLog(fmt.Sprintf("[KAFKA CONSUMER] Failed to unescape message on topic %s : %s", *&msg.TopicPartition.Topic, err.Error()))
			continue
		}

		var consumerErr error
		switch *msg.TopicPartition.Topic {
		case constants.EXECUTE_SKILL_TOPIC:
			var message models.ExecuteSkillMessage
			err := json.Unmarshal([]byte(cleanedMessage), &message)
			if err != nil {
				helpers.WriteOutLog(fmt.Sprintf("[KAFKA CONSUMER] Failed to parse message on topic %s : %s", *&msg.TopicPartition.Topic, err.Error()))
				continue
			}
			consumerErr = kafkaImpl.ExecuteSkillConsumer(message)
			if consumerErr != nil {
				helpers.WriteOutLog(fmt.Sprintf("[KAFKA CONSUMER] Error when consuming message on topic %s : %s", *msg.TopicPartition.Topic, consumerErr.Error()))
			}
		case constants.END_GAME_TOPIC:
			var message models.EndGameMessage
			err := json.Unmarshal([]byte(cleanedMessage), &message)
			if err != nil {
				helpers.WriteOutLog(fmt.Sprintf("[KAFKA CONSUMER] Failed to parse message on topic %s : %s", *&msg.TopicPartition.Topic, err.Error()))
				continue
			}
			consumerErr = kafkaImpl.EndGameConsumer(message)
			if consumerErr != nil {
				helpers.WriteOutLog(fmt.Sprintf("[KAFKA CONSUMER] Error when consuming message on topic %s : %s", *msg.TopicPartition.Topic, consumerErr.Error()))
			}
		}
	}
}
