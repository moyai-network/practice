package form

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type settings struct{}

func NewSettings() form.Form {
	return form.NewMenu(settings{}, text.Colourf("<dark-red>» <red>Settings</red> «</dark-red>")).WithButtons(
		form.NewButton("Display", ""),
		//form.NewButton("Visual", ""),
		//form.NewButton("Gameplay", ""),
		form.NewButton("Privacy", ""),
		//form.NewButton("Matchmaking", ""),
		//form.NewButton("Advanced", ""),
	)
}

func (settings) Submit(sub form.Submitter, btn form.Button) {
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
		p.SendForm(NewGameplay(s.u))*/
	case "Privacy":
		p.SendForm(NewPrivacy(p))
		/*case "Matchmaking":
			p.SendForm(NewMatchmaking(s.u))
		case "Advanced":
			p.SendForm(NewAdvanced(s.u))
		*/
	}
}
