package interfaces

type WebSocketRoom interface {
	GetClients() map[string]bool
	AddClient(token string)
	RemoveClient(token string)
}

// func NewGameRoom() WebSocketRoom {
// 	return &models.GameRoom{}
// }
