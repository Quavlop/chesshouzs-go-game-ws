package models

import "time"

type PoolParams struct {
	Type        string
	TimeControl string
}

type PlayerPool struct {
	EloPoints int32
	JoinTime  time.Time
}
