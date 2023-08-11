package duel

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/carrot/sets"
	"github.com/moyai-network/practice/moyai/game"
	"time"
)

var lobby func(p *player.Player)

func InitializeLobby(f func(*player.Player)) {
	lobby = f
}

var queued = sets.New[*player.Player]()

func init() {
	t := time.NewTicker(time.Second)
	go func() {
		for range t.C {
			v := queued.Values()
			if len(v) < 2 {
				continue
			}
			Start(v[0], v[1], game.NoDebuff())
		}
	}()
}

func Queue(p *player.Player, g game.Game) {
	queued.Add(p)
}

func Queued(p *player.Player) bool {
	return queued.Contains(p)
}

func UnQueue(p *player.Player) {
	queued.Delete(p)
}
