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
	wsRoom.clients.GetLock().Lock()
	defer wsRoom.clients.GetLock().Unlock()
	return wsRoom.clients.GetMap()
}

func (wsRoom *GameRoom) IsClientInRoom(token string) bool {
	wsRoom.clients.GetLock().Lock()
	defer wsRoom.clients.GetLock().Unlock()
	data, exists := wsRoom.GetClients()[token]
	return data && exists
}

func (wsRoom *GameRoom) AddClient(token string) {
	wsRoom.clients.GetLock().Lock()
	defer wsRoom.clients.GetLock().Unlock()
	if wsRoom.GetClients() == nil {
		wsRoom.clients.NewMap()
	}
	wsRoom.GetClients()[token] = true
}

func (wsRoom *GameRoom) RemoveClient(token string) {
	wsRoom.clients.GetLock()
	delete(wsRoom.GetClients(), token)
	wsRoom.clients.GetLock().Unlock()
}
