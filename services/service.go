package services

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"ingenhouzs.com/chesshouzs/go-game/config/websocket"
	"ingenhouzs.com/chesshouzs/go-game/interfaces"
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
