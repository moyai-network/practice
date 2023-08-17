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

type Queue struct{}

func NewQueue() form.Form {
	var buttons []form.Button
	m := form.NewMenu(Queue{}, text.Colourf("<orange>» <black>Duel Queue</black> «</orange>"))
	for _, g := range game.Games() {
		buttons = append(buttons, form.NewButton(text.Colourf("<dark-grey>%s</dark-grey>\n<grey>%d in queue</grey>", g.Name(), len(game.Queued(g))), g.Texture()))
	}
	return m.WithBody(text.Colourf("<orange>»</orange> Welcome to the <black>Queue</black> form. You may choose a game mode.")).WithButtons(buttons...)
}

func (q Queue) Submit(sub form.Submitter, btn form.Button) {
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

	game.Queue(p, g)
	kit.Apply(kit.Queue{}, p)
}
