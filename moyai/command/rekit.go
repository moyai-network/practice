package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/practice/moyai/game/ffa"
	"github.com/moyai-network/practice/moyai/game/kit"
)

type ReKit struct{}

func (ReKit) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	h, ok := p.Handler().(*ffa.Handler)
	if !ok {
		return
	}

	kit.Apply(h.Game().Kit(), p)
}
