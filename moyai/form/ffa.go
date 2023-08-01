package form

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/moyai-network/moose"
	"github.com/moyai-network/practice/moyai/game"
	"github.com/moyai-network/practice/moyai/game/ffa"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type FFA struct {
	lobby func(*player.Player)
}

func NewFFA(lobby func(*player.Player)) form.Form {
	var buttons []form.Button
	m := form.NewMenu(FFA{
		lobby: lobby,
	}, moose.GlyphFont("FFA", item.ColourOrange()))
	for _, g := range game.Games() {
		buttons = append(buttons, form.NewButton(text.Colourf("<purple>%s</purple>", g.Name()), g.Texture()))
	}
	return m.WithButtons(buttons...)
}

func (f FFA) Submit(sub form.Submitter, btn form.Button) {
	p, ok := sub.(*player.Player)
	if !ok {
		return
	}

	g := game.ByName(btn.Text)

	ffa.AddPlayer(p, g, f.lobby)
}
