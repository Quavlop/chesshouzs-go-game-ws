package tests

import (
	"time"

	"ingenhouzs.com/chesshouzs/go-game/helpers"
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
	location := helpers.GetLocalTimeZone()

	baseTime, err := time.ParseInLocation("2006-01-02 15:04:05", "2024-01-01 10:00:00", location)
	if err != nil {
		return []models.PlayerPool{}
	}
	return []models.PlayerPool{
		{
			User: models.User{
				EloPoints: 800,
			},
			JoinTime: baseTime.Add(3 * time.Minute),
		},
		{
			User: models.User{
				EloPoints: 840,
			},
			JoinTime: baseTime.Add(2 * time.Minute),
		},
		{
			User: models.User{
				EloPoints: 800,
			},
			JoinTime: baseTime.Add(4 * time.Minute),
		},
		{
			User: models.User{
				EloPoints: 700,
			},
			JoinTime: baseTime.Add(5 * time.Minute),
		},
	}
}
