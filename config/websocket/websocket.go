package websocket

import (
	"errors"
	"log"
	"os"
	"time"

	ws "github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"ingenhouzs.com/chesshouzs/go-game/constants"
	"ingenhouzs.com/chesshouzs/go-game/helpers"
	"ingenhouzs.com/chesshouzs/go-game/interfaces"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocketClientMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func initConnection(c echo.Context, connectionList *Connections) (*ws.Conn, error) {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return conn, err
	}

	token := c.QueryParams().Get("sid")
	connectionList.addConnection(token, conn)

	helpers.WriteOutLog("[WEBSOCKET] CONNECTION ESTABLISHED : \"" + c.Request().RemoteAddr + " | " + c.Request().Host + " | " + time.Now().Format(os.Getenv("TIME_FORMAT")) + "\"")
	return conn, nil
}

func NewWebSocketHandler(e *echo.Echo, service interfaces.WebsocketService, connectionList *Connections, gameRoomList map[string]*models.GameRoom) {
	e.GET("/ws", func(c echo.Context) error {
		conn, err := initConnection(c, connectionList)
		if err != nil {
			return err
		}
		defer conn.Close()
		return handleIO(c, service, conn, connectionList)
	})
}

func handleIO(c echo.Context, service interfaces.WebsocketService, conn *ws.Conn, connectionList *Connections) error {

	for {

		var message WebSocketClientMessage
		err := conn.ReadJSON(&message)
		if err != nil {
			token := c.QueryParams().Get("sid")
			connectionList.deleteConnection(token, conn)
			break
		}

		if message.Event == "" {
			continue
		}
		response, err := handleEvents(service, conn, connectionList, message.Event)

		err = conn.WriteJSON(response)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}

func handleEvents(service interfaces.WebsocketService, conn *ws.Conn, connectionList *Connections, event string) (models.WebSocketResponse, error) {
	var eventHandler map[string]func(channel models.WebSocketChannel) (models.WebSocketResponse, error) = map[string]func(channel models.WebSocketChannel) (models.WebSocketResponse, error){
		constants.WS_EVENT_INIT_MATCHMAKING: service.HandleMatchmaking,
	}

	handler, eventExists := eventHandler[event]
	if !eventExists {
		errMessage := "[ERROR] : 404 - Unrecognizable Event"
		return models.WebSocketResponse{
			Status: constants.WS_SERVER_RESPONSE_ERROR,
			Data:   errMessage,
		}, errors.New(errMessage)
	}
	response, err := handler(models.WebSocketChannel{})
	if err != nil {
		return response, err
	}
	return response, nil
}
