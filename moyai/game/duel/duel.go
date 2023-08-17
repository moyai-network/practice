package duel

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/practice/moyai/game"
	"time"
)

var lobby func(p *player.Player)

func InitializeLobby(f func(*player.Player)) {
	lobby = f
}

func init() {
	t := time.NewTicker(time.Second)
	go func() {
		for range t.C {
			for _, g := range game.Games() {
				q := game.Queued(g, false)
				rq := game.Queued(g, true)
				if len(q) >= 2 {
					NewMatch(q[0], q[1], g, false).Start()
				}
				if len(rq) >= 2 {
					NewMatch(rq[0], rq[1], g, true).Start()
				}
			}
		}
	}()
}

var pairs = [][]world.Block{
	{
		block.Concrete{Colour: item.ColourLime()},
		block.Concrete{Colour: item.ColourGreen()},
	},
	{
		block.Concrete{Colour: item.ColourLightBlue()},
		block.Concrete{Colour: item.ColourBlue()},
	},
	{
		block.Concrete{Colour: item.ColourLightGrey()},
		block.Concrete{Colour: item.ColourBlack()},
	},
}
