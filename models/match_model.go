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

type GetGameTimeDurationParams struct {
	GameID string `json:"game_id"`
}

type GetGameTimeDurationResponse struct {
	White int64 `json:"white"`
	Black int64 `json:"black"`
}

type HandleGamePublishActionResponse struct {
	State              string `json:"state"`
	Turn               bool   `json:"turn"`
	Duration           int64  `json:"duration"`
	WhiteSpentDuration int64  `json:"white_spent_duration"`
	BlackSpentDuration int64  `json:"black_spent_duration"`
}

type GameData struct {
	ID string `json:"id"`
}

type PlayerGameState struct {
	ID          string
	Type        string
	TimeControl string
}

type SkillPosition struct {
	Row int `cql:"row"`
	Col int `cql:"col"`
}

type SkillStatus struct {
	Position     SkillPosition `cql:"position"`
	DurationLeft int           `cql:"duration_left"`
}

type SkillState struct {
	DurationLeft int           `cql:"duration_left"`
	List         []SkillStatus `cql:"list"`
}

type PlayerState struct {
	PlayerID    string                `cql:"player_id"`
	GameID      string                `cql:"game_id"`
	BuffState   map[string]SkillState `cql:"buff_state"`
	DebuffState map[string]SkillState `cql:"debuff_state"`
}
