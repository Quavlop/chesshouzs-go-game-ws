package interfaces

type Repository interface {
	GameRepository
}

type GameRepository interface {
	GetUnderMatchmakingPlayers() string
}
