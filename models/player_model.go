package models

type Player struct {
	EloPoints int32
}
type FilterEligibleOpponentParams struct {
	Filter PoolParams
	Client PlayerPool
}
type FilterEligibleOpponentResponse struct {
	Player PlayerPool
}
