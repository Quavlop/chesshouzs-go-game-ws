package websocket

import (
	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
	"ingenhouzs.com/chesshouzs/go-game/constants"
	"ingenhouzs.com/chesshouzs/go-game/helpers/errs"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

type Connections struct {
	clients   SafeMapClient
	gameRooms SafeMapGameRoom
}

func (c *Connections) Init() {
	c.clients.NewMap()
	c.gameRooms.NewMap()
	// c.gameRooms[0].AddClient()
}

func (c *Connections) GetConnections() map[string]*models.WebSocketClientConnection {
	return c.clients.GetMap()
}

func (c *Connections) GetRooms() map[string]*models.GameRoom {
	return c.gameRooms.GetMap()
}

func (c *Connections) GetClientConnection(token string) *models.WebSocketClientConnection {
	return c.GetConnections()[token]
}

func (c *Connections) GetClientActiveRooms(token string) []models.GameRoom {
	var rooms []models.GameRoom

	for _, room := range c.gameRooms.GetMap() {
		if room.IsClientInRoom(token) {
			rooms = append(rooms, models.GameRoom{
				Name: room.Name,
				Type: room.Type,
			})
		}
	}
	return rooms
}

func (c *Connections) IsClientInRoom(roomType string, token string) bool {
	for _, room := range c.GetRooms() {
		if room.IsClientInRoom(token) && room.Type == roomType {
			return true
		}
	}
	return false
}

func (c *Connections) IsClientActive(token string) *models.WebSocketClientConnection {
	client, active := c.GetConnections()[token]
	if !active {
		return nil
	}
	return client
}

func (c *Connections) addConnection(token string, conn *ws.Conn) {
	c.GetConnections()[token] = &models.WebSocketClientConnection{
		Connection: conn,
		Token:      token,
	}
}

func (c *Connections) deleteConnection(token string, conn *ws.Conn) {
	// delete from global connections
	delete(c.GetConnections(), token)

	// delete from room connections
	for _, room := range c.GetRooms() {
		delete(room.GetClients(), token)
	}
}

func (c *Connections) CreateRoom(params *models.GameRoom) *models.GameRoom {
	id := uuid.New()
	c.GetRooms()[id.String()] = params
	return c.GetRooms()[id.String()]
}

func (c *Connections) EmitOneOnOne(params models.WebSocketChannel) error {
	sourceClient := c.IsClientActive(params.Source)
	targetClient := c.IsClientActive(params.TargetClient)
	if sourceClient == nil || targetClient == nil {
		return errs.WS_CLIENT_CONNECTION_NOT_FOUND
	}

	err := targetClient.Connection.WriteJSON(models.WebSocketResponse{
		Status: constants.WS_SERVER_RESPONSE_SUCCESS,
		Event:  params.Event,
		Data:   params.Data,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Connections) EmitToRoom(params models.WebSocketChannel) error {
	sourceClient := c.IsClientActive(params.Source)
	if sourceClient == nil {
		return errs.WS_CLIENT_CONNECTION_NOT_FOUND
	}
	room, exists := c.GetRooms()[params.TargetRoom]
	if !exists {
		return errs.WS_ROOM_NOT_FOUND
	}

	for clientID := range room.GetClients() {
		c.EmitOneOnOne(models.WebSocketChannel{
			Source:       params.Source,
			TargetClient: clientID,
			Event:        params.Event,
			Data:         params.Data,
		})
	}

	return nil
}

func (c *Connections) EmitGlobalBroadcast(params models.WebSocketChannel) bool {
	for clientID, client := range c.GetConnections() {
		if client == nil {
			continue
		}
		c.EmitOneOnOne(models.WebSocketChannel{
			Source:       params.Source,
			TargetClient: clientID,
			Event:        params.Event,
			Data:         params.Data,
		})
	}
	return true
}
