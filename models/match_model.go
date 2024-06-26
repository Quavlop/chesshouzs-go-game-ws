package models

type HandleMatchmakingParams struct {
	Type        string `json:"type"`
	TimeControl string `json:"time_control"`
}

type HandleMatchmakingResponse struct {
	ID       string                    `json:"id"`
	Opponent PlayerMatchmakingResponse `json:"opponent"`
}

type HandleRecoverMatchSocketConnectionParams struct{}

type HandleRecoverMatchSocketConnectionResponse struct{}

type GameData struct {
	ID string `json:"id"`
}

type PlayerGameState struct {
	ID          string
	Type        string
	TimeControl string
}
