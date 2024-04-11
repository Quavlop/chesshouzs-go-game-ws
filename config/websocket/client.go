package websocket

import (
	ws "github.com/gorilla/websocket"
	"ingenhouzs.com/chesshouzs/go-game/constants"
	"ingenhouzs.com/chesshouzs/go-game/helpers/errs"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

type Connections struct {
	clients   map[string]*models.WebSocketClientConnection
	gameRooms map[string]*models.GameRoom
}

func (c *Connections) Init() {
	c.clients = make(map[string]*models.WebSocketClientConnection)
	c.gameRooms = make(map[string]*models.GameRoom)
	// c.gameRooms[0].AddClient()
}

func (c *Connections) GetConnections() map[string]*models.WebSocketClientConnection {
	return c.clients
}

func (c *Connections) GetClientConnection(token string) *models.WebSocketClientConnection {
	return c.clients[token]
}

func (c *Connections) GetClientActiveRooms(token string) map[string]models.GameRoom {
	var rooms map[string]models.GameRoom

	for _, room := range c.gameRooms {
		if room.IsClientInRoom(token) {
			rooms[token] = models.GameRoom{
				Name: room.Name,
				Type: room.Type,
			}
		}
	}
	return rooms
}

func (c *Connections) IsClientActive(token string) *models.WebSocketClientConnection {
	client, active := c.clients[token]
	if !active {
		return nil
	}
	return client
}

func (c *Connections) addConnection(token string, conn *ws.Conn) {
	c.clients[token] = &models.WebSocketClientConnection{
		Connection: conn,
		Token:      token,
	}
}

func (c *Connections) deleteConnection(token string, conn *ws.Conn) {
	// delete from global connections
	delete(c.clients, token)

	// delete from room connections
	for _, room := range c.gameRooms {
		delete(room.GetClients(), token)
	}

}

func (c *Connections) EmitOneOnOne(params models.WebSocketChannel) error {
	sourceClient := c.IsClientActive(params.Source)
	targetClient := c.IsClientActive(params.TargetClient)
	if sourceClient == nil || targetClient == nil {
		return errs.WS_CLIENT_CONNECTION_NOT_FOUND
	}

	err := targetClient.Connection.WriteJSON(models.WebSocketResponse{
		Status: constants.WS_SERVER_RESPONSE_SUCCESS,
		Event:  params.Source,
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
	room, exists := c.gameRooms[params.TargetRoom]
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
	for clientID, client := range c.clients {
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
