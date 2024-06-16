package websocket

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	ws "github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"ingenhouzs.com/chesshouzs/go-game/constants"
	"ingenhouzs.com/chesshouzs/go-game/controllers"
	"ingenhouzs.com/chesshouzs/go-game/helpers"
	"ingenhouzs.com/chesshouzs/go-game/helpers/errs"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		if os.Getenv("ENVIRONMENT") == "development" {
			return true
		}
		var allowedOrigins = strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
		origin := r.Header.Get("Origin")
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				return true
			}
		}
		return false
	},
}

func initConnection(c echo.Context, connectionList *Connections) (*ws.Conn, error) {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		helpers.WriteOutLog(err.Error())
		helpers.WriteErrLog(err.Error())
		return conn, err
	}

	token := c.Get("user").(models.User)
	connectionList.addConnection(token.ID.String(), conn)

	helpers.WriteOutLog("[WEBSOCKET] CONNECTION ESTABLISHED : \"" + c.Request().RemoteAddr + " | " + c.Request().Host + " | " + time.Now().Format(os.Getenv("TIME_FORMAT")) + "\"")
	return conn, nil
}

func NewWebSocketHandler(e *echo.Echo, controller *controllers.Controller, connectionList *Connections) {
	e.GET("/ws", func(c echo.Context) error {
		conn, err := initConnection(c, connectionList)
		if err != nil {
			helpers.WriteOutLog("[WEBSOCKET] FAILED TO INITIALIZE CONNECTION : " + err.Error())
			return err
		}
		defer conn.Close()
		token := c.QueryParams().Get("sid")
		return handleIO(c, controller, conn, token, connectionList)
	})
}

func handleIO(c echo.Context, controller *controllers.Controller, conn *ws.Conn, token string, connectionList *Connections) error {

	for {

		var message models.WebSocketClientMessage
		err := conn.ReadJSON(&message)
		if err != nil {
			token := c.Get("user").(models.User).ID.String()
			connectionList.deleteConnection(token, conn)
			helpers.WriteOutLog("[WEBSOCKET] Failed to read message : " + err.Error())
			break
		}
		if message.Event == "" {
			continue
		}
		response, err := handleEvents(c, controller, conn, token, connectionList, message)
		if err != nil {
			helpers.LogErrorCallStack(c, err)
			helpers.WriteOutLog("[WEBSOCKET] Error response : " + err.Error())
			response = helpers.ErrorWebSocketResponseWrap(message.Event, err.Error())
		}

		err = conn.WriteJSON(response)
		if err != nil {
			helpers.LogErrorCallStack(c, err)
			helpers.WriteErrLog("[WEBSOCKET] Failed to send response : " + err.Error())
			continue
		}
	}
	return nil
}

func handleEvents(c echo.Context, controller *controllers.Controller, conn *ws.Conn, token string, connectionList *Connections, message models.WebSocketClientMessage) (models.WebSocketResponse, error) {
	var response models.WebSocketResponse
	var eventHandler = map[string]func(models.WebSocketClientData) (models.WebSocketResponse, error){
		constants.WS_EVENT_INIT_MATCHMAKING: controller.HandleMatchmaking,
	}

	handler, eventExists := eventHandler[message.Event]
	if !eventExists {
		err := errs.WS_EVENT_NOT_FOUND
		return models.WebSocketResponse{
			Status: constants.WS_SERVER_RESPONSE_ERROR,
			Data:   err.Error(),
		}, err
	}

	c.Set("ws-event", message.Event)

	// handler()'s argument contains the connection data which belongs to the request initiator.
	response, err := handler(models.WebSocketClientData{
		Connection: conn,
		Token:      token,
		Event:      message.Event,
		Context:    &c,
		User:       c.Get("user").(models.User),
		Data:       message.Data,
	})
	if err != nil {
		return response, err
	}
	return response, nil
}
