package form

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/moyai-network/practice/moyai/game"
	"github.com/moyai-network/practice/moyai/game/ffa"
	"github.com/moyai-network/practice/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"golang.org/x/exp/slices"
	"strings"
)

type FFA struct{}

func NewFFA() form.Form {
	var buttons []form.Button
	playing := map[game.Game]int{}

	for _, u := range user.All() {
		if g, ok := u.Handler().(*ffa.Handler); ok {
			playing[g.Game()]++
		}
	}

	m := form.NewMenu(FFA{}, text.Colourf("<orange>» <black>Free For All</black> «</orange>"))
	for _, g := range game.Games() {
		if !g.FFA() {
			continue
		}
		buttons = append(buttons, form.NewButton(text.Colourf("<dark-grey>%s</dark-grey>\n<grey>%d playing</grey>", g.Name(), playing[g]), g.Texture()))
	}
	return m.WithBody(text.Colourf("<orange>»</orange> Welcome to the <black>FFA</black> form. You may choose a game mode.")).WithButtons(buttons...)
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

	g := game.ByName(strings.Split(btn.Text, "\n")[0])

	ffa.AddPlayer(p, g)
}
