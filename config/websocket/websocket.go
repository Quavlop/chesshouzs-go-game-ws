package websocket

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	ws "github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"ingenhouzs.com/chesshouzs/go-game/interfaces"
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

	return conn, nil
}

func NewWebSocketHandler(e *echo.Echo, service interfaces.WebsocketService, connectionList []*WebSocketClientConnection) {
	e.GET("/ws", func(c echo.Context) error {
		conn, err := initConnection(c, connectionList)
		if err != nil {
			return err
		}
		defer conn.Close()
		return handleIO(conn)
	})
}

func handleIO(conn *ws.Conn) error {
	go func(){
		for {
			// Write
			err := conn.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
			if err != nil {
				log.Println(err)
			}
	
			// Read
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("%s\n", msg)
		}		
	}
}
