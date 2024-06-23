package repositories

import (
	"errors"

	"ingenhouzs.com/chesshouzs/go-game/models"
)

func (r *Repository) GetGameTypeVariant(params models.GameTypeVariant) ([]models.GameTypeVariant, error) {
	var data []models.GameTypeVariant

	db := r.postgres.Table("game_type gt").Select(`
		gt.name, 
		gtv.id,
		gtv.duration, 
		gtv.increment
	`)

	if params.Duration != 0 {
		db = db.Where("gtv.duration = ?", params.Duration)
	}

	if params.Increment != 0 {
		db = db.Where("gtv.increment = ?", params.Increment)
	}

	db = db.Joins("JOIN game_type_variant gtv ON gt.id = gtv.game_type_id")

	result := db.Find(&data)

	if result.Error != nil {
		return data, result.Error
	}

	if result.RecordNotFound() {
		return data, errors.New("game data not found")
	}

	return data, nil
}
