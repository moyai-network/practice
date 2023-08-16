package form

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/moyai-network/carrot"
	"github.com/moyai-network/practice/moyai/game"
	"github.com/moyai-network/practice/moyai/game/ffa"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"golang.org/x/exp/slices"
)

type FFA struct{}

func NewFFA() form.Form {
	var buttons []form.Button
	m := form.NewMenu(FFA{}, carrot.GlyphFont("FFA"))
	for _, g := range game.Games() {
		if !g.FFA() {
			continue
		}
		buttons = append(buttons, form.NewButton(text.Colourf("<purple>%s</purple>", g.Name()), g.Texture()))
	}
	return m.WithButtons(buttons...)
}

func (f FFA) Submit(sub form.Submitter, btn form.Button) {
	p, ok := sub.(*player.Player)
	if !ok {
		return
	}

	for _, gm := range game.Games() {
		if slices.Contains(game.Queued(gm), p) {
			return
		}
	}

	g := game.ByName(btn.Text)

	ffa.AddPlayer(p, g)
}
