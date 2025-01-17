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
	ID                 uuid.UUID
	Move               string
	Turn               bool
	LastMovement       time.Time
	WhiteTotalDuration int64
	BlackTotalDuration int64
}

type InitMatchSkillStats struct {
	ID           uuid.UUID
	GameSkills   []GameSkill
	GameSkillMap map[string]int
}

type GameActiveData struct {
	ID                uuid.UUID
	WhitePlayerID     uuid.UUID
	BlackPlayerID     uuid.UUID
	GameTypeVariantID uuid.UUID
	RoomID            uuid.UUID
	MovesCacheRef     uuid.UUID
	Moves             string
	IsDone            bool
	WinnerPlayerID    uuid.UUID
	StartTime         string
	EndTime           string

	Duration  int64
	Increment int64
}

type EloBounds struct {
	Upper int32
	Lower int32
}

type GameSkill struct {
	ID               uuid.UUID
	Name             string
	Description      string
	ForSelf          bool
	ForEnemy         bool
	RadiusX          int
	RadiusY          int
	AutoTrigger      bool
	Duration         int
	UsageCount       int
	Type             string
	Permanent        bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	CurrentUserCount int `gorm:"-"`
}

type SkillUsageCount struct {
	ID uuid.UUID
}
