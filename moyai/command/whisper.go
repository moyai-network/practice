package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/carrot"
	"github.com/moyai-network/carrot/lang"
	"github.com/moyai-network/carrot/role"
	"github.com/moyai-network/practice/moyai/data"
	"github.com/moyai-network/practice/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strings"
)

// Whisper is a command that allows a player to send a private message to another player.
type Whisper struct {
	Target  []cmd.Target `cmd:"target"`
	Message cmd.Varargs  `cmd:"message"`
}

// Run ...
func (w Whisper) Run(s cmd.Source, o *cmd.Output) {
	p := s.(*player.Player)
	l := p.Locale()

	u, ok := data.LoadUser(p.Name())
	if !ok {
		// The user somehow left in the middle of this, so just stop in our tracks.
		return
	}
	if !u.Settings.Privacy.PrivateMessages {
		o.Error(lang.Translatef(l, "user.whisper.disabled"))
		return
	}
	msg := strings.TrimSpace(string(w.Message))
	if len(msg) <= 0 {
		o.Error(lang.Translatef(l, "message.empty"))
		return
	}
	if len(w.Target) > 1 {
		o.Error(lang.Translatef(l, "command.targets.exceed"))
		return
	}

	tP, ok := w.Target[0].(*player.Player)
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	t, ok := data.LoadUser(tP.Name())
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	if !t.Settings.Privacy.PrivateMessages {
		o.Error(lang.Translatef(l, "target.whisper.disabled"))
		return
	}

	str := "<white>%s</white>"
	ur, tr := u.Roles.Highest(), t.Roles.Highest()

	uTag, uMsg := text.Colourf(str, u.DisplayName), text.Colourf(str, msg)
	tTag, tMsg := text.Colourf(str, t.DisplayName), text.Colourf(str, msg)
	if _, ok := ur.(role.Default); !ok {
		uMsg, uTag = tr.Colour(msg), ur.Colour(u.DisplayName)
	}
	if _, ok := tr.(role.Default); !ok {
		tMsg, tTag = ur.Colour(msg), tr.Colour(t.DisplayName)
	}

	th := tP.Handler().(user.UserHandler).UserHandler()
	th.SetLastMessageFrom(p)
	carrot.SendCustomSound(tP, "random.orb", 1, 1, false)
	p.Message(lang.Translatef(l, "command.whisper.to", tTag, tMsg))
	tP.Message(lang.Translatef(l, "command.whisper.from", uTag, uMsg))
}

// Allow ...
func (Whisper) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}
