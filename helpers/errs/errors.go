package errs

import "errors"

var (
	WS_CLIENT_CONNECTION_NOT_FOUND = errors.New("Client connection not found.")
	WS_ROOM_NOT_FOUND              = errors.New("Websocket room not found.")
	WS_EVENT_NOT_FOUND             = errors.New("Event not found.")
)

var (
	ERR_USER_NOT_FOUND            = errors.New("user data not found")
	ERR_UNAUTHENTICATED           = errors.New("unauthenticated")
	ERR_NO_AVAILABLE_PLAYERS      = errors.New("no available players currently.")
	ERR_ERROR_TIME_CONTROL_PARSE  = errors.New("failed to parse time control format")
	ERR_INVALID_GAME_TYPE         = errors.New("invalid game type")
	ERR_REDIS_DATA_NOT_FOUND      = errors.New("redis data not found")
	ERR_PLAYER_IN_GAME            = errors.New("player is already on an existing game")
	ERR_PLAYER_IN_POOL            = errors.New("player is under matchmaking")
	ERR_ACTIVE_GAME_NOT_FOUND     = errors.New("player is not in an active game")
	ERR_INVALID_MOVE              = errors.New("invalid move")
	ERR_GAME_DATA_NOT_FOUND       = errors.New("game data not found")
	ERR_GAME_SKILL_DATA_NOT_FOUND = errors.New("game skill data not found")
	ERR_UNAVAILABLE_SKILL_COUNT   = errors.New("skill cannot be used anymore")
)
