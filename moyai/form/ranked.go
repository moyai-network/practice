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

type Ranked struct{}

func NewRanked() form.Form {
	var buttons []form.Button
	m := form.NewMenu(Ranked{}, text.Colourf("<redstone>» <red>Ranked Queue</red> «</redstone>"))
	for _, g := range game.Games() {
		buttons = append(buttons, form.NewButton(text.Colourf("<dark-grey>%s</dark-grey>\n<grey>%d Queuing</grey>", g.Name(), len(game.Queued(g, true))), g.Texture()))
	}
	return m.WithBody(text.Colourf("<redstone>»</redstone> Welcome to the <red>Ranked</red> form. You may choose a game mode.")).WithButtons(buttons...)
}

func (Ranked) Submit(sub form.Submitter, btn form.Button) {
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

	game.Queue(p, g, true)
	kit.Apply(kit.Queue{}, p)
}
