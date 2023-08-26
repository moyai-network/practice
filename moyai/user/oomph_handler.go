package user

import (
	"github.com/df-mc/dragonfly/server/player/chat"
	"golang.org/x/text/language"
	"strings"

	"github.com/df-mc/dragonfly/server/event"
	p "github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/carrot/lang"
	"github.com/oomph-ac/oomph/check"
	"github.com/oomph-ac/oomph/player"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
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
	p, ok := Lookup(h.p.Name())
	if !ok {
		return
	}
	if h, ok := p.Handler().(*Handler); ok {
		h.History()[ch] = ch.Violations()
	}
	Broadcast("staff.alert",
		h.p.Name(),
		name,
		variant,
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

	_, _ = chat.Global.WriteString(lang.Translatef(language.English, "oomph.kick.broadcast", n+v))
}

func (h *OomphHandler) HandleClientPacket(ctx *event.Context, pk packet.Packet) {
	p, ok := Lookup(h.p.Name())
	if !ok {
		return
	}
	if h, ok := p.Handler().(replayHandler); ok {
		switch pk := pk.(type) {
		case *packet.PlayerAuthInput:
			//fmt.Println(p.Name())
			h.AddReplayAction(p, pk)
		}
	}
	if !ok {
		return
	}
}

type replayHandler interface {
	AddReplayAction(*p.Player, packet.Packet)
}
