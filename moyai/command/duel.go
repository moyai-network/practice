package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/practice/moyai/form"
	"github.com/moyai-network/practice/moyai/game/duel"
	"github.com/moyai-network/practice/moyai/game/lobby"
	"github.com/moyai-network/practice/moyai/user"
)

type Duel struct {
	Target []cmd.Target
}

type DuelAccept struct {
	Sub  cmd.SubCommand `cmd:"accept"`
	Duel duelRequests
}

func (d Duel) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	_, ok = p.Handler().(user.UserHandler)
	if !ok {
		return
	}

	t, ok := d.Target[0].(*player.Player)
	if !ok {
		return
	}

	if t == p {
		out.Error("You may not duel yourself")
		return
	}
	p.SendForm(form.NewDuel(t))
}

func (d DuelAccept) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	h, ok := p.Handler().(*lobby.Handler)
	if !ok {
		out.Error("Your opponent must be in the lobby in order to accept their duel request.")
		return
	}

	t, ok := user.Lookup(string(d.Duel))
	if !ok {
		return
	}

	_, ok = p.Handler().(*lobby.Handler)
	if !ok {
		out.Error("Your opponent must be in the lobby in order to accept their duel request.")
		return
	}

	g := h.UserHandler().AcceptDuel(t)
	duel.NewMatch(p, t, g).Start()
}

type duelRequests string

func (duelRequests) Type() string {
	return "duels"
}

func (duelRequests) Options(src cmd.Source) []string {
	p, ok := src.(*player.Player)
	if !ok {
		return []string{}
	}
	if h, ok := p.Handler().(user.UserHandler); ok {
		return h.UserHandler().DuelRequests()
	}
	return []string{}
}
