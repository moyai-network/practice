package form

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/moyai-network/practice/moyai/game"
	"github.com/moyai-network/practice/moyai/game/kit"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"golang.org/x/exp/slices"
	"strings"
)

type casualQueue struct{}

func NewCasualQueue() form.Form {
	var buttons []form.Button
	m := form.NewMenu(casualQueue{}, text.Colourf("<dark-red>» <red>Casual Queue</red> «</dark-red>"))
	for _, g := range game.Games() {
		if !g.Duel() {
			continue
		}
		buttons = append(buttons, form.NewButton(text.Colourf("<dark-grey>%s</dark-grey>\n<grey>%d Queuing</grey>", g.Name(), len(game.Queued(g, false))), g.Texture()))
	}
	return m.WithBody(text.Colourf("<dark-red>»</dark-red> Welcome to the <red>Casual Queue</red> form. You may choose a game mode.")).WithButtons(buttons...)
}

func (casualQueue) Submit(sub form.Submitter, btn form.Button) {
	p, ok := sub.(*player.Player)
	if !ok {
		return
	}

	for _, gm := range game.Games() {
		if slices.Contains(game.Queued(gm, false), p) || slices.Contains(game.Queued(gm, true), p) {
			return
		}
	}

	g := game.ByName(strings.Split(btn.Text, "\n")[0])

	game.Queue(p, g, false)
	kit.Apply(kit.Queue{}, p)
}
