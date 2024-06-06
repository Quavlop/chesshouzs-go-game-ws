package tests

import (
	"time"

	"ingenhouzs.com/chesshouzs/go-game/models"
)

func GenerateWebSocketClientData(eloPoints int32) models.WebSocketClientData {
	return models.WebSocketClientData{
		Data: models.Player{
			EloPoints: eloPoints,
		},
	}
}

func GeneratePlayerPool() []models.PlayerPool {
	baseTime, err := time.Parse("2006-01-02 15:04:05", "2024-01-01 10:00:00")
	if err != nil {
		return []models.PlayerPool{}
	}
	return []models.PlayerPool{
		{
			EloPoints: 800,
			JoinTime:  baseTime.Add(3 * time.Minute),
		},
		{
			EloPoints: 840,
			JoinTime:  baseTime.Add(2 * time.Minute),
		},
		{
			EloPoints: 800,
			JoinTime:  baseTime.Add(4 * time.Minute),
		},
		{
			EloPoints: 700,
			JoinTime:  baseTime.Add(5 * time.Minute),
		},
	}
}
