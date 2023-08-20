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

// Reply is a command that allows a player to reply to their most recent private message.
type Reply struct {
	Message cmd.Varargs `cmd:"message"`
}

// Run ...
func (r Reply) Run(s cmd.Source, o *cmd.Output) {
	p := s.(*player.Player)
	l := p.Locale()

	h := p.Handler().(user.UserHandler).UserHandler()
	u, ok := data.LoadUser(p.Name())
	if !ok {
		return
	}
	if !u.Settings.Privacy.PrivateMessages {
		o.Error(lang.Translatef(l, "user.whisper.disabled"))
		return
	}
	msg := strings.TrimSpace(string(r.Message))
	if len(msg) <= 0 {
		o.Error(lang.Translatef(l, "message.empty"))
		return
	}

	t, ok := h.LastMessageFrom()
	if !ok {
		o.Error(lang.Translatef(l, "command.reply.none"))
		return
	}
	tu, ok := data.LoadUser(t.Name())
	if !ok {
		return
	}
	if !tu.Settings.Privacy.PrivateMessages {
		o.Error(lang.Translatef(l, "target.whisper.disabled"))
		return
	}

	uTag, uMsg := text.Colourf("<white>%s</white>", u.DisplayName), text.Colourf("<white>%s</white>", msg)
	tTag, tMsg := text.Colourf("<white>%s</white>", tu.DisplayName), text.Colourf("<white>%s</white>", msg)
	if _, ok := u.Roles.Highest().(role.Default); !ok {
		uMsg = tu.Roles.Highest().Colour(msg)
		uTag = u.Roles.Highest().Colour(u.DisplayName)
	}
	if _, ok := tu.Roles.Highest().(role.Default); !ok {
		tMsg = u.Roles.Highest().Colour(msg)
		tTag = tu.Roles.Highest().Colour(tu.DisplayName)
	}

	t.Handler().(user.UserHandler).UserHandler().SetLastMessageFrom(t)
	carrot.SendCustomSound(t, "random.orb", 1, 1, false)
	o.Print(lang.Translatef(l, "command.whisper.to", tTag, tMsg))
	t.Message(lang.Translatef(l, "command.whisper.from", uTag, uMsg))
}

// Allow ...
func (Reply) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}
