package repositories

import (
	"github.com/google/uuid"
	"ingenhouzs.com/chesshouzs/go-game/helpers/errs"
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
		return data, errs.ERR_GAME_DATA_NOT_FOUND
	}

	return data, nil
}

func (r *Repository) GetGameSkills(params models.GameSkill) ([]models.GameSkill, error) {
	var data []models.GameSkill

	db := r.postgres.Table("game_skill gs").Select(`
		gs.id, 
		gs.name, 
		gs.description,
		gs.for_self, 
		gs.for_enemy, 
		gs.radius_x, 
		gs.radius_y, 
		gs.auto_trigger, 
		gs.duration, 
		gs.usage_count
	`)

	if params.ID != uuid.Nil {
		db = db.Where("gs.id = ?", params.ID.String())
	}

	result := db.Find(&data)
	if result.Error != nil {
		return data, result.Error
	}

	if result.RecordNotFound() {
		return data, errs.ERR_GAME_SKILL_DATA_NOT_FOUND
	}

	return data, nil
}
