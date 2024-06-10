package repositories

import (
	"errors"

	"ingenhouzs.com/chesshouzs/go-game/models"
)

func (r *Repository) GetUserDataByID(id string) (models.User, error) {
	var user models.User
	db := r.postgres.Table("users u").Select("*")
	db = db.Where("u.id = ?", id)

	result := db.Find(&user)

	if result.Error != nil {
		return user, result.Error
	}

	if result.RecordNotFound() {
		return user, errors.New("user data not found")
	}

	return user, nil
}
