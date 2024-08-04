package models

type HandleMatchmakingParams struct {
	Type        string `json:"type"`
	TimeControl string `json:"time_control"`
}

type HandleMatchmakingResponse struct {
	ID       string                    `json:"id"`
	Opponent PlayerMatchmakingResponse `json:"opponent"`
}

type HandleConnectMatchSocketConnectionParams struct{}

type HandleConnectMatchSocketConnectionResponse struct{}

type HandleGamePublishActionParams struct {
	State string `json:"state"`
}

type HandleGamePublishActionResponse struct {
	State string `json:"state"`
	Turn  bool   `json:"turn"`
}
type GameData struct {
	ID string `json:"id"`
}

type PlayerGameState struct {
	ID          string
	Type        string
	TimeControl string
}
