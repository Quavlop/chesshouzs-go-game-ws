package helpers

import (
	"os"
	"strconv"
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
	return "pool_" + params.Type + "_" + params.TimeControl
}
