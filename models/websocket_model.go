package models

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

type SafeMap struct {
	mtx *sync.RWMutex
	m   interface{}
}

type SafeMapInterface interface {
	GetLock() *sync.RWMutex
}

type SafeMapClient struct {
	SafeMap
}

type SafeMapGameRoom struct {
	SafeMap
}

type SafeMapRoomClient struct {
	SafeMap
}

func (sm *SafeMap) GetLock() *sync.RWMutex {
	return sm.mtx
}

func (sm *SafeMapClient) GetMap() map[string]*WebSocketClientConnection {
	res, ok := sm.m.(map[string]*WebSocketClientConnection)
	if !ok {
		return nil
	}
	return res
}

func (sm *SafeMapGameRoom) GetMap() map[string]*GameRoom {
	res, ok := sm.m.(map[string]*GameRoom)
	if !ok {
		return nil
	}
	return res
}

func (sm *SafeMapRoomClient) GetMap() map[string]bool {
	res, ok := sm.m.(map[string]bool)
	if !ok {
		return nil
	}
	return res
}

func (sm *SafeMapClient) NewMap() {
	sm.m = make(map[string]*WebSocketClientConnection)
	sm.mtx = &sync.RWMutex{}
}

func (sm *SafeMapGameRoom) NewMap() {
	sm.m = make(map[string]*GameRoom)
	sm.mtx = &sync.RWMutex{}
}

func (sm *SafeMapRoomClient) NewMap() {
	sm.m = make(map[string]bool)
	sm.mtx = &sync.RWMutex{}
}

type WebSocketClientConnection struct {
	Connection *websocket.Conn
	Token      string
}

type WebSocketClientData struct {
	Connection *websocket.Conn
	Token      string
	Event      string
	Context    *echo.Context
	User       User
	Data       interface{}
}

type WebSocketClientMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type WebSocketChannel struct {
	Source       string
	TargetClient string
	TargetRoom   string
	Event        string
	Data         interface{}
}

type WebSocketResponse struct {
	Status string      `json:"status,omitempty"`
	Event  string      `json:"event,omitempty"`
	Data   interface{} `json:"data"`
}

type GameRoom struct {
	id   string
	Type string // METADATA
	// etc ...

	// key : user's session token
	// [BETA] value : boolean value to show if user is in the room ()
	clients SafeMapRoomClient
}
