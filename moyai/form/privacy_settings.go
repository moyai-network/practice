package form

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/moyai-network/practice/moyai/data"
)

// privacy is the form that handles the modification of privacy settings.
type privacy struct {
	// PrivateMessages is a dropdown that allows the user to enable or disable private messages from others.
	PrivateMessages form.Dropdown
	// PublicStatistics is a dropdown that allows the user to enable or disable public statistics.
	PublicStatistics form.Dropdown
	// DuelRequests is a dropdown that allows the user to enable or disable duel requests from others.
	DuelRequests form.Dropdown
	// p is the player that is using the form.
	p *player.Player
}

// NewPrivacy creates a new form for the player to modify their privacy settings.
func NewPrivacy(p *player.Player) form.Form {
	u, _ := data.LoadUser(p.Name())
	s := u.Settings
	return form.New(privacy{
		PrivateMessages: newToggleDropdown("Allow others to private message me:", s.Privacy.PrivateMessages),
		p:               p,
	}, "Privacy Settings")
}

// Submit ...
func (d privacy) Submit(form.Submitter) {
	u, _ := data.LoadUser(d.p.XUID())
	s := u.Settings
	s.Privacy.PrivateMessages = indexBool(d.PrivateMessages)
	_ = data.SaveUser(u)
	d.p.SendForm(NewPrivacy(d.p))
}

// Close ...
func (d privacy) Close(form.Submitter) {
	d.p.SendForm(NewSettings())
}
