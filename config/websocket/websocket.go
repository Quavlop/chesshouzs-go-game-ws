package websocket

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
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

func initConnection(c echo.Context, connectionList *Connections, isGuest bool) (*ws.Conn, error) {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		helpers.WriteOutLog(err.Error())
		helpers.WriteErrLog(err.Error())
		return conn, err
	}

	token, ok := c.Get("user").(models.User)
	if ok {
		connectionList.addConnection(token.ID.String(), conn)
	} else if isGuest {
		guestID := "G:" + uuid.New().String()
		c.Set("guest", guestID)
		connectionList.addConnection(guestID, conn)
	}

	helpers.WriteOutLog("[WEBSOCKET] CONNECTION ESTABLISHED : \"" + c.Request().RemoteAddr + " | " + c.Request().Host + " | " + time.Now().Format(os.Getenv("TIME_FORMAT")) + "\"")
	return conn, nil
}

func NewWebSocketHandler(e *echo.Echo, controller *controllers.Controller, connectionList *Connections) {
	e.GET("/ws", func(c echo.Context) error {

		// handle guest
		var isGuest bool
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			authHeader = c.QueryParam("sid")
			if authHeader == "" {
				isGuest = true
			}
		}

		if !isGuest {
			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method: " + token.Header["alg"].(string))
				}
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
			if err != nil {
				isGuest = true
			} else {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					userID := claims["id"].(string)
					user, err := controller.Repository.GetUserDataByID(userID)
					if err != nil {
						isGuest = true
					} else {
						c.Set("user", user)
					}
				}
			}

		}

		conn, err := initConnection(c, connectionList, isGuest)
		if err != nil {
			helpers.WriteOutLog("[WEBSOCKET] FAILED TO INITIALIZE CONNECTION : " + err.Error())
			return err
		}
		defer conn.Close()

		var sid string
		if !isGuest {
			sid = c.QueryParams().Get("sid")
		}
		return handleIO(c, controller, conn, sid, connectionList)
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
			err = controller.WebSocketService.CleanMatchupState(c, c.Get("user").(models.User))
			if err != nil {
				helpers.WriteErrLog("[WEBSOCKET] MATCH STATE CLEAN UP FAIL : " + err.Error())
				break
			}
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
		constants.WS_EVENT_INIT_MATCHMAKING:        controller.HandleMatchmaking,
		constants.WS_EVENT_RESTORE_GAME_CONNECTION: nil,
	}

	handler, eventExists := eventHandler[message.Event]
	if !eventExists {
		err := errs.WS_EVENT_NOT_FOUND
		return models.WebSocketResponse{
			Status: constants.WS_SERVER_RESPONSE_ERROR,
			Data:   err.Error(),
		}, err
	}

	// auth middleware for websocket events
	eventTypeParse := strings.Split(message.Event, ":")
	if len(eventTypeParse) == 1 {
		guestToken, _ := c.Get("guest").(string)
		if guestToken != "" {
			errAuth := errs.ERR_UNAUTHENTICATED
			return models.WebSocketResponse{
				Status: constants.WS_SERVER_RESPONSE_ERROR,
				Data:   errAuth.Error(),
			}, errAuth
		}
	}

	c.Set("ws-event", message.Event)

	user, ok := c.Get("user").(models.User)
	if !ok {
		user = models.User{}
	}

	// handler()'s argument contains the connection data which belongs to the request initiator.
	response, err := handler(models.WebSocketClientData{
		Connection: conn,
		Token:      token,
		Event:      message.Event,
		Context:    &c,
		User:       user,
		Data:       message.Data,
	})
	if err != nil {
		return response, err
	}
	return response, nil
}
