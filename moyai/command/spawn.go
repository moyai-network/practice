package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/practice/moyai/game/lobby"
)

type Spawn struct{}

func (Spawn) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	if lobby.Contains(p) {
		p.Teleport(p.World().Spawn().Vec3Middle())
		return
	}
	lobby.AddPlayer(p)
}
