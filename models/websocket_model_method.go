package models

func (wsRoom *GameRoom) GetRoomID() string {
	return wsRoom.id
}

func (wsRoom *GameRoom) GetRoomData() GameRoom {
	return GameRoom{
		Name: wsRoom.Name,
		Type: wsRoom.Type,
	}
}

func (wsRoom *GameRoom) GetClients() map[string]bool {
	return wsRoom.clients
}

func (wsRoom *GameRoom) IsClientInRoom(token string) bool {
	data, exists := wsRoom.clients[token]
	return data && exists
}

func (wsRoom *GameRoom) AddClient(token string) {
	wsRoom.clients[token] = true
}

func (wsRoom *GameRoom) RemoveClient(token string) {
	delete(wsRoom.clients, token)
}
