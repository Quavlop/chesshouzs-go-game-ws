package models

type HandleMatchmakingParams struct {
	Type        string `json:"type"`
	TimeControl string `json:"time_control"`
}

type HandleMatchmakingResponse struct {
	// BETa
	ID string `json:"id"`
}
