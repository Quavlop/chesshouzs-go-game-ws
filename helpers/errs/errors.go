package errs

import "errors"

var (
	WS_CLIENT_CONNECTION_NOT_FOUND = errors.New("Client connection not found.")
	WS_ROOM_NOT_FOUND              = errors.New("Websocket room not found.")
	WS_EVENT_NOT_FOUND             = errors.New("Event not found.")
)

var (
	ERR_NO_AVAILABLE_PLAYERS     = errors.New("no available players currently.")
	ERR_ERROR_TIME_CONTROL_PARSE = errors.New("failed to parse time control format")
	ERR_INVALID_GAME_TYPE        = errors.New("invalid game type")
	ERR_REDIS_DATA_NOT_FOUND     = errors.New("redis data not found")
	ERR_PLAYER_IN_GAME           = errors.New("player is already on an existing game")
	ERR_PLAYER_IN_POOL           = errors.New("player is under matchmaking")
)
