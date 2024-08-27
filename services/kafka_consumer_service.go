package services

import (
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"ingenhouzs.com/chesshouzs/go-game/constants"
	"ingenhouzs.com/chesshouzs/go-game/helpers"
	"ingenhouzs.com/chesshouzs/go-game/helpers/errs"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

func (s *kafkaConsumer) ExecuteSkillConsumer(message models.ExecuteSkillMessage) error {

	// get redis state data
	game, err := s.repository.GetPlayerCurrentGameState(message.ExecutorUserId.String())
	if err != nil {
		return err
	}

	// get redis skill data
	skillParam := models.InitMatchSkillStats{
		ID: message.ExecutorUserId,
	}
	skill, err := s.repository.GetPlayerSkillCountUsageData(skillParam)
	if err != nil {
		return err
	}

	skillUsage := skill[message.SkillId.String()]
	if skillUsage <= 0 {
		return errs.ERR_UNAVAILABLE_SKILL_COUNT
	}

	skill[message.SkillId.String()] = skillUsage - 1

	var opponentID uuid.UUID
	if message.ExecutorUserId == game.BlackPlayerID {
		opponentID = game.WhitePlayerID
	} else {
		opponentID = game.BlackPlayerID
	}

	gameRoom := s.wsConnections.GetRoomByID(game.ID.String())

	keys := []string{
		helpers.GetPlayerMatchSkillState(skillParam),
		helpers.GetGameMoveCacheKey(models.MoveCache{
			ID: game.MovesCacheRef,
		}),
	}

	gameMove, err := s.repository.GetMoveCacheIdentifier(models.MoveCache{
		ID: game.MovesCacheRef,
	})
	if err != nil {
		return err
	}

	var newTurn bool
	if gameMove["turn"] == "1" {
		newTurn = false
	} else {
		newTurn = true
	}

	err = s.repository.WithRedisTrx(s.context, keys, func(pipe redis.Pipeliner) error {
		// set new state redis data
		// - skip turn
		// - update new state
		err = s.repository.InsertMoveCacheIdentifier(models.MoveCache{
			ID:   game.MovesCacheRef,
			Move: message.State,
			Turn: newTurn,
		}, pipe)
		if err != nil {
			return err
		}

		// set count skill usage redis data (decrement by one)
		err = s.repository.InsertMatchSkillCount(models.InitMatchSkillStats{
			ID:           message.ExecutorUserId,
			GameSkillMap: skill,
		}, pipe)
		if err != nil {
			return err
		}

		// apply the player state increment here
		err = s.BaseService.WebSocketService.ApplySkillEffects(
			message.GameId,
			message.SkillId,
			message.ExecutorUserId,
			opponentID,
			message.Position,
		)
		if err != nil {
			return err
		}

		// emit event to clients
		s.wsConnections.EmitToRoom(models.WebSocketChannel{
			Source:       message.ExecutorUserId.String(),
			TargetClient: opponentID.String(),
			Event:        constants.WS_EVENT_EMIT_UPDATE_GAME_STATE,
			Data: models.HandleGamePublishActionResponse{
				State: message.State,
				Turn:  message.ExecutorUserId == game.BlackPlayerID,
			},
			TargetRoom: gameRoom.GetRoomID(),
		})

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
