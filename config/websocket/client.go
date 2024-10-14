package websocket

import (
	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
	"ingenhouzs.com/chesshouzs/go-game/constants"
	"ingenhouzs.com/chesshouzs/go-game/helpers/errs"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

type Connections struct {
	clients   models.SafeMapClient
	gameRooms models.SafeMapGameRoom
}

func (c *Connections) Init() {
	c.clients.NewMap()
	c.gameRooms.NewMap()
	// c.gameRooms[0].AddClient()
}

func (c *Connections) GetConnections() map[string]*models.WebSocketClientConnection {
	c.clients.GetLock().RLock()
	defer c.clients.GetLock().RUnlock()
	return c.clients.GetMap()
}

func (c *Connections) GetRooms() map[string]*models.GameRoom {
	c.gameRooms.GetLock().RLock()
	defer c.gameRooms.GetLock().RUnlock()
	return c.gameRooms.GetMap()
}

func (c *Connections) GetRoomByID(id string) *models.GameRoom {
	c.gameRooms.GetLock().RLock()
	defer c.gameRooms.GetLock().RUnlock()
	rooms := c.gameRooms.GetMap()
	return rooms[id]
}

func (c *Connections) GetClientConnection(token string) *models.WebSocketClientConnection {
	c.clients.GetLock().RLock()
	defer c.clients.GetLock().RUnlock()
	return c.clients.GetMap()[token]
}

func (c *Connections) GetClientActiveRooms(token string) []models.GameRoom {
	var rooms []models.GameRoom

	c.gameRooms.GetLock().RLock()
	defer c.gameRooms.GetLock().RUnlock()

	for _, room := range c.gameRooms.GetMap() {
		if room.IsClientInRoom(token) {
			rooms = append(rooms, models.GameRoom{
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
	c.clients.GetLock().Lock()
	defer c.clients.GetLock().Unlock()
	c.clients.GetMap()[token] = &models.WebSocketClientConnection{
		Connection: conn,
		Token:      token,
	}
}

func (c *Connections) deleteConnection(token string, conn *ws.Conn) {
	// delete from global connections
	c.clients.GetLock().Lock()
	delete(c.clients.GetMap(), token)
	// c.clients.GetLock().Unlock()

	// delete from room connections
	// c.gameRooms.GetLock().Lock()
	rooms := c.GetRooms()
	for _, room := range rooms {
		roomClients := room.GetClients()
		delete(roomClients, token)

		if len(roomClients) <= 0 {
			delete(rooms, room.GetRoomID())
		}
	}

	c.clients.GetLock().Unlock()
}

func (c *Connections) CreateRoom(params *models.GameRoom, roomID string) *models.GameRoom {
	c.gameRooms.NewMap()
	c.gameRooms.GetLock().Lock()
	defer c.gameRooms.GetLock().Unlock()

	connectionPoolRoom := c.gameRooms.GetMap()

	if roomID == "" {
		id := uuid.New()
		connectionPoolRoom[id.String()] = params
		connectionPoolRoom[id.String()].SetRoomID(id.String())
		return connectionPoolRoom[id.String()]
	}

	connectionPoolRoom[roomID] = params
	connectionPoolRoom[roomID].SetRoomID(roomID)
	return connectionPoolRoom[roomID]
}

func (c *Connections) EmitOneOnOne(params models.WebSocketChannel) error {
	// sourceClient := c.IsClientActive(params.Source)
	targetClient := c.IsClientActive(params.TargetClient)
	if targetClient == nil {
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
