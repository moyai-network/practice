package form

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Settings struct{}

func NewSettings() form.Form {
	return form.NewMenu(Settings{}, text.Colourf("Settings")).WithButtons(
		form.NewButton("Display", ""),
		//form.NewButton("Visual", ""),
		//form.NewButton("Gameplay", ""),
		//form.NewButton("Privacy", ""),
		//form.NewButton("Matchmaking", ""),
		//form.NewButton("Advanced", ""),
	)
}

func (Settings) Submit(sub form.Submitter, btn form.Button) {
	p, ok := sub.(*player.Player)
	if !ok {
		return
	}

	switch btn.Text {
	case "Display":
		p.SendForm(NewDisplay(p))
		/*case "Visual":
			p.SendForm(NewVisual(s.u))
		case "Gameplay":
			p.SendForm(NewGameplay(s.u))
		case "Privacy":
			p.SendForm(NewPrivacy(s.u))
		case "Matchmaking":
			p.SendForm(NewMatchmaking(s.u))
		case "Advanced":
			p.SendForm(NewAdvanced(s.u))
		*/
	}
}