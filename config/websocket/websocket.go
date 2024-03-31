package websocket

import (
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
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

type wsHandler struct {
	conn    *ws.Conn
	service interfaces.WebsocketService
}

type WebSocketClientConnection struct {
	Token      string
	Connection *ws.Conn
}

type WebSocketClientMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func initConnection(c echo.Context, connectionList []*WebSocketClientConnection) (*ws.Conn, error) {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return conn, err
	}

	token := c.QueryParams().Get("sid")
	connectionList = append(connectionList, &WebSocketClientConnection{
		Token:      token,
		Connection: conn,
	})
	helpers.WriteOutLog("[WEBSOCKET] CONNECTION ESTABLISHED : \"" + c.Request().RemoteAddr + " | " + c.Request().Host + " | " + time.Now().Format(os.Getenv("TIME_FORMAT")) + "\"")
	return conn, nil
}

func NewWebSocketHandler(e *echo.Echo, service interfaces.WebsocketService, connectionList []*WebSocketClientConnection) {
	e.GET("/ws", func(c echo.Context) error {
		conn, err := initConnection(c, connectionList)
		if err != nil {
			return err
		}
		defer conn.Close()
		return handleIO(service, conn, connectionList)
	})
}

func handleIO(service interfaces.WebsocketService, conn *ws.Conn, connectionList []*WebSocketClientConnection) error {

	for {
		err := conn.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			log.Println(err)
			continue
		}

		var message WebSocketClientMessage
		err = conn.ReadJSON(&message)
		if err != nil {
			terminateConnection(conn, connectionList)
			break
		}

		if message.Event == "" {
			continue
		}
		handleEvents(service, conn, connectionList, message.Event)
	}
	return nil
}

func terminateConnection(conn *ws.Conn, connectionList []*WebSocketClientConnection) []*WebSocketClientConnection {
	for idx, connection := range connectionList {
		if connection.Connection == conn {
			connectionList := append(connectionList[:idx], connectionList[idx+1:]...)
			return connectionList
		}
	}
	return connectionList
}

func handleEvents(service interfaces.WebsocketService, conn *ws.Conn, connectionList []*WebSocketClientConnection, event string) {
	var eventHandler map[string]func(channel models.WebSocketChannel) (models.WebSocketResponse, error) = map[string]func(channel models.WebSocketChannel) (models.WebSocketResponse, error){
		constants.WS_EVENT_INIT_MATCHMAKING: service.HandleMatchmaking,
	}

	eventHandler[event](models.WebSocketChannel{})
}
