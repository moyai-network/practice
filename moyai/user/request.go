package user

import (
	"github.com/moyai-network/practice/moyai/game"
	"time"
)

type request struct {
	g          game.Game
	expiration time.Time
}

func newRequest(g game.Game) request {
	return request{
		g:          g,
		expiration: time.Now().Add(5 * time.Minute),
	}
}

func (r request) Expired() bool {
	return r.expiration.Before(time.Now())
}

func (r request) Game() game.Game {
	return r.g
}
