package form

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/moyai-network/practice/moyai/data"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type gameplay struct {
	// CriticalEffect is a dropdown that allows the user to enable or disable critical effects.
	CriticalEffect form.Dropdown
	// p is the player that is using the form.
	p *player.Player
}

func NewGameplay(p *player.Player) form.Form {
	u, _ := data.LoadUser(p.Name())
	s := u.Settings
	return form.New(gameplay{
		CriticalEffect: newToggleDropdown("Critical Effect:", s.Gameplay.CriticalEffect),
		p:              p,
	}, text.Colourf("<dark-red>» <red>Gameplay Settings</red> «</dark-red>"))
}

func (g gameplay) Submit(submitter form.Submitter) {
	u, _ := data.LoadUser(g.p.Name())
	s := u.Settings
	s.Gameplay.CriticalEffect = indexBool(g.CriticalEffect)
	_ = data.SaveUser(u.WithSettings(s))

	g.p.SendForm(NewGameplay(g.p))
}
