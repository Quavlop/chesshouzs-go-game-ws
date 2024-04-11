package errs

import "errors"

var (
	WS_CLIENT_CONNECTION_NOT_FOUND = errors.New("Client connection not found.")
	WS_ROOM_NOT_FOUND              = errors.New("Websocket room not found.")
	WS_EVENT_NOT_FOUND             = errors.New("Event not found.")
)
