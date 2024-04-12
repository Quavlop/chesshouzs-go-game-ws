package models

type PoolParams struct {
	Type        string
	TimeControl string
}

type PlayerPool struct {
	EloPoints int32
}
