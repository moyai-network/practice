package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/practice/moyai/game"
	"github.com/moyai-network/practice/moyai/game/duel"
)

type Duel struct {
	Target []cmd.Target
}

func (d Duel) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	t, ok := d.Target[0].(*player.Player)
	if !ok {
		return
	}

	m := duel.NewDuel(p, t, game.NoDebuff())
	m.Start()
}
