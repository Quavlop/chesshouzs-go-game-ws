package tests

import (
	"github.com/google/uuid"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

func GenerateUserStub() models.User {
	userID, err := uuid.Parse("4c897f19-aa3e-43df-8ec6-369d70c5dc7d")
	if err != nil {
		return models.User{}
	}
	return models.User{
		ID:        userID,
		Username:  "IngenHouzs",
		Email:     "farreldinarta@gmail.com",
		EloPoints: 72,
	}
}
