package user

import (
	"strings"

	"github.com/df-mc/dragonfly/server/event"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/carrot/lang"
	"github.com/oomph-ac/oomph/check"
	"github.com/oomph-ac/oomph/player"
	"github.com/oomph-ac/oomph/utils"
	"github.com/unickorn/strutils"
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

func (h *OomphHandler) HandleFlag(ctx *event.Context, ch check.Check, params map[string]any, _ *bool) {
	name, variant := ch.Name()
	Broadcast("oomph.staff.alert",
		h.p.Name(),
		name,
		variant,
		utils.PrettyParameters(params, true),
		mgl64.Round(ch.Violations(), 2),
	)
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
