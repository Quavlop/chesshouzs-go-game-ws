package tests

import (
	"context"

	"ingenhouzs.com/chesshouzs/go-game/interfaces"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

// test construction
type Fields struct {
	Repository interfaces.Repository
}

type Args struct {
	Ctx    context.Context
	Params map[string]interface{}
	Client models.WebSocketClientData
}

type Errs struct {
	ExpectErr bool
}

type Test struct {
	Fields
	Args
	Errs
}
