package form

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/moyai-network/carrot/lang"
	"github.com/moyai-network/practice/moyai/game"
	"github.com/moyai-network/practice/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Duel struct {
	t *player.Player
}

func NewDuel(t *player.Player) form.Form {
	var buttons []form.Button
	m := form.NewMenu(Duel{t: t}, text.Colourf("<redstone>» <red>Duel %s</red> «</redstone>", t.Name()))
	for _, g := range game.Games() {
		buttons = append(buttons, form.NewButton(text.Colourf("<dark-grey>%s</dark-grey>", g.Name()), g.Texture()))
	}
	return m.WithBody(text.Colourf("<redstone>»</redstone> Welcome to the <red>Duel</red> form. You may choose a game mode")).WithButtons(buttons...)
}

func (d Duel) Submit(sub form.Submitter, btn form.Button) {
	p, ok := sub.(*player.Player)
	if !ok {
		return
	}

	g := game.ByName(btn.Text)
	h, ok := d.t.Handler().(user.UserHandler)
	if !ok {
		return
	}

	h.UserHandler().Duel(p, g)

	p.Message(lang.Translatef(p.Locale(), "duel.request", d.t.Name(), g.Name()))
	d.t.Message(lang.Translatef(p.Locale(), "duel.requested", p.Name(), g.Name()))
}