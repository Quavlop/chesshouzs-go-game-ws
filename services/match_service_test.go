package services

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"ingenhouzs.com/chesshouzs/go-game/constants"
	"ingenhouzs.com/chesshouzs/go-game/mocks"
	"ingenhouzs.com/chesshouzs/go-game/models"
	"ingenhouzs.com/chesshouzs/go-game/tests"
)

func TestHandleMatchmaking(t *testing.T) {
}

func TestFilterEligibleOpponent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// helpers
	timeFormat := "2006-01-02 15:04:05"

	// mock arguments
	playerPool := tests.GeneratePlayerPool() // ini hasil dri GetUnderMatchMakingPlayers
	client := tests.GenerateWebSocketClientData(900)

	baseTime, err := time.Parse(timeFormat, "2024-01-01 10:00:00")
	if err != nil {
		t.Errorf("Failed to parse time format : " + err.Error())
		return
	}

	t.Run("POSITIVE CASE : Success", func(t *testing.T) {
		mockRepository := mocks.NewMockRepository(ctrl)
		mockWebSocketService := mocks.NewMockWebsocketService(ctrl)

		params := make(map[string]interface{})
		params["user"] = tests.GenerateUserStub()
		params["filter_eligible_opponent_params"] = models.FilterEligibleOpponentParams{
			Filter: models.PoolParams{
				Type:        constants.GAME_RAPID_TYPE,
				TimeControl: constants.GAME_60_1_TIME_CONTROL,
			},
			Client: models.PlayerPool{
				User: params["user"].(models.User),
			},
		}
		test := tests.Test{
			Fields: tests.Fields{
				Repository: mockRepository,
			},
			Args: tests.Args{
				Ctx:    context.TODO(),
				Params: params,
				Client: client,
			},
			Errs: tests.Errs{
				ExpectErr: false,
			},
		}

		mockRepository.
			EXPECT().
			GetUserDataByID(gomock.Any()).
			Return(params["user"].(models.User), nil).
			AnyTimes()

		// get under matchmaking player
		mockRepository.
			EXPECT().
			GetUnderMatchmakingPlayers(models.PoolParams{
				Type:        test.Args.Params["filter_eligible_opponent_params"].(models.FilterEligibleOpponentParams).Filter.Type,
				TimeControl: test.Args.Params["filter_eligible_opponent_params"].(models.FilterEligibleOpponentParams).Filter.TimeControl,
			}).
			Return(playerPool, nil)

		// filter out opps
		mockWebSocketService.
			EXPECT().
			FilterOutOpponents(client, playerPool).
			Return([]models.PlayerPool{
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
			}, nil)

		// sort player pool
		mockWebSocketService.
			EXPECT().
			SortPlayerPool(client, []models.PlayerPool{
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
			}).
			Return([]models.PlayerPool{
				{
					User: models.User{
						EloPoints: 800,
					},
					JoinTime: baseTime.Add(3 * time.Minute),
				},
				{
					User: models.User{
						EloPoints: 800,
					},
					JoinTime: baseTime.Add(4 * time.Minute),
				},
				{
					User: models.User{
						EloPoints: 840,
					},
					JoinTime: baseTime.Add(2 * time.Minute),
				},
			}, nil)

		s := &webSocketService{
			repository:    test.Fields.Repository,
			wsConnections: nil,
			BaseService: &BaseService{
				WebSocketService: mockWebSocketService,
			},
		}

		if result, err := s.FilterEligibleOpponent(test.Args.Client, test.Args.Params["filter_eligible_opponent_params"].(models.FilterEligibleOpponentParams)); err != nil {
			t.Errorf("match_service.HandleMatchmaking() fails : " + err.Error())
			if result.Player != playerPool[0] {
				t.Errorf("match_service.HandleMatchmaking() fails : didn't get the expected player.")
			}
		}
	})
}
