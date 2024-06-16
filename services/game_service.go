package services

import (
	"os"
	"strconv"

	"ingenhouzs.com/chesshouzs/go-game/models"
)

// func (s *gameRoomService) GetClients() map[string]bool {
// 	return s.room.GetClients()
// }

// func (s *gameRoomService) AddClient(token string) {
// 	s.room.AddClient(token)
// }

// func (s *gameRoomService) RemoveClient(token string) {
// 	s.room.RemoveClient(token)
// }

func (s *httpService) IsValidGameType(params models.GameTypeVariant) (bool, error) {

	gameTypes, err := s.repository.GetGameTypeVariant(models.GameTypeVariant{})
	if err != nil {
		return false, err
	}

	for _, gameType := range gameTypes {
		if gameType.Name == params.Name && gameType.Duration == params.Duration && gameType.Increment == params.Increment {
			return true, nil
		}
	}

	return false, nil
}

func (s *httpService) CalculateEloBounds(params models.User) models.EloBounds {
	threshold, err := strconv.Atoi(os.Getenv("ELO_GAP_THRESHOLD"))
	if err != nil {
		threshold = 150
	}
	return models.EloBounds{
		Upper: params.EloPoints + int32(threshold),
		Lower: params.EloPoints - int32(threshold),
	}
}
