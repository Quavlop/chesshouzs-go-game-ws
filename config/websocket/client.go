package websocket

import (
	ws "github.com/gorilla/websocket"
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

func (c *Connections) addConnection(token string, conn *ws.Conn) {
	c.clients[token] = &models.WebSocketClientConnection{
		Connection: conn,
	}
}

func (c *Connections) deleteConnection(token string, conn *ws.Conn) {
	delete(c.clients, token) // delete from global connections
	// delete from room connections
	// ...
}
