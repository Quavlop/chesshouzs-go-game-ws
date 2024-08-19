package models

import (
	"github.com/google/uuid"
)

type ExecuteSkillMessage struct {
	State          string    `json:"state"`
	GameId         uuid.UUID `json:"gameId"`
	ExecutorUserId uuid.UUID `json:"executorUserId"`
	SkillId        uuid.UUID `json:"skillId"`
}

type KafkaQueueMessage struct {
	Value ExecuteSkillMessage `json:"value"`
}
