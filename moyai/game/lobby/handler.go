package lobby

import (
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
	"github.com/moyai-network/carrot/lang"
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
			h.p.SendForm(form.NewQueue())
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

	var kdr float64
	if u.Stats.Deaths > 0 {
		kdr = float64(u.Stats.Kills) / float64(u.Stats.Deaths)
	} else {
		kdr = float64(u.Stats.Kills)
	}

	sb := scoreboard.New(carrot.GlyphFont("PRACTICE"))
	sb.RemovePadding()
	_, _ = sb.WriteString("§r\uE000")

	_, _ = sb.WriteString("\uE142\uE143\uE144\uE143\uE142")
	_, _ = sb.WriteString(text.Colourf("<black>\uE141 </black>K<grey>:</grey> <black>%d</black> D<grey>:</grey> <black>%d</black>", u.Stats.Kills, u.Stats.Deaths))
	_, _ = sb.WriteString(text.Colourf("<black>\uE141 </black>KDR<grey>:</grey> <black>%.2f</black>", kdr))
	_, _ = sb.WriteString(text.Colourf("<black>\uE141 </black>KS<grey>:</grey> <black>%d</black>", u.Stats.KillStreak))

	_, _ = sb.WriteString("\uE000")
	for i, li := range sb.Lines() {
		if !strings.Contains(li, "\uE000") {
			sb.Set(i, "  "+li)
		}
	}
	_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.footer"))
	h.p.RemoveScoreboard()
	h.p.SendScoreboard(sb)
}
func (h *Handler) HandleQuit() {
	user.Remove(h.p)
	game.DeQueue(h.p)
}
