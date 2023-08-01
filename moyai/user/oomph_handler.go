package user

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/moyai-network/carrot/lang"
	"github.com/oomph-ac/oomph/check"
	"github.com/oomph-ac/oomph/player"
	"github.com/unickorn/strutils"
	"strings"
)

type OomphHandler struct {
	player.NopHandler
	p *player.Player
}

func NewOomphHandler(p *player.Player) *OomphHandler {
	return &OomphHandler{
		p: p,
	}
}

func (h *OomphHandler) HandlePunishment(ctx *event.Context, ch check.Check, msg *string) {
	ctx.Cancel()
	n, v := ch.Name()
	l := h.p.Locale()
	h.p.Disconnect(strutils.CenterLine(strings.Join([]string{
		lang.Translatef(l, "user.kick.header.oomph"),
		lang.Translatef(l, "user.kick.description", n+v),
	}, "\n")))
}
