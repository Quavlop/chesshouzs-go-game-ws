package models

import (
	"time"

	"github.com/google/uuid"
)

type PoolParams struct {
	ID          uuid.UUID
	Type        string
	TimeControl string
}

type PlayerPool struct {
	JoinTime time.Time
	User     User
}

type PlayerPoolParams struct {
	PoolParams
	User
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

type InsertGameParams struct {
	WhitePlayerID     uuid.UUID
	BlackPlayerID     uuid.UUID
	GameTypeVariantID uuid.UUID
	MovesCacheRef     uuid.UUID
}
