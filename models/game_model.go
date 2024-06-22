package models

import (
	"time"

	"github.com/google/uuid"
)

type PoolParams struct {
	ID          uuid.UUID
	Type        string
	TimeControl string
	UpperBound  int32
	LowerBound  int32
}

type PlayerPool struct {
	JoinTime time.Time
	User     User
}

type PlayerMatchmakingResponse struct {
	JoinTime string
	User     User
}

type PlayerPoolParams struct {
	PoolParams
	User   User
	Player PlayerPool
}

type GameTypeVariant struct {
	ID        uuid.UUID
	Name      string
	Duration  int32
	Increment int32
}

type MoveCache struct {
	ID   uuid.UUID
	Turn bool
}

type GameActiveData struct {
	ID                uuid.UUID
	WhitePlayerID     uuid.UUID
	BlackPlayerID     uuid.UUID
	GameTypeVariantID uuid.UUID
	MovesCacheRef     uuid.UUID
	Moves             string
	IsDone            bool
	WinnerPlayerID    uuid.UUID
	StartTime         string
	EndTime           string
}

type EloBounds struct {
	Upper int32
	Lower int32
}
