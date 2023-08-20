package form

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/moyai-network/practice/moyai/data"
)

// display is the form that handles the modification of display settings.
type display struct {
	// Scoreboard is a dropdown that allows the user to enable or disable scoreboards.
	Scoreboard form.Dropdown
	// CPS is a dropdown that allows the user to enable or disable the CPS counter.
	//CPS form.Dropdown
}

// NewDisplay creates a new form for the player to modify their display settings.
func NewDisplay(p *player.Player) form.Form {
	u, _ := data.LoadUser(p.Name())
	s := u.Settings
	return form.New(display{
		Scoreboard: newToggleDropdown("Scoreboard:", s.Display.Scoreboard),
	}, "Display Settings")
}

// Submit ...
func (d display) Submit(sub form.Submitter) {
	p, ok := sub.(*player.Player)
	if !ok {
		return
	}

	u, _ := data.LoadUser(p.Name())
	s := u.Settings
	//s.Display.CPS = indexBool(d.CPS)
	s.Display.Scoreboard = indexBool(d.Scoreboard)
	_ = data.SaveUser(u.WithSettings(s))

	if s.Display.Scoreboard {
		if h, ok := p.Handler().(interface{ SendScoreBoard() }); ok {
			h.SendScoreBoard()
		}
	} else if !s.Display.Scoreboard {
		p.RemoveScoreboard()
	}
	p.SendForm(NewDisplay(p))
}

// Close ...
func (d display) Close(sub form.Submitter) {
	p, ok := sub.(*player.Player)
	if !ok {
		return
	}

	p.SendForm(NewSettings())
}