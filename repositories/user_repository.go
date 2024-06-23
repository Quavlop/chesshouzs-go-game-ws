package repositories

import (
	"ingenhouzs.com/chesshouzs/go-game/helpers/errs"
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
		return user, errs.ERR_USER_NOT_FOUND
	}

	return user, nil
}
