package helpers

import (
	"os"
	"strconv"
	"strings"
	"time"

	"ingenhouzs.com/chesshouzs/go-game/models"
)

func GetTimeoutThreshold(env string) time.Duration {
	timeout := os.Getenv(env)
	timeoutInt, err := strconv.Atoi(timeout)
	if err != nil {
		return 10
	}
	return time.Duration(timeoutInt)
}

func GetPoolKey(params models.PoolParams) string {
	return "pool:" + params.Type + ":" + params.TimeControl
}

func GetGameMoveCacheKey(params models.MoveCache) string {
	return "game_move:" + params.ID.String()
}

func GetPlayerPoolCloneKey(params models.PlayerPoolParams) string {
	return "pool_player:" + params.User.ID.String()
}

func GameNotationBuilder(size int) string {
	var notation string

	firstLine := "rnbnbnkqnbnbnr"
	secondLine := "ppnbpprrppbnpp"
	thirdLine := "pppppppppppppp"

	notation = notation + firstLine + "|" + secondLine + "|" + thirdLine

	for i := 3; i < size-3; i++ {
		notation += "|.............."
	}

	notation = notation + "|" + strings.ToUpper(thirdLine) + "|" + strings.ToUpper(secondLine) + "|" + strings.ToUpper(firstLine)
	return notation
}
