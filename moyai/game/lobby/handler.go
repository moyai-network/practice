package lobby

import (
	"github.com/moyai-network/carrot/lang"
	"github.com/moyai-network/practice/moyai/form"
	"github.com/moyai-network/practice/moyai/game"
	"github.com/moyai-network/practice/moyai/game/kit"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/carrot"
	"github.com/moyai-network/practice/moyai/data"
	"github.com/moyai-network/practice/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Handler represents the player handler for the lobby.
type Handler struct {
	*user.Handler
	p *player.Player
}

// newHandler returns a new lobby handler.
func newHandler(p *player.Player) *Handler {
	var uHandler *user.Handler
	if uh, ok := p.Handler().(user.UserHandler); ok {
		uHandler = uh.UserHandler()
	} else {
		uHandler = user.NewHandler(p)
	}

	h := &Handler{
		Handler: uHandler,
		p:       p,
	}

	u := data.LoadOrCreateUser(p.Name())
	u.DisplayName = p.Name()
	u.XUID = p.XUID()
	_ = data.SaveUser(u)

	return h
}

// HandleItemUse ...
func (h *Handler) HandleItemUse(_ *event.Context) {
	held, _ := h.p.HeldItems()

	val, ok := held.Value("lobby")
	if ok {
		switch val {
		case 0:
			h.p.SendForm(form.NewFFA())
		case 1:
			h.p.SendForm(form.NewUnranked())
		case 2:
			h.p.SendForm(form.NewRanked())
		case 8:
			h.p.SendForm(form.NewSettings())
		}
	}

	val, ok = held.Value("queue")
	if ok {
		switch val {
		case 8:
			kit.Apply(kit.Lobby{}, h.p)
			game.DeQueue(h.p)
		}
	}
}

// HandleAttackEntity ...
func (h *Handler) HandleAttackEntity(ctx *event.Context, _ world.Entity, _, _ *float64, _ *bool) {
	ctx.Cancel()
}

// HandleHurt ...
func (h *Handler) HandleHurt(ctx *event.Context, _ *float64, _ *time.Duration, _ world.DamageSource) {
	ctx.Cancel()
}

// UserHandler ...
func (h *Handler) UserHandler() *user.Handler {
	return h.Handler
}

func (h *Handler) SendScoreBoard() {
	l := h.p.Locale()
	u, _ := data.LoadUser(h.p.Name())

	if !u.Settings.Display.Scoreboard {
		return
	}

	var kdr float64
	if u.Stats.Deaths > 0 {
		kdr = float64(u.Stats.Kills) / float64(u.Stats.Deaths)
	} else {
		kdr = float64(u.Stats.Kills)
	}

	sb := scoreboard.New(carrot.GlyphFont(" Moyai"))
	sb.RemovePadding()
	_, _ = sb.WriteString("§r\uE002")

	_, _ = sb.WriteString("\uE142\uE143\uE144\uE143\uE142")
	_, _ = sb.WriteString(text.Colourf("\uE141 K<grey>:</grey> <red>%d</red> D<grey>:</grey> <red>%d</red>", u.Stats.Kills, u.Stats.Deaths))
	_, _ = sb.WriteString(text.Colourf("\uE141 KDR<grey>:</grey> <red>%.2f</red>", kdr))
	_, _ = sb.WriteString(text.Colourf("\uE141 KS<grey>:</grey> <red>%d</red>", u.Stats.KillStreak))

	for i, li := range sb.Lines() {
		if !strings.Contains(li, "\uE002") {
			sb.Set(i, "  "+li)
		}
	}
	_, _ = sb.WriteString("§a")
	_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.footer"))

	_, _ = sb.WriteString("\uE002")
	h.p.RemoveScoreboard()
	h.p.SendScoreboard(sb)
}
func (h *Handler) HandleQuit() {
	user.Remove(h.UserHandler())
	game.DeQueue(h.p)
}
